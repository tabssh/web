package config

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"strings"
	"testing"
)

func TestNormalizePath(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"/myadmin/", "myadmin"},
		{"//admin", "admin"},
		{"/my//admin", "my/admin"},
		{"///a//b///", "a/b"},
		{"/server/admin/./x", "server/admin/x"},
		{"admin", "admin"},
	}
	for _, tt := range tests {
		if got := normalizePath(tt.input); got != tt.want {
			t.Errorf("normalizePath(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestValidatePathSegment(t *testing.T) {
	tests := []struct {
		segment string
		wantErr error
	}{
		{"admin", nil},
		{"my-path_2", nil},
		{"", ErrInvalidPath},
		{strings.Repeat("a", 65), ErrPathTooLong},
		{"Admin", ErrInvalidPath},
		{"a b", ErrInvalidPath},
		{"<script>", ErrInvalidPath},
	}
	for _, tt := range tests {
		if err := validatePathSegment(tt.segment); !errors.Is(err, tt.wantErr) {
			t.Errorf("validatePathSegment(%q) = %v, want %v", tt.segment, err, tt.wantErr)
		}
	}
}

func TestSafePath(t *testing.T) {
	tests := []struct {
		input   string
		want    string
		wantErr error
	}{
		{"/myadmin/", "myadmin", nil},
		{"//admin", "admin", nil},
		{"/my//admin", "my/admin", nil},
		{"///a//b///", "a/b", nil},
		{"/../admin", "", ErrPathTraversal},
		{"/server/admin/../secret", "", ErrPathTraversal},
		{"/server/admin/..", "", ErrPathTraversal},
		{"....", "", ErrPathTraversal},
		{"/Admin", "", ErrInvalidPath},
		{"/server/admin/<script>", "", ErrInvalidPath},
		{strings.Repeat("a/", 1500), "", ErrPathTooLong},
	}
	for _, tt := range tests {
		got, err := SafePath(tt.input)
		if !errors.Is(err, tt.wantErr) {
			t.Errorf("SafePath(%q) error = %v, want %v", tt.input, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("SafePath(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestSafeFilePath(t *testing.T) {
	base := t.TempDir()

	got, err := SafeFilePath(base, "uploads/file-1")
	if err != nil {
		t.Fatalf("SafeFilePath valid path: %v", err)
	}
	want := filepath.Join(base, "uploads", "file-1")
	if got != want {
		t.Errorf("SafeFilePath = %q, want %q", got, want)
	}

	if _, err := SafeFilePath(base, "../escape"); !errors.Is(err, ErrPathTraversal) {
		t.Errorf("SafeFilePath traversal error = %v, want ErrPathTraversal", err)
	}
	if _, err := SafeFilePath(base, "Bad Name"); !errors.Is(err, ErrInvalidPath) {
		t.Errorf("SafeFilePath invalid chars error = %v, want ErrInvalidPath", err)
	}
}

func TestPathSecurityMiddleware(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		wantStatus int
		wantPath   string
	}{
		{"collapse double slashes", "/server/admin//config//settings", http.StatusOK, "/server/admin/config/settings"},
		{"collapse leading slashes", "//api///v1//users", http.StatusOK, "/api/v1/users"},
		{"root", "///", http.StatusOK, "/"},
		{"trailing slash preserved", "/docs/", http.StatusOK, "/docs/"},
		{"plain traversal blocked", "/static/../server/admin", http.StatusBadRequest, ""},
		{"encoded traversal blocked", "/api/v1/files/..%2F..%2Fetc/passwd", http.StatusBadRequest, ""},
		{"dot spam blocked", "/server/admin/....//secret", http.StatusBadRequest, ""},
		{"encoded dot blocked", "/files/%2e%2e/secret", http.StatusBadRequest, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotPath string
			handler := PathSecurityMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotPath = r.URL.Path
				w.WriteHeader(http.StatusOK)
			}))
			// Build the URL manually so encoded sequences survive intact
			u, err := url.ParseRequestURI(tt.path)
			if err != nil {
				t.Fatalf("parse %q: %v", tt.path, err)
			}
			req := httptest.NewRequest(http.MethodGet, "http://example.test/", nil)
			req.URL = u
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if tt.wantStatus == http.StatusOK && gotPath != tt.wantPath {
				t.Errorf("normalized path = %q, want %q", gotPath, tt.wantPath)
			}
		})
	}
}
