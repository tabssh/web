package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServerAppliesRequestIDMiddleware(t *testing.T) {
	s := New("127.0.0.1:0")
	s.Mux().HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	tests := []struct {
		name     string
		clientID string
	}{
		{"generated when absent", ""},
		{"echoed when valid", "550e8400-e29b-41d4-a716-446655440000"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "http://example.com/ping", nil)
			if tt.clientID != "" {
				r.Header.Set("X-Request-ID", tt.clientID)
			}
			w := httptest.NewRecorder()
			s.HTTPServer().Handler.ServeHTTP(w, r)
			got := w.Header().Get("X-Request-ID")
			if got == "" {
				t.Fatal("X-Request-ID missing from response")
			}
			if tt.clientID != "" && got != tt.clientID {
				t.Errorf("request ID = %q, want %q", got, tt.clientID)
			}
		})
	}
}

func TestServerAddr(t *testing.T) {
	s := New("0.0.0.0:8080")
	if s.HTTPServer().Addr != "0.0.0.0:8080" {
		t.Errorf("Addr = %q, want 0.0.0.0:8080", s.HTTPServer().Addr)
	}
}
