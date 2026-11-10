package handler

import (
	"net/http"
	"time"

	_ "github.com/amarnathdbg101-coder/GoProjectStructure/docs"
	"github.com/amarnathdbg101-coder/GoProjectStructure/pkg/config"
	"github.com/amarnathdbg101-coder/GoProjectStructure/pkg/middleware"
	"github.com/amarnathdbg101-coder/GoProjectStructure/internal/repository"
	"github.com/amarnathdbg101-coder/GoProjectStructure/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
)

// RouteSetup initializes all global middlewares, health probes, and API domain routes.
func RouteSetup(db *pgxpool.Pool, logger *zap.Logger, cfg config.Config) *chi.Mux {
	r := chi.NewRouter()

	// Global Middlewares Pipeline (Order: RequestID -> Recoverer -> Security -> Logger -> CORS -> RateLimit)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer(logger))
	r.Use(middleware.SecurityHeaders)
	r.Use(middleware.Logger)
	r.Use(middleware.CORS)
	limiter := middleware.NewIPRateLimiter(100, time.Minute)
	r.Use(limiter.Middleware())

	// 1. Interactive Swagger Documentation UI
	r.Get("/swagger/*", httpSwagger.WrapHandler)
	r.Get("/swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/index.html", http.StatusMovedPermanently)
	})

	// 2. Health & Monitoring Probes
	healthH := NewHealthHandler(db)
	r.Get("/health", healthH.Check)
	r.Get("/health/live", healthH.Live)
	r.Get("/health/ready", healthH.Ready)

	// 3. Domain Repositories
	userRepo := repository.NewUserRepository(db, logger)

	// 4. Domain Services
	authService := service.NewAuthService(userRepo, cfg.JwtSecret)

	// 5. Domain Handlers
	authH := NewAuthHandler(authService)

	// 6. API v1 Routes
	r.Route("/api/v1", func(r chi.Router) {
		// Public Auth routes
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authH.Register)
			r.Post("/login", authH.Login)

			// Protected user profile route
			r.With(middleware.JWTAuth(cfg.JwtSecret)).Get("/profile", authH.Me)
		})
	})

	return r
}
