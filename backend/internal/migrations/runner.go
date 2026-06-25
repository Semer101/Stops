package migrations

import (
	"context"
	"embed"
	"fmt"
	"path"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed sql/*.up.sql
var migrationFiles embed.FS

type Runner struct {
	db *pgxpool.Pool
}

func NewRunner(db *pgxpool.Pool) *Runner {
	return &Runner{db: db}
}

func (runner *Runner) Up(ctx context.Context) error {
	if err := runner.ensureSchemaMigrations(ctx); err != nil {
		return err
	}

	entries, err := migrationFiles.ReadDir("sql")
	if err != nil {
		return fmt.Errorf("read migration directory: %w", err)
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".up.sql") {
			names = append(names, entry.Name())
		}
	}
	sort.Strings(names)

	for _, name := range names {
		applied, err := runner.isApplied(ctx, name)
		if err != nil {
			return err
		}
		if applied {
			continue
		}

		content, err := migrationFiles.ReadFile(path.Join("sql", name))
		if err != nil {
			return fmt.Errorf("read migration %s: %w", name, err)
		}

		if err := runner.apply(ctx, name, string(content)); err != nil {
			return err
		}
	}

	return nil
}

func (runner *Runner) ensureSchemaMigrations(ctx context.Context) error {
	const query = `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`

	if _, err := runner.db.Exec(ctx, query); err != nil {
		return fmt.Errorf("ensure schema migrations table: %w", err)
	}

	return nil
}

func (runner *Runner) isApplied(ctx context.Context, name string) (bool, error) {
	const query = `SELECT EXISTS (SELECT 1 FROM schema_migrations WHERE version = $1)`

	var exists bool
	if err := runner.db.QueryRow(ctx, query, name).Scan(&exists); err != nil {
		return false, fmt.Errorf("check migration %s: %w", name, err)
	}

	return exists, nil
}

func (runner *Runner) apply(ctx context.Context, name string, sql string) error {
	tx, err := runner.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin migration %s: %w", name, err)
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, sql); err != nil {
		return fmt.Errorf("apply migration %s: %w", name, err)
	}

	if _, err := tx.Exec(ctx, `INSERT INTO schema_migrations (version) VALUES ($1)`, name); err != nil {
		return fmt.Errorf("record migration %s: %w", name, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit migration %s: %w", name, err)
	}

	return nil
}
