package controller

import (
	"context"
	"net/http"
	"goProjectStructure/internal/reuse"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var startTime = time.Now()

type HealthController struct {
	db *pgxpool.Pool
}

func NewHealthController(db *pgxpool.Pool) *HealthController {
	return &HealthController{db: db}
}

// Check returns system health, uptime, and database connectivity status
func (c *HealthController) Check(w http.ResponseWriter, r *http.Request) {
	dbStatus := "connected"
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := c.db.Ping(ctx); err != nil {
		dbStatus = "unreachable: " + err.Error()
	}

	reuse.Success(w, "Server is operational", map[string]interface{}{
		"status":    "healthy",
		"database":  dbStatus,
		"uptime":    time.Since(startTime).String(),
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}