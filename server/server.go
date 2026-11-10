// Package server provides a minimal, production-grade HTTP server lifecycle manager with graceful shutdown.
package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Config holds basic server operational settings.
type Config struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// DefaultConfig returns production defaults for timeouts and port.
func DefaultConfig() Config {
	return Config{
		Port:         "8080",
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

// Server wraps http.Server with signal trapping and graceful shutdown.
type Server struct {
	httpServer *http.Server
	cfg        Config
}

// New creates an initialized Server instance wrapping any http.Handler.
func New(handler http.Handler, cfg ...Config) *Server {
	c := DefaultConfig()
	if len(cfg) > 0 {
		if cfg[0].Port != "" {
			c.Port = cfg[0].Port
		}
		if cfg[0].ReadTimeout > 0 {
			c.ReadTimeout = cfg[0].ReadTimeout
		}
		if cfg[0].WriteTimeout > 0 {
			c.WriteTimeout = cfg[0].WriteTimeout
		}
		if cfg[0].IdleTimeout > 0 {
			c.IdleTimeout = cfg[0].IdleTimeout
		}
	}

	return &Server{
		cfg: c,
		httpServer: &http.Server{
			Addr:         ":" + c.Port,
			Handler:      handler,
			ReadTimeout:  c.ReadTimeout,
			WriteTimeout: c.WriteTimeout,
			IdleTimeout:  c.IdleTimeout,
		},
	}
}

// Addr returns the configured listening address.
func (s *Server) Addr() string {
	return s.httpServer.Addr
}

// Start begins listening and blocks until an OS interrupt signal (SIGINT/SIGTERM) is received,
// shutting down gracefully with a 10-second timeout.
func (s *Server) Start() error {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	serverErr := make(chan error, 1)
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	select {
	case err := <-serverErr:
		return fmt.Errorf("server: %w", err)
	case <-shutdown:
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return s.httpServer.Shutdown(ctx)
	}
}

// Shutdown gracefully terminates the server using the provided context.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
