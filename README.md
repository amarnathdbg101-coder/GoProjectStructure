# 🚀 GoProjectStructure (Minimal Clean Architecture Starter)

A lightweight, unopinionated **Go Clean Architecture Starter Kit & Library** designed for building fast, maintainable backend services and REST APIs in Go.

Zero bloat. Zero duplicate loggers, hardcoded configs, or redundant middlewares. Built specifically to be composed with your modular library collection (`C:\Library` / `@amarnathdbg101-coder/go/...`).

---

## 🎯 Philosophy

- **Zero Bloat**: No hardcoded Zap, no hardcoded Godotenv, no hardcoded database drivers, no hardcoded migrations.
- **Unopinionated**: You pick your database (`pgx`, `database/sql`, `mongo`), your logger (`go-colorlog`, `slog`), your config (`go-env`), and your auth (`go-jwt`).
- **Clean Contracts**: Provides clean abstractions for Domain, Service, Repository, and Transport layers.
- **Production Server Lifecycle**: Built-in graceful HTTP server runner (`server.New(handler)`) with context-aware signal trapping (SIGINT / SIGTERM).
- **Uniform Responses**: Standardized JSON response envelope (`response.Success`, `response.Error`, `response.Created`).

---

## 📁 Architecture Overview

```
GoProjectStructure/
├── cmd/
│   └── api/
│       └── main.go          # Minimal runnable server entrypoint
├── internal/
│   ├── domain/              # Business entities & error definitions (pure Go)
│   ├── repository/          # Storage interface contracts
│   ├── service/             # Business logic interface contracts
│   └── handler/
│       └── router.go        # Chi router setup without hardcoded bloat
├── response/
│   ├── response.go          # Standard JSON API envelopes (Success, Error, Message)
│   └── response_test.go     # Unit tests
├── server/
│   ├── server.go            # Production HTTP server runner with Graceful Shutdown
│   └── server_test.go       # Unit tests
├── go.mod                   # Zero bloat (only go-chi/chi/v5)
└── README.md
```

---

## ⚡ Quick Start

### 1. Installation

```bash
go get github.com/amarnathdbg101-coder/GoProjectStructure@latest
```

Or clone directly as a project starter template:
```bash
git clone git@github.com:amarnathdbg101-coder/GoProjectStructure.git my-project
cd my-project
```

### 2. Plug in Your Modular Libraries

Compose your project using your specialized libraries:

```go
package main

import (
	"log"
	"net/http"

	"github.com/amarnathdbg101-coder/GoProjectStructure/internal/handler"
	"github.com/amarnathdbg101-coder/GoProjectStructure/response"
	"github.com/amarnathdbg101-coder/GoProjectStructure/server"
	"github.com/go-chi/chi/v5"

	// Your custom library collection from C:\Library:
	// "github.com/amarnathdbg101-coder/go/go-colorlog"
	// "github.com/amarnathdbg101-coder/go/go-env"
	// "github.com/amarnathdbg101-coder/go/go-cors"
	// "github.com/amarnathdbg101-coder/go/go-shutdown"
	// "github.com/amarnathdbg101-coder/go/go-jwt"
	// "github.com/amarnathdbg101-coder/go/go-tinylfu"
)

func main() {
	// 1. Initialize lightweight router
	r := handler.NewRouter()

	// 2. Mount domain routes
	r.Route("/api/v1", func(api chi.Router) {
		api.Get("/welcome", func(w http.ResponseWriter, r *http.Request) {
			response.Success(w, "Welcome to Go Clean Architecture!", map[string]string{
				"version": "1.0.0",
			})
		})
	})

	// 3. Start server with automatic graceful shutdown
	srv := server.New(r, server.Config{
		Port: "8080",
	})

	log.Printf("Server listening on http://localhost%s", srv.Addr())
	if err := srv.Start(); err != nil {
		log.Fatalf("Server stopped: %v", err)
	}
}
```

### 3. Run the Server

```bash
go run ./cmd/api
```

---

## 🧪 Testing

```bash
go test -v ./...
```

---

## 📄 License

MIT License. Crafted with clean architecture principles.
