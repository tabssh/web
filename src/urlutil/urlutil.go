// Package urlutil resolves {proto}, {fqdn}, {port} per request and builds
// absolute URLs, per AI.md PART 8 (URL Variables, Reverse Proxy Header
// Support). Reverse-proxy headers are preferred; :80 and :443 are always
// stripped from generated URLs.
package urlutil

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/tabssh/web/src/mode"
)

// GetURLVars returns resolved URL variables from the request. Reverse
// proxy headers are checked first; port is the empty string for 80/443
// (always stripped).
func GetURLVars(r *http.Request) (proto, fqdn, port string) {
	proto = detectProto(r)
	fqdn, port = detectHostPort(r)
	if port == "" {
		port = detectPort(r, proto)
	}
	// Default ports are never included in generated URLs.
	if port == "80" || port == "443" {
		port = ""
	}
	return proto, fqdn, port
}

// BuildURL constructs a full URL with automatic default-port stripping;
// :80 and :443 are NEVER included.
func BuildURL(r *http.Request, path string) string {
	proto, fqdn, port := GetURLVars(r)
	if port == "" {
		return fmt.Sprintf("%s://%s%s", proto, fqdn, path)
	}
	return fmt.Sprintf("%s://%s:%s%s", proto, fqdn, port, path)
}

// detectProto resolves {proto}: X-Forwarded-Proto > X-Forwarded-Ssl >
// X-Url-Scheme > Front-End-Https > TLS detection > http.
func detectProto(r *http.Request) string {
	if v := strings.ToLower(r.Header.Get("X-Forwarded-Proto")); v == "https" || v == "http" {
		return v
	}
	switch strings.ToLower(r.Header.Get("X-Forwarded-Ssl")) {
	case "on":
		return "https"
	case "off":
		return "http"
	}
	if v := strings.ToLower(r.Header.Get("X-Url-Scheme")); v == "https" || v == "http" {
		return v
	}
	if strings.EqualFold(r.Header.Get("Front-End-Https"), "on") {
		return "https"
	}
	if r.TLS != nil {
		return "https"
	}
	return "http"
}

// detectHostPort resolves {fqdn} (and an embedded port, if any):
// X-Forwarded-Host > X-Real-Host > X-Original-Host > Host header >
// DOMAIN env > os.Hostname() > HOSTNAME env > localhost. Invalid env
// values are skipped silently and detection continues.
func detectHostPort(r *http.Request) (string, string) {
	for _, h := range []string{"X-Forwarded-Host", "X-Real-Host", "X-Original-Host"} {
		if v := strings.TrimSpace(r.Header.Get(h)); v != "" {
			return splitHostPort(v)
		}
	}
	if r.Host != "" {
		return splitHostPort(r.Host)
	}
	devMode := mode.IsAppModeDev()
	if v := strings.TrimSpace(os.Getenv("DOMAIN")); v != "" && IsValidHost(v, devMode, "tabssh") {
		return strings.ToLower(v), ""
	}
	if h, err := os.Hostname(); err == nil && IsValidHost(h, devMode, "tabssh") {
		return strings.ToLower(h), ""
	}
	if v := strings.TrimSpace(os.Getenv("HOSTNAME")); v != "" && IsValidHost(v, devMode, "tabssh") {
		return strings.ToLower(v), ""
	}
	return "localhost", ""
}

// detectPort resolves {port}: X-Forwarded-Port > X-Real-Port > implied by
// protocol (returned empty so default ports stay stripped).
func detectPort(r *http.Request, proto string) string {
	for _, h := range []string{"X-Forwarded-Port", "X-Real-Port"} {
		if v := strings.TrimSpace(r.Header.Get(h)); v != "" {
			return v
		}
	}
	return ""
}

// splitHostPort splits an optional :port suffix off a host value,
// tolerating bare hosts and bracketed IPv6 literals.
func splitHostPort(hostport string) (string, string) {
	host, port, err := net.SplitHostPort(hostport)
	if err != nil {
		return strings.ToLower(hostport), ""
	}
	return strings.ToLower(host), port
}

// ClientIP resolves the client IP for logging, rate limiting, and GeoIP.
// Proxy headers are honored only when trustedPeer is true (the immediate
// TCP peer passed the trusted_proxies gate); otherwise r.RemoteAddr is
// used directly. Priority: CF-Connecting-IP > True-Client-IP > X-Real-IP
// > X-Forwarded-For (leftmost) > X-Client-IP > r.RemoteAddr.
func ClientIP(r *http.Request, trustedPeer bool) string {
	if trustedPeer {
		for _, h := range []string{"CF-Connecting-IP", "True-Client-IP", "X-Real-IP"} {
			if v := strings.TrimSpace(r.Header.Get(h)); v != "" && net.ParseIP(v) != nil {
				return v
			}
		}
		if v := r.Header.Get("X-Forwarded-For"); v != "" {
			first := strings.TrimSpace(strings.Split(v, ",")[0])
			if net.ParseIP(first) != nil {
				return first
			}
		}
		if v := strings.TrimSpace(r.Header.Get("X-Client-IP")); v != "" && net.ParseIP(v) != nil {
			return v
		}
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

// ExtractAuthToken extracts the auth token from a request using the PART 8
// priority order (first found wins): Authorization (Bearer stripped) >
// X-API-Key/API-Key variants > X-Auth-Token/X-Access-Token > X-Token/Token
// > ?token= query parameter (least preferred).
func ExtractAuthToken(r *http.Request) string {
	if v := strings.TrimSpace(r.Header.Get("Authorization")); v != "" {
		if len(v) > 7 && strings.EqualFold(v[:7], "bearer ") {
			return strings.TrimSpace(v[7:])
		}
		return v
	}
	for _, h := range []string{"X-API-Key", "X-Api-Key", "API-Key", "ApiKey", "X-Auth-Token", "X-Access-Token", "X-Token", "Token"} {
		if v := strings.TrimSpace(r.Header.Get(h)); v != "" {
			return v
		}
	}
	return strings.TrimSpace(r.URL.Query().Get("token"))
}
