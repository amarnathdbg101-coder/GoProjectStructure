package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/amarnathdbg101-coder/GoProjectStructure/pkg/config"
	"github.com/amarnathdbg101-coder/GoProjectStructure/pkg/database"
	"github.com/amarnathdbg101-coder/GoProjectStructure/internal/handler"
	"github.com/amarnathdbg101-coder/GoProjectStructure/pkg/middleware"
	"github.com/amarnathdbg101-coder/GoProjectStructure/internal/worker"

	"go.uber.org/zap"
)

func main() {
	logger := middleware.New()
	defer logger.Sync()

	cfg := config.MustLoad()

	// 1. Connect to PostgreSQL database connection pool
	db, err := database.ConnectDB(cfg.DBURL)
	if err != nil {
		logger.Warn("database connection failed (running in offline/disconnected mode)", zap.Error(err))
	} else {
		defer db.Close()

		// Run automated schema checks and database migrations
		database.EnsureSchemaExtensions(db)
		if err := database.RunAutoMigrations(cfg.DBURL); err != nil {
			logger.Warn("auto-migration notice", zap.Error(err))
		} else {
			logger.Info("database auto-migrations verified successfully")
		}
	}

	// 2. Setup routes and middleware
	router := handler.RouteSetup(db, logger, cfg)

	// 3. Application context for background concurrent workers
	appCtx, stopApp := context.WithCancel(context.Background())
	defer stopApp()

	if db != nil {
		go worker.StartSampleBackgroundWorker(appCtx, db, logger)
	}

	// 4. Configure Production HTTP Server with Slowloris-safe timeouts
	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1MB
	}

	// 5. Start Server in background goroutine
	go func() {
		logger.Info("server starting", zap.String("addr", server.Addr))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("server failed to start", zap.Error(err))
		}
	}()

	// 6. Graceful Shutdown on OS interrupt signals (SIGINT, SIGTERM)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	logger.Info("shutting down server gracefully...")
	stopApp() // Signal background workers to stop

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("server forced to shutdown", zap.Error(err))
	}

	logger.Info("server exited cleanly")
}
