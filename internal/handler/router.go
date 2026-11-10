// Package handler provides standard HTTP routing and endpoint registration.
package handler

import (
	"net/http"

	"github.com/amarnathdbg101-coder/GoProjectStructure/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// NewRouter constructs a lightweight Chi router with minimal standard middleware.
// It has zero hardcoded dependencies on specific loggers, databases, or configs,
// allowing you to compose your own preferred libraries (from C:\Library or external).
func NewRouter() *chi.Mux {
	r := chi.NewRouter()

	// Lightweight standard Chi middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	// Basic liveness probe
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		response.Success(w, "server is operational", map[string]string{
			"status": "healthy",
		})
	})

	return r
}
