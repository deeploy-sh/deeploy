package db

import (
	"embed"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var migrations embed.FS

// dialectMap maps database drivers to Goose dialect names
var dialectMap = map[string]string{
	"sqlite": "sqlite3",
	"pgx":    "postgres",
}

// getDialect returns the Goose dialect for the given driver
func getDialect(driver string) string {
	dialect, ok := dialectMap[driver]
	if ok {
		return dialect
	}
	return driver
}

func Init(driver, connection string) (*sqlx.DB, error) {
	// SQLite: create data directory if needed
	if driver == "sqlite" {
		// Extract path before query params (e.g., "./data/deeploy.db?_pragma=...")
		dbPath := connection
		if idx := strings.Index(connection, "?"); idx != -1 {
			dbPath = connection[:idx]
		}
		dir := filepath.Dir(dbPath)
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return nil, fmt.Errorf("failed to create data directory: %w", err)
		}
	}

	db, err := sqlx.Connect(driver, connection)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	// Connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	slog.Info("database connected", "driver", driver)

	// Run migrations
	err = runMigrations(db, driver)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func runMigrations(db *sqlx.DB, driver string) error {
	// Set dialect based on driver
	err := goose.SetDialect(getDialect(driver))
	if err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	// Get migrations subdirectory from embed.FS
	migrationsDir, err := fs.Sub(migrations, "migrations")
	if err != nil {
		return fmt.Errorf("failed to get migrations directory: %w", err)
	}

	goose.SetBaseFS(migrationsDir)

	err = goose.Up(db.DB, ".")
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	slog.Info("migrations completed")
	return nil
}
