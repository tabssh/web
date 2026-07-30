package urlutil

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetURLVarsAndBuildURL(t *testing.T) {
	tests := []struct {
		name    string
		host    string
		headers map[string]string
		path    string
		want    string
	}{
		{
			"plain http host",
			"example.com",
			nil,
			"/foo",
			"http://example.com/foo",
		},
		{
			"forwarded proto https strips 443",
			"example.com",
			map[string]string{"X-Forwarded-Proto": "https", "X-Forwarded-Port": "443"},
			"/foo",
			"https://example.com/foo",
		},
		{
			"port 80 stripped",
			"example.com:80",
			nil,
			"/",
			"http://example.com/",
		},
		{
			"custom port kept",
			"example.com:8080",
			nil,
			"/api",
			"http://example.com:8080/api",
		},
		{
			"forwarded host wins",
			"internal:9999",
			map[string]string{"X-Forwarded-Host": "public.example.com", "X-Forwarded-Proto": "https"},
			"/x",
			"https://public.example.com/x",
		},
		{
			"forwarded ssl on",
			"example.com",
			map[string]string{"X-Forwarded-Ssl": "on"},
			"",
			"https://example.com",
		},
		{
			"forwarded port kept",
			"example.com",
			map[string]string{"X-Forwarded-Port": "8443", "X-Forwarded-Proto": "https"},
			"/p",
			"https://example.com:8443/p",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "http://placeholder/", nil)
			r.Host = tt.host
			for k, v := range tt.headers {
				r.Header.Set(k, v)
			}
			if got := BuildURL(r, tt.path); got != tt.want {
				t.Errorf("BuildURL() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestClientIP(t *testing.T) {
	tests := []struct {
		name        string
		remoteAddr  string
		headers     map[string]string
		trustedPeer bool
		want        string
	}{
		{"remote addr fallback", "203.0.113.7:5321", nil, false, "203.0.113.7"},
		{"untrusted peer ignores headers", "203.0.113.7:5321", map[string]string{"X-Real-IP": "198.51.100.1"}, false, "203.0.113.7"},
		{"cf connecting ip wins", "10.0.0.1:80", map[string]string{"CF-Connecting-IP": "198.51.100.2", "X-Real-IP": "198.51.100.3"}, true, "198.51.100.2"},
		{"xff leftmost", "10.0.0.1:80", map[string]string{"X-Forwarded-For": "198.51.100.4, 10.0.0.2"}, true, "198.51.100.4"},
		{"invalid header ip skipped", "10.0.0.1:80", map[string]string{"X-Real-IP": "not-an-ip"}, true, "10.0.0.1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "http://example.com/", nil)
			r.RemoteAddr = tt.remoteAddr
			for k, v := range tt.headers {
				r.Header.Set(k, v)
			}
			if got := ClientIP(r, tt.trustedPeer); got != tt.want {
				t.Errorf("ClientIP() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestExtractAuthToken(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		headers map[string]string
		want    string
	}{
		{"bearer stripped", "/", map[string]string{"Authorization": "Bearer abc123"}, "abc123"},
		{"authorization wins over api key", "/", map[string]string{"Authorization": "Bearer a", "X-API-Key": "b"}, "a"},
		{"api key", "/", map[string]string{"X-API-Key": "key1"}, "key1"},
		{"auth token header", "/", map[string]string{"X-Auth-Token": "tok"}, "tok"},
		{"query param last resort", "/?token=qtok", nil, "qtok"},
		{"nothing", "/", nil, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "http://example.com"+tt.url, nil)
			for k, v := range tt.headers {
				r.Header.Set(k, v)
			}
			if got := ExtractAuthToken(r); got != tt.want {
				t.Errorf("ExtractAuthToken() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestIsValidHost(t *testing.T) {
	tests := []struct {
		name     string
		host     string
		devMode  bool
		wantProd bool
		wantDev  bool
	}{
		{"valid etld+1", "api.example.com", false, true, true},
		{"co.uk domain", "my.server.domain.co.uk", false, true, true},
		{"dev tld local", "dev.local", false, false, true},
		{"dev tld test", "app.test", false, false, true},
		{"project tld", "app.tabssh", false, false, true},
		{"localhost", "localhost", false, false, true},
		{"suffix only", "co.uk", false, false, false},
		{"ip address", "192.168.1.1", false, false, false},
		{"single label", "myhost", false, false, false},
		{"empty", "", false, false, false},
		{"onion always valid", "abc.onion", false, true, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidHost(tt.host, false, "tabssh"); got != tt.wantProd {
				t.Errorf("IsValidHost(%q, prod) = %v, want %v", tt.host, got, tt.wantProd)
			}
			if got := IsValidHost(tt.host, true, "tabssh"); got != tt.wantDev {
				t.Errorf("IsValidHost(%q, dev) = %v, want %v", tt.host, got, tt.wantDev)
			}
		})
	}
}

func TestIsValidSSLHost(t *testing.T) {
	tests := []struct {
		name string
		host string
		want bool
	}{
		{"public domain", "example.com", true},
		{"onion rejected", "abc.onion", false},
		{"dev tld rejected", "dev.local", false},
		{"localhost rejected", "localhost", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidSSLHost(tt.host); got != tt.want {
				t.Errorf("IsValidSSLHost(%q) = %v, want %v", tt.host, got, tt.want)
			}
		})
	}
}

func TestRequestIDMiddleware(t *testing.T) {
	tests := []struct {
		name      string
		headers   map[string]string
		wantEcho  string
		wantFresh bool
	}{
		{"client uuid honored", map[string]string{"X-Request-ID": "550e8400-e29b-41d4-a716-446655440000"}, "550e8400-e29b-41d4-a716-446655440000", false},
		{"correlation id honored", map[string]string{"X-Correlation-ID": "550e8400-e29b-41d4-a716-446655440001"}, "550e8400-e29b-41d4-a716-446655440001", false},
		{"trace id honored", map[string]string{"X-Trace-ID": "550e8400-e29b-41d4-a716-446655440002"}, "550e8400-e29b-41d4-a716-446655440002", false},
		{"invalid format regenerated", map[string]string{"X-Request-ID": "not-a-uuid"}, "", true},
		{"missing generated", nil, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ctxID string
			handler := RequestIDMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctxID = RequestIDFromContext(r.Context())
			}))
			r := httptest.NewRequest(http.MethodGet, "http://example.com/", nil)
			for k, v := range tt.headers {
				r.Header.Set(k, v)
			}
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, r)
			got := w.Header().Get("X-Request-ID")
			if got == "" {
				t.Fatal("X-Request-ID response header missing")
			}
			if tt.wantFresh {
				if !isValidUUID(got) {
					t.Errorf("generated request ID %q is not a valid UUID", got)
				}
				if tt.headers["X-Request-ID"] != "" && got == tt.headers["X-Request-ID"] {
					t.Error("invalid client ID was not regenerated")
				}
			} else if got != tt.wantEcho {
				t.Errorf("request ID = %q, want %q", got, tt.wantEcho)
			}
			if ctxID != got {
				t.Errorf("context request ID = %q, header = %q", ctxID, got)
			}
		})
	}
}
