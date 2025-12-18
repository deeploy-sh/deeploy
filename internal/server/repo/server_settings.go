package repo

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type ServerSettingsRepo struct {
	db *sqlx.DB
}

func NewServerSettingsRepo(db *sqlx.DB) *ServerSettingsRepo {
	return &ServerSettingsRepo{db: db}
}

// Get returns a setting value by key. Returns empty string if not found.
func (r *ServerSettingsRepo) Get(key string) (string, error) {
	var value string
	query := `SELECT value FROM server_settings WHERE key = $1`

	err := r.db.Get(&value, query, key)
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	return value, nil
}

// Set creates or updates a setting.
func (r *ServerSettingsRepo) Set(key, value string) error {
	query := `
		INSERT INTO server_settings (key, value, updated_at)
		VALUES ($1, $2, CURRENT_TIMESTAMP)
		ON CONFLICT (key) DO UPDATE SET value = $2, updated_at = CURRENT_TIMESTAMP
	`
	_, err := r.db.Exec(query, key, value)
	return err
}

// Delete removes a setting by key.
func (r *ServerSettingsRepo) Delete(key string) error {
	query := `DELETE FROM server_settings WHERE key = $1`
	_, err := r.db.Exec(query, key)
	return err
}
