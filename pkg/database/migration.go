package database

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

// RunAutoMigrations executes all pending migrations from migrations folder.
func RunAutoMigrations(dbURL string) error {
	m, err := migrate.New("file://internal/migration", dbURL)
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run up migrations: %w", err)
	}

	return nil
}

// EnsureSchemaExtensions runs idempotent SQL statements to ensure required Postgres extensions exist.
func EnsureSchemaExtensions(pool *pgxpool.Pool) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	queries := []string{
		"CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";",
		"CREATE EXTENSION IF NOT EXISTS \"pgcrypto\";",
	}

	for _, q := range queries {
		_, _ = pool.Exec(ctx, q)
	}
}
