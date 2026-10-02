# ðŸš€ Universal Go Backend Architecture & Hackathon Starter Template

A production-grade, highly modular, scalable Go backend starter template designed for **Hackathons, Rapid Prototyping, and Microservices**. No project-specific lock-in; completely generalized for any domain.

---

## âš¡ Quick Start (Under 1 Minute)

### 1. Copy or Clone
```bash
xcopy /E /I /H C:\Go\GoProjectStructure C:\Go\MyNewProject
```
Or instruct your AI:
> *"Use `C:\Go\GoProjectStructure` as the base architecture template for this project. Keep the same directory layout, error handling, database pool, router, middlewares, and pagination."*

### 2. Start PostgreSQL (Option A: Docker, Option B: Local DB)
```bash
# Start PostgreSQL 16 + Redis automatically via Docker:
docker compose up -d
```
Or configure `.env` with your own Supabase, Neon, or local database URL:
```ini
PORT=8080
ENV=development
DATABASE_URL=postgres://postgres:password@localhost:5432/appdb?sslmode=disable
JWT_SECRET=super-secret-key-change-in-production
JWT_EXPIRATION_HOURS=72
```

### 3. Run Migrations & Seed
```bash
go run ./cmd/migrate up
go run ./cmd/seed
```
*Seeds default admin: `admin@example.com` / `Admin@123456`*

### 4. Run Server
```bash
.\run-backend.bat
# Or: go run ./cmd/api
```

---

## ðŸŒ Endpoints & Interactive Swagger Documentation

Server starts at: `http://localhost:8080`

- **Interactive Swagger UI**: `http://localhost:8080/swagger/index.html` (Test all endpoints directly from browser)
- **Health Check**: `GET /health`
- **Register**: `POST /api/v1/auth/register`
- **Login**: `POST /api/v1/auth/login`
- **User Profile (Protected)**: `GET /api/v1/auth/profile`

---

## âš¡ Superpower: Instant Domain Scaffolder (Save 15 Mins)

Need to add a new entity (e.g. `Product`, `Order`, `Todo`, `Article`)?
Run this single command:

```powershell
.\scripts\generate-domain.bat Product
```
This automatically generates:
1. `internal/handler/model/product.go` (Entity struct)
2. `internal/handler/dto/product_dto.go` (Create, Update, Response DTOs with validation)
3. `internal/handler/repository/product_repository.go` (CRUD queries + pagination)
4. `internal/handler/services/product_service.go` (Business logic)
5. `internal/handler/controller/product_controller.go` (HTTP handlers with JSON validation)
6. `internal/migration/<timestamp>_create_products_table.up.sql` & `.down.sql` (Versioned SQL migrations)

Then just mount it in `internal/handler/routes/router.go` in 4 lines!

---

## ðŸ› ï¸ Stack & Included Features

| Category | Component | Description |
|---|---|---|
| **Router** | `go-chi/chi/v5` | Idiomatic, net/http compatible, lightweight HTTP router |
| **Database** | `jackc/pgx/v5` (`pgxpool`) | High-performance PostgreSQL driver with connection pooling |
| **Transactions** | `utils.WithTx()` | Auto-commit & rollback transaction wrapper |
| **API Docs** | `Swagger / OpenAPI` | Interactive web UI at `/swagger/index.html` |
| **Pagination** | `reuse.GetPagination()` | Standard `?page=1&limit=10&order=desc` query parser and envelope |
| **Caching** | `reuse.GetCache()` | Thread-safe in-memory key-value cache with TTL |
| **Migrations** | `golang-migrate/migrate/v4` | CLI and auto-migrations for versioned SQL schema |
| **Logger** | `uber-go/zap` | High-performance structured JSON logging |
| **Validation** | `go-playground/validator/v10` | Struct & request validation with formatted field errors |
| **Auth & Security** | `golang-jwt/jwt/v5` + `bcrypt` | Stateless JWT tokens & password hashing (cost 12) |
| **Cloud Storage** | `aws/aws-sdk-go-v2` | S3 & Cloudflare R2 file upload utility |
| **Container** | `Docker & Compose` | Multi-stage Dockerfile + Compose for Postgres & Redis |

---

## ðŸ—ï¸ Folder Structure Overview

```text
â”œâ”€â”€ bin/                        # Compiled executable binaries
â”œâ”€â”€ cmd/
â”‚   â”œâ”€â”€ api/                    # Application entrypoint & dependency injection
â”‚   â”œâ”€â”€ migrate/                # CLI migration tool (up / down)
â”‚   â””â”€â”€ seed/                   # Database seeder (admin user & mock records)
â”œâ”€â”€ docs/                       # Auto-generated Swagger OpenAPI specifications
â”œâ”€â”€ internal/
â”‚   â”œâ”€â”€ handler/                # Domain Layer
â”‚   â”‚   â”œâ”€â”€ controller/         # HTTP handlers (parse JSON, invoke service, send response)
â”‚   â”‚   â”œâ”€â”€ dto/                # Request & response payloads with validation tags
â”‚   â”‚   â”œâ”€â”€ model/              # Database entities
â”‚   â”‚   â”œâ”€â”€ repository/         # Parameterized raw SQL queries using pgxpool
â”‚   â”‚   â”œâ”€â”€ routes/             # Chi router & route grouping
â”‚   â”‚   â””â”€â”€ services/           # Business logic & validations
â”‚   â”œâ”€â”€ middleware/             # Middlewares (Auth JWT, CORS, Zap Logger, Rate Limiter)
â”‚   â”œâ”€â”€ migration/              # SQL migrations (.up.sql / .down.sql) & auto-migration hook
â”‚   â”œâ”€â”€ reuse/                  # General utilities (Pagination, Cache, JWT, Bcrypt, S3, Errors)
â”‚   â”œâ”€â”€ servers/                # Production HTTP server with graceful shutdown
â”‚   â”œâ”€â”€ utils/                  # DB pool, Transaction helper, .env configuration loader
â”‚   â””â”€â”€ worker/                 # Background job runners & scheduled workers
â”œâ”€â”€ scripts/                    # Scaffolding automation (generate-domain.bat / .ps1)
â”œâ”€â”€ docker-compose.yml          # 1-click Postgres + Redis setup
â”œâ”€â”€ Dockerfile                  # Multi-stage minimal production container
â””â”€â”€ taskfile.yaml               # Task runner recipes
```

---

## ðŸ¤– Prompt Template to Give AI for New Features

When you ask an AI assistant to build a new feature or project using this template, paste this prompt:

```text
I want to build a feature using the clean architecture defined in C:\Go\GoProjectStructure.
Follow these rules strictly:
1. Architecture Flow: Route -> Controller -> Service -> Repository -> Database.
2. DTOs: Define request/response structs in `internal/handler/dto/` with `validate:"..."` tags.
3. Validation: Use `reuse.ValidateStruct()` to return standardized 422 errors.
4. Response Format: Always use `reuse.Success()`, `reuse.Created()`, `reuse.PaginatedResponse()`, or `reuse.Error()`.
5. Database: Use raw parameterized SQL queries with `pgxpool.Pool` (no heavy ORM). Use `utils.WithTx()` for multi-step transactions.
6. Security: Protect sensitive routes using `middleware.AuthMiddleware(cfg)` and `middleware.RequireRole(...)`.
```