# 🏛️ Architecture & Coding Standards (Lean Production Architecture)

This project follows a **Lean, Pragmatic 3-Tier Layered Architecture** optimized for high performance, maintainability, and rapid developer velocity.

---

## 🔄 Request Lifecycle & Data Flow

```text
[ Client Request ]
       │
       ▼
[ Middlewares ]  ───► (RequestID, Recoverer, Security, Logger, CORS, RateLimiter)
       │
       ▼
[ Handler ]      ───► Chi Router + Endpoint Handlers (Binds JSON, Validates, Sends JSON response)
       │
       ▼
[ Service ]      ───► Core Business Logic (Passwords, JWT Tokens, Complex Rules)
       │
       ▼
[ Repository ]   ───► Parameterized Raw SQL Queries via pgxpool.Pool
       │
       ▼
[ PostgreSQL ]   ───► Relational Database
```

---

## 📁 Clean Directory Layout

```text
C:\Go\GoProjectStructure\
├── cmd/
│   ├── api/main.go              # Server bootstrap: config, db, router, graceful shutdown
│   ├── migrate/main.go          # Standalone migration CLI runner
│   └── seed/main.go             # Database seeder (default admin user)
│
├── internal/
│   ├── config/
│   │   └── config.go            # Type-safe .env configuration loader
│   ├── database/
│   │   ├── connect.go           # High-performance pgxpool connection
│   │   ├── tx.go                # Atomic transaction helper: database.WithTx()
│   │   └── migration.go         # Automatic database schema migration runner
│   ├── model/
│   │   └── user.go              # Database entity structs & Request/Response DTOs
│   ├── repository/
│   │   └── user_repository.go   # Pure SQL queries using pgxpool.Pool
│   ├── service/
│   │   └── auth_service.go      # Core business rules, password hash, JWT generation
│   ├── handler/
│   │   ├── router.go            # Chi router setup & route registration
│   │   ├── auth_handler.go      # HTTP handlers for auth endpoints
│   │   ├── health_handler.go    # System health, live & ready probes
│   │   └── response.go          # Standard JSON envelopes: Success(), Created(), Error()
│   ├── middleware/              # Production middlewares: auth, cors, logger, ratelimit, recoverer
│   ├── worker/                  # Background worker routines with cancelable context
│   └── pkg/                     # Standalone reusable utility packages:
│       ├── hash/                # Bcrypt password hashing & checking
│       ├── jwt/                 # JWT token generation & parsing
│       ├── validator/           # Struct validation with human-readable error formatting
│       ├── pagination/          # Query pagination extractor & envelope builder
│       ├── cache/               # Concurrency-safe in-memory cache with TTL
│       ├── otp/                 # Cryptographically secure 6-digit OTP generator
│       └── slug/                # URL-safe slugifier
│
└── migrations/                  # Raw SQL migration files (*.up.sql, *.down.sql)
```

---

## 🛠️ General Utility Helpers

### 1. Database Transactions (`internal/database/tx.go`)
Wrap multiple database operations in an atomic transaction:
```go
err := database.WithTx(ctx, db, func(tx pgx.Tx) error {
    if err := repo.DeductBalanceTx(ctx, tx, userID, amount); err != nil {
        return err
    }
    return repo.CreateOrderTx(ctx, tx, order)
})
```

### 2. Standardized JSON Responses (`internal/handler/response.go`)
```go
handler.Success(w, "Profile fetched successfully", user)
handler.Created(w, "User registered successfully", user)
handler.Error(w, http.StatusBadRequest, "Invalid payload")
```

### 3. Pagination & Sorting (`internal/pkg/pagination/pagination.go`)
```go
params := pagination.GetParams(r) // parses ?page=1&limit=10&order=desc&search=foo
items, total, err := repo.FindAll(r.Context(), params)
pagination.Respond(w, items, total, params)
```

---

## 🚀 Domain Generation Command

To scaffold any new entity in 2 seconds:
```powershell
.\scripts\generate-domain.bat Product
```
Creates:
1. `internal/model/product.go` ➡️ Entity Struct + Request/Response DTOs
2. `internal/repository/product_repository.go` ➡️ Database CRUD + Pagination Queries
3. `internal/handler/product_handler.go` ➡️ HTTP Handlers & Input Validation
4. `internal/migration/<timestamp>_create_products_table.up.sql` ➡️ Table SQL migration
