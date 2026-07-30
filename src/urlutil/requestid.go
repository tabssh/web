// Request ID handling per AI.md PART 8 (Request ID Handling): every
// request carries a UUID v4 request ID for tracing, echoed in the
// X-Request-ID response header and stored in the request context.
package urlutil

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

// contextKey is a private context key type to avoid collisions.
type contextKey string

// requestIDKey stores the request ID in the request context.
const requestIDKey contextKey = "request_id"

// RequestIDMiddleware ensures every request has a valid request ID,
// accepting X-Request-ID > X-Correlation-ID > X-Trace-ID from the client
// and generating a fresh UUID when absent or invalid.
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for an existing request ID from client or upstream proxy.
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = r.Header.Get("X-Correlation-ID")
		}
		if requestID == "" {
			requestID = r.Header.Get("X-Trace-ID")
		}

		// Generate a new ID if none provided or invalid.
		if requestID == "" || !isValidUUID(requestID) {
			requestID = uuid.New().String()
		}

		// Add to response headers.
		w.Header().Set("X-Request-ID", requestID)

		// Add to request context for logging and downstream calls.
		ctx := context.WithValue(r.Context(), requestIDKey, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequestIDFromContext returns the request ID stored by the middleware,
// or the empty string when none is present.
func RequestIDFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(requestIDKey).(string); ok {
		return v
	}
	return ""
}

// isValidUUID reports whether s parses as a UUID.
func isValidUUID(s string) bool {
	_, err := uuid.Parse(s)
	return err == nil
}
