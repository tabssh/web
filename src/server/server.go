// Package server holds the HTTP server wrapper the binary starts in
// PHASE 5 of the startup sequence (AI.md PART 8). Route registration is
// the PART 13/14 seam: handlers register on Mux(); the request-ID
// middleware from PART 8 is always applied at the outermost layer.
package server

import (
	"net/http"
	"time"

	"github.com/tabssh/web/src/urlutil"
)

// Server wraps the standard library HTTP server with the project mux.
type Server struct {
	httpServer *http.Server
	mux        *http.ServeMux
}

// New creates a server listening on addr with the PART 8 request-ID
// middleware wrapping the mux and sane connection timeouts.
func New(addr string) *Server {
	mux := http.NewServeMux()
	return &Server{
		mux: mux,
		httpServer: &http.Server{
			Addr:              addr,
			Handler:           urlutil.RequestIDMiddleware(mux),
			ReadHeaderTimeout: 10 * time.Second,
			IdleTimeout:       120 * time.Second,
		},
	}
}

// Mux returns the route mux; PART 13/14 handlers register here.
func (s *Server) Mux() *http.ServeMux {
	return s.mux
}

// HTTPServer returns the underlying http.Server for signal handling and
// graceful shutdown.
func (s *Server) HTTPServer() *http.Server {
	return s.httpServer
}

// ListenAndServe starts serving; it returns http.ErrServerClosed on
// graceful shutdown.
func (s *Server) ListenAndServe() error {
	return s.httpServer.ListenAndServe()
}
