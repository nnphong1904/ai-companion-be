package database

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"ai-companion-be/internal/config"
)

// NewPostgresPool creates a connection pool configured for Supabase.
//
// When UsePooler is true (Supabase PgBouncer on port 6543, transaction mode):
//   - Disables prepared statement caching (best practice 2.4)
//   - Uses QueryExecModeExec to avoid extended protocol conflicts
//
// When UsePooler is false (direct connection on port 5432):
//   - Uses standard pgx extended protocol with prepared statements
func NewPostgresPool(ctx context.Context, cfg config.DatabaseConfig) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("parsing database config: %w", err)
	}

	poolConfig.MaxConns = 20
	poolConfig.MinConns = 2

	// Supabase pooler compatibility (transaction-mode PgBouncer).
	// Named prepared statements are tied to individual connections.
	// In transaction mode, connections are shared between requests,
	// causing "prepared statement does not exist" errors.
	// Using QueryExecModeExec sends queries as simple protocol messages.
	if cfg.UsePooler {
		poolConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeExec
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("creating connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("pinging database: %w", err)
	}

	return pool, nil
}

// RunMigrations reads and executes SQL migration files in order.
func RunMigrations(ctx context.Context, pool *pgxpool.Pool, migrationsDir string) error {
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("reading migrations directory: %w", err)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		path := filepath.Join(migrationsDir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading %s: %w", entry.Name(), err)
		}

		if _, err := pool.Exec(ctx, string(data)); err != nil {
			return fmt.Errorf("executing %s: %w", entry.Name(), err)
		}
	}

	return nil
}
