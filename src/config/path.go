package config

import (
	"errors"
	"net/http"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

// Path security errors, per AI.md PART 5 Path Normalization & Validation.
var (
	ErrPathTraversal = errors.New("path traversal attempt detected")
	ErrInvalidPath   = errors.New("invalid path characters")
	ErrPathTooLong   = errors.New("path exceeds maximum length")
)

// validPathSegment matches a valid path segment: lowercase alphanumeric,
// hyphens, underscores.
var validPathSegment = regexp.MustCompile(`^[a-z0-9_-]+$`)

// normalizePath cleans a path for safe use: strips leading/trailing slashes,
// collapses multiple slashes, removes traversal (.., .), and returns an empty
// string for invalid input.
func normalizePath(input string) string {
	if input == "" {
		return ""
	}

	// path.Clean handles .., ., and //
	cleaned := path.Clean(input)

	// Strip leading/trailing slashes
	cleaned = strings.Trim(cleaned, "/")

	// Reject if still contains .. after cleaning (defense in depth)
	if strings.Contains(cleaned, "..") {
		return ""
	}

	return cleaned
}

// validatePathSegment checks a single path segment (e.g. "admin" in
// "/server/admin/dashboard").
func validatePathSegment(segment string) error {
	if segment == "" {
		return ErrInvalidPath
	}
	if len(segment) > 64 {
		return ErrPathTooLong
	}
	if !validPathSegment.MatchString(segment) {
		return ErrInvalidPath
	}
	if segment == "." || segment == ".." {
		return ErrPathTraversal
	}
	return nil
}

// validatePath checks an entire path.
func validatePath(p string) error {
	if len(p) > 2048 {
		return ErrPathTooLong
	}

	// Check for traversal attempts before normalization
	if strings.Contains(p, "..") {
		return ErrPathTraversal
	}

	// Check each segment
	segments := strings.Split(strings.Trim(p, "/"), "/")
	for _, seg := range segments {
		if seg == "" {
			// Skip empty segments produced by //
			continue
		}
		if err := validatePathSegment(seg); err != nil {
			return err
		}
	}

	return nil
}

// SafePath normalizes and validates a path, returning an error if invalid.
func SafePath(input string) (string, error) {
	if err := validatePath(input); err != nil {
		return "", err
	}
	return normalizePath(input), nil
}

// SafeFilePath ensures a user-supplied path stays within the base directory.
func SafeFilePath(baseDir, userPath string) (string, error) {
	// Normalize user input
	safe, err := SafePath(userPath)
	if err != nil {
		return "", err
	}

	// Construct full path
	fullPath := filepath.Join(baseDir, safe)

	// Resolve to absolute
	absPath, err := filepath.Abs(fullPath)
	if err != nil {
		return "", err
	}

	absBase, err := filepath.Abs(baseDir)
	if err != nil {
		return "", err
	}

	// Verify path is still within base
	if !strings.HasPrefix(absPath, absBase+string(filepath.Separator)) && absPath != absBase {
		return "", ErrPathTraversal
	}

	return absPath, nil
}

// PathSecurityMiddleware normalizes request paths and blocks traversal
// attempts. Per AI.md PART 5 Middleware Order it MUST execute third in the
// chain: after URLNormalizeMiddleware and RequestIDMiddleware, before
// security headers, allowlist, blocklist, rate limiting, GeoIP, auth, and
// logging (those middlewares belong to later PARTs).
func PathSecurityMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		original := r.URL.Path

		// r.URL.Path is already decoded by net/http; check RawPath too so
		// encoded traversal (%2e%2e) is caught
		rawPath := r.URL.RawPath
		if rawPath == "" {
			rawPath = r.URL.Path
		}

		// Block path traversal attempts, encoded (%2e = .) and decoded
		if strings.Contains(original, "..") ||
			strings.Contains(rawPath, "..") ||
			strings.Contains(strings.ToLower(rawPath), "%2e") {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		// Normalize the path
		cleaned := path.Clean(original)

		// Ensure leading slash
		if !strings.HasPrefix(cleaned, "/") {
			cleaned = "/" + cleaned
		}

		// Preserve trailing slash for directory paths
		if original != "/" && strings.HasSuffix(original, "/") && !strings.HasSuffix(cleaned, "/") {
			cleaned += "/"
		}

		// Update request
		r.URL.Path = cleaned

		next.ServeHTTP(w, r)
	})
}
