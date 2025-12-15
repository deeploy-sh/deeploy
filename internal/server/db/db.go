package db

import (
	"embed"
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrations embed.FS

func Init(databaseURL string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	// Connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Run migrations
	err = runMigrations(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func runMigrations(db *sqlx.DB) error {
	goose.SetBaseFS(migrations)

	err := goose.SetDialect("postgres")
	if err != nil {
		return err
	}

	err = goose.Up(db.DB, "migrations")
	if err != nil {
		return err
	}

	slog.Info("migrations completed")
	return nil
}
