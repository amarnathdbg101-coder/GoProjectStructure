// Minimal Go Clean Architecture Starter Demo
package main

import (
	"log"
	"net/http"

	"github.com/amarnathdbg101-coder/GoProjectStructure/internal/handler"
	"github.com/amarnathdbg101-coder/GoProjectStructure/response"
	"github.com/amarnathdbg101-coder/GoProjectStructure/server"
	"github.com/go-chi/chi/v5"
)

func main() {
	// 1. Initialize clean router
	r := handler.NewRouter()

	// 2. Mount domain routes (Example: /api/v1/welcome)
	r.Route("/api/v1", func(api chi.Router) {
		api.Get("/welcome", func(w http.ResponseWriter, r *http.Request) {
			response.Success(w, "Welcome to Go Clean Architecture!", map[string]string{
				"starter": "GoProjectStructure",
				"version": "v1.0.0",
			})
		})
	})

	// 3. Initialize and start HTTP server with graceful shutdown
	srv := server.New(r, server.Config{
		Port: "8080",
	})

	log.Printf("Starting server on http://localhost%s ...", srv.Addr())
	if err := srv.Start(); err != nil {
		log.Fatalf("Server stopped with error: %v", err)
	}
}
