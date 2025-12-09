package repo

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

type PodEnvVar struct {
	ID        string    `json:"id" db:"id"`
	PodID     string    `json:"pod_id" db:"pod_id"`
	Key       string    `json:"key" db:"key"`
	Value     string    `json:"value" db:"value"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type PodEnvVarRepoInterface interface {
	Create(envVar *PodEnvVar) error
	EnvVar(id string) (*PodEnvVar, error)
	EnvVarsByPod(podID string) ([]PodEnvVar, error)
	Update(envVar PodEnvVar) error
	Delete(id string) error
	DeleteByPod(podID string) error
}

type PodEnvVarRepo struct {
	db *sqlx.DB
}

func NewPodEnvVarRepo(db *sqlx.DB) *PodEnvVarRepo {
	return &PodEnvVarRepo{db: db}
}

func (r *PodEnvVarRepo) Create(envVar *PodEnvVar) error {
	query := `INSERT INTO pod_env_vars (id, pod_id, key, value) VALUES (?, ?, ?, ?)`

	_, err := r.db.Exec(query, envVar.ID, envVar.PodID, envVar.Key, envVar.Value)
	if err != nil {
		return err
	}

	return nil
}

func (r *PodEnvVarRepo) EnvVar(id string) (*PodEnvVar, error) {
	envVar := &PodEnvVar{}
	query := `SELECT id, pod_id, key, value, created_at, updated_at FROM pod_env_vars WHERE id = ?`

	err := r.db.Get(envVar, query, id)
	if err == sql.ErrNoRows {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	return envVar, nil
}

func (r *PodEnvVarRepo) EnvVarsByPod(podID string) ([]PodEnvVar, error) {
	envVars := []PodEnvVar{}
	query := `SELECT id, pod_id, key, value, created_at, updated_at FROM pod_env_vars WHERE pod_id = ?`

	err := r.db.Select(&envVars, query, podID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return envVars, nil
}

func (r *PodEnvVarRepo) Update(envVar PodEnvVar) error {
	query := `UPDATE pod_env_vars SET key = ?, value = ? WHERE id = ?`

	result, err := r.db.Exec(query, envVar.Key, envVar.Value, envVar.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("env var not found")
	}

	return nil
}

func (r *PodEnvVarRepo) Delete(id string) error {
	query := `DELETE FROM pod_env_vars WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("env var not found")
	}

	return nil
}

func (r *PodEnvVarRepo) DeleteByPod(podID string) error {
	query := `DELETE FROM pod_env_vars WHERE pod_id = ?`

	_, err := r.db.Exec(query, podID)
	if err != nil {
		return err
	}

	return nil
}
