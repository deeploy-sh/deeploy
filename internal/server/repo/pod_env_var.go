package repo

import (
	"database/sql"
	"errors"

	"github.com/deeploy-sh/deeploy/internal/shared/model"
	"github.com/jmoiron/sqlx"
)

type PodEnvVarRepoInterface interface {
	Create(envVar *model.PodEnvVar) error
	EnvVar(id string) (*model.PodEnvVar, error)
	EnvVarsByPod(podID string) ([]model.PodEnvVar, error)
	Update(envVar model.PodEnvVar) error
	Delete(id string) error
	DeleteByPod(podID string) error
}

type PodEnvVarRepo struct {
	db *sqlx.DB
}

func NewPodEnvVarRepo(db *sqlx.DB) *PodEnvVarRepo {
	return &PodEnvVarRepo{db: db}
}

func (r *PodEnvVarRepo) Create(envVar *model.PodEnvVar) error {
	query := `INSERT INTO pod_env_vars (id, pod_id, key, value) VALUES ($1, $2, $3, $4)`

	_, err := r.db.Exec(query, envVar.ID, envVar.PodID, envVar.Key, envVar.Value)
	if err != nil {
		return err
	}

	return nil
}

func (r *PodEnvVarRepo) EnvVar(id string) (*model.PodEnvVar, error) {
	envVar := &model.PodEnvVar{}
	query := `SELECT id, pod_id, key, value, created_at, updated_at FROM pod_env_vars WHERE id = $1`

	err := r.db.Get(envVar, query, id)
	if err == sql.ErrNoRows {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	return envVar, nil
}

func (r *PodEnvVarRepo) EnvVarsByPod(podID string) ([]model.PodEnvVar, error) {
	envVars := []model.PodEnvVar{}
	query := `SELECT id, pod_id, key, value, created_at, updated_at FROM pod_env_vars WHERE pod_id = $1`

	err := r.db.Select(&envVars, query, podID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return envVars, nil
}

func (r *PodEnvVarRepo) Update(envVar model.PodEnvVar) error {
	query := `UPDATE pod_env_vars SET key = $1, value = $2 WHERE id = $3`

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
	query := `DELETE FROM pod_env_vars WHERE id = $1`

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
	query := `DELETE FROM pod_env_vars WHERE pod_id = $1`

	_, err := r.db.Exec(query, podID)
	if err != nil {
		return err
	}

	return nil
}
