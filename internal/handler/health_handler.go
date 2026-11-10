package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var startTime = time.Now()

// HealthHandler provides health check, liveness, and readiness endpoints.
type HealthHandler struct {
	db *pgxpool.Pool
}

// NewHealthHandler constructs a HealthHandler with database awareness.
func NewHealthHandler(db *pgxpool.Pool) *HealthHandler {
	return &HealthHandler{db: db}
}

// Check returns overall system health, uptime, and database connection status.
func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	dbStatus := "connected"
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if h.db == nil {
		dbStatus = "disconnected (offline mode)"
	} else if err := h.db.Ping(ctx); err != nil {
		dbStatus = "unreachable: " + err.Error()
	}

	Success(w, "Server is operational", map[string]interface{}{
		"status":    "healthy",
		"database":  dbStatus,
		"uptime":    time.Since(startTime).String(),
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// Live is the Kubernetes/Cloud liveness probe verifying the process is running.
func (h *HealthHandler) Live(w http.ResponseWriter, r *http.Request) {
	Success(w, "Server process is alive", map[string]interface{}{
		"status": "alive",
		"uptime": time.Since(startTime).String(),
	})
}

// Ready is the readiness probe ensuring the database is healthy before receiving live traffic.
func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
	if h.db == nil {
		Error(w, http.StatusServiceUnavailable, "Database pool is not initialized")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := h.db.Ping(ctx); err != nil {
		Error(w, http.StatusServiceUnavailable, "Database ping failed: "+err.Error())
		return
	}

	Success(w, "Server is ready to accept traffic", map[string]interface{}{
		"status": "ready",
	})
}
