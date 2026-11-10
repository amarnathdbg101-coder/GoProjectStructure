# 🚀 Complete User Guide: Kaise Use Karein (How to Use This Template)

Is guide me step-by-step bataya gaya hai ki kisi bhi naye project ya hackathon ke liye is template ko kaise use karna hai.

---

## 📑 Table of Contents
1. [Naye Project Ke Liye Setup (Step-by-Step)](#1-naye-project-ke-liye-setup-step-by-step)
2. [Database Start Karna (Docker ya Local)](#2-database-start-karna-docker-ya-local)
3. [Migrations & Seeder Run Karna](#3-migrations--seeder-run-karna)
4. [Server Ko Run Karna (3 Easy Ways)](#4-server-ko-run-karna-3-easy-ways)
5. [APIs Test Karna (Swagger Interactive UI)](#5-apis-test-karna-swagger-interactive-ui)
6. [2 Second Me Naya Feature/Entity Add Karna (Domain Scaffolder)](#6-2-second-me-naya-featureentity-add-karna-domain-scaffolder)
7. [Unit Tests Run Karna](#7-unit-tests-run-karna)
8. [AI Assistants (ChatGPT / Claude / Copilot) Ke Liye Ready Prompt](#8-ai-assistants-ke-liye-ready-prompt)

---

## 1. Naye Project Ke Liye Setup (Step-by-Step)

Jab bhi aapko koi naya project ya hackathon shuru karna ho, is template folder ko copy kar lijiye:

```powershell
# Windows PowerShell me:
xcopy /E /I /H C:\Go\GoProjectStructure C:\Go\MyNewProject
```

Phir apne naye project me jakar:
1. `.env.example` file ko copy karke `.env` banayein:
   ```powershell
   Copy-Item .env.example .env
   ```
2. `.env` me apni settings check kar lijiye (Default settings local development ke liye already configured hain).

---

## 2. Database Start Karna (Docker ya Local)

### Option A: Docker (Sabse Aasan & Recommended)
Agar aapke laptop me PostgreSQL install nahi hai, toh Docker use karein. Ek single command se PostgreSQL 16 aur Redis 7 start ho jayenge:
```powershell
docker compose up -d
```
*(Stop karne ke liye: `docker compose down`)*

### Option B: Local PostgreSQL ya Cloud (Supabase / Neon / Render)
Agar aap local PostgreSQL ya Supabase use kar rahe hain, toh bas `.env` me `DATABASE_URL` update karein:
```ini
DATABASE_URL=postgres://username:password@localhost:5432/appdb?sslmode=disable
```

---

## 3. Migrations & Seeder Run Karna

Database me required tables create karne aur initial admin user create karne ke liye:

```powershell
# 1. Database tables create karein:
go run ./cmd/migrate up
# Ya Taskfile se: task migrate:up

# 2. Default Admin User create karein:
go run ./cmd/seed
# Ya Taskfile se: task seed
```

Default credentials:
- **Email:** `admin@example.com`
- **Password:** `Admin@123456`

---

## 4. Server Ko Run Karna (3 Easy Ways)

Aap in 3 tarikon me se kisi bhi tarike se backend start kar sakte hain:

### Tarika 1: 1-Click Batch Runner (Windows Shortcut)
Bas terminal me likhein ya file par double-click karein:
```powershell
.\run-backend.bat
```

### Tarika 2: Hot-Reload Mode (`air`) — *Best for Coding!*
Code change karne par server automatically 0.5s me reload ho jayega:
```powershell
air
```

### Tarika 3: Standard Go Command
```powershell
go run ./cmd/api
# Ya: task run
```

Server start ho jayega: `http://localhost:8080`

---

## 5. APIs Test Karna (Swagger Interactive UI)

Postman kholne ki bhi zaroorat nahi hai! Apne browser me open karein:

🔗 **`http://localhost:8080/swagger/index.html`**

- **Health Checks:** `/health`, `/health/live`, `/health/ready` check karein.
- **Register User:** `/api/v1/auth/register` se naya user banayein.
- **Login:** `/api/v1/auth/login` se token generate karein.
- **Authorize (JWT):** Top-right me green **Authorize** button pe click karke token daalein (`Bearer <token>`) aur protected `/api/v1/auth/profile` test karein!

---

## 6. 2 Second Me Naya Feature/Entity Add Karna (Domain Scaffolder)

Maan lijiye aapko hackathon me **Product**, **Order**, **Todo**, ya **Article** ka naya feature banana hai:

### Step 1: Scaffolding Command Chalayein
```powershell
.\scripts\generate-domain.bat Product
```

Ye command automatically 2 second me create kar degi (Lean 3-Tier Production Architecture):
1. `internal/model/product.go` ➡️ Database Struct + Request/Response DTOs
2. `internal/repository/product_repository.go` ➡️ Database CRUD + Pagination Queries
3. `internal/handler/product_handler.go` ➡️ HTTP Handlers & Input Validation
4. `internal/migration/<timestamp>_create_products_table.up.sql` ➡️ Table SQL migration

### Step 2: Route Mount Karein
`internal/handler/router.go` me jaakar command dwara batayi gayi lines paste kar dein:
```go
productRepo := repository.NewProductRepository(db)
productH := NewProductHandler(productRepo)

r.Route("/api/v1/products", func(r chi.Router) {
    r.Get("/", productH.List)
    r.Get("/{id}", productH.GetByID)
    r.Post("/", productH.Create)
    r.Delete("/{id}", productH.Delete)
})
```

Bas! Pura naya feature with Database, API endpoints, Validation, aur Pagination ready hai!

---

## 7. Unit Tests Run Karna

Code me koi bug toh nahi hai, ye verify karne ke liye:
```powershell
go test -v ./...
# Ya Taskfile se:
task test
```

---

## 8. AI Assistants Ke Liye Ready Prompt

Agar aap kisi AI (ChatGPT, Claude, Gemini, Cursor) se code likhwa rahe hain, toh use bas ye prompt paste karein:

```text
I am using the clean, lean Go backend template located in C:\Go\GoProjectStructure.
Follow these strict architectural rules:
1. Flow: Route -> Handler -> Service -> Repository -> Database (or Handler -> Repository for simple CRUD).
2. Models & DTOs: Define all entity and request/response schemas together in `internal/model/<entity>.go` with `validate:"..."` tags.
3. Validation: Use `validator.ValidateStruct()` to return standardized 422 JSON errors.
4. Response Format: Always use `handler.Success()`, `handler.Created()`, `pagination.Respond()`, or `handler.Error()`.
5. Database: Use raw parameterized SQL queries with `pgxpool.Pool` (no heavy ORM). Use `database.WithTx()` for multi-step atomic transactions.
6. Security: Protect sensitive routes using `middleware.JWTAuth(cfg.JwtSecret)`.
```
