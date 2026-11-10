package server_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/amarnathdbg101-coder/GoProjectStructure/server"
)

func TestServer_Initialization(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	srv := server.New(mux, server.Config{
		Port: "9999",
	})

	if srv.Addr() != ":9999" {
		t.Fatalf("expected addr :9999, got %s", srv.Addr())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Shutdown should return nil when not yet running or cleanly stopped
	_ = srv.Shutdown(ctx)
}
