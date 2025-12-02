package db

import (
	"embed"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var migrations embed.FS

func Init(connection string) (*sqlx.DB, error) {
	// Create data directory if it doesn't exist
	dir := filepath.Dir(connection)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return nil, err
	}

	// Open database connection
	db, err := sqlx.Connect("sqlite", connection)
	if err != nil {
		return nil, err
	}

	// Connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// SQLite pragmas
	_, err = db.Exec("PRAGMA foreign_keys=ON")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("PRAGMA journal_mode=WAL")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("PRAGMA busy_timeout=5000")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("PRAGMA synchronous=NORMAL")
	if err != nil {
		return nil, err
	}

	// Run migrations
	err = runMigrations(connection)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func runMigrations(connection string) error {
	source, err := iofs.New(migrations, "migrations")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithSourceInstance("iofs", source, "sqlite://"+connection)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	slog.Info("migrations completed")
	return nil
}
