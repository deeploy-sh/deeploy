package repo

import (
	"database/sql"
	"fmt"

	"github.com/deeploy-sh/deeploy/internal/shared/errs"
	"github.com/deeploy-sh/deeploy/internal/shared/model"
	"github.com/jmoiron/sqlx"
)

type PodRepoInterface interface {
	Create(pod *model.Pod) error
	Pod(id string) (*model.Pod, error)
	PodsByProject(id string) ([]model.Pod, error)
	PodsByUser(id string) ([]model.Pod, error)
	CountByProject(id string) (int, error)
	Update(pod model.Pod) error
	Delete(id string) error
}

type PodRepo struct {
	db *sqlx.DB
}

func NewPodRepo(db *sqlx.DB) *PodRepo {
	return &PodRepo{db: db}
}

func (r *PodRepo) Create(pod *model.Pod) error {
	query := `INSERT INTO pods (id, user_id, project_id, title, repo_url, branch, dockerfile_path, git_token_id, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.db.Exec(query, pod.ID, pod.UserID, pod.ProjectID, pod.Title, pod.RepoURL, pod.Branch, pod.DockerfilePath, pod.GitTokenID, pod.Status)
	if err != nil {
		return err
	}

	return nil
}

func (r *PodRepo) Pod(id string) (*model.Pod, error) {
	pod := &model.Pod{}
	query := `SELECT id, user_id, project_id, title, repo_url, branch, dockerfile_path, git_token_id, container_id, status, created_at, updated_at FROM pods WHERE id = $1`

	err := r.db.Get(pod, query, id)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("pod %s: %w", id, errs.ErrNotFound)
	}
	if err != nil {
		return nil, err
	}

	return pod, nil
}

func (r *PodRepo) PodsByProject(id string) ([]model.Pod, error) {
	pods := []model.Pod{}
	query := `SELECT id, user_id, project_id, title, repo_url, branch, dockerfile_path, git_token_id, container_id, status, created_at, updated_at FROM pods WHERE project_id = $1`

	err := r.db.Select(&pods, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return pods, nil
}

func (r *PodRepo) CountByProject(id string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM pods WHERE project_id = $1`

	err := r.db.Get(&count, query, id)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *PodRepo) PodsByUser(id string) ([]model.Pod, error) {
	pods := []model.Pod{}
	query := `SELECT id, user_id, project_id, title, repo_url, branch, dockerfile_path, git_token_id, container_id, status, created_at, updated_at FROM pods WHERE user_id = $1`

	err := r.db.Select(&pods, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return pods, nil
}

func (r *PodRepo) Update(pod model.Pod) error {
	query := `UPDATE pods SET title = $1, repo_url = $2, branch = $3, dockerfile_path = $4, git_token_id = $5, container_id = $6, status = $7 WHERE id = $8`

	result, err := r.db.Exec(query, pod.Title, pod.RepoURL, pod.Branch, pod.DockerfilePath, pod.GitTokenID, pod.ContainerID, pod.Status, pod.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("pod %s: %w", pod.ID, errs.ErrNotFound)
	}

	return nil
}

func (r *PodRepo) Delete(id string) error {
	query := `DELETE FROM pods WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("pod %s: %w", id, errs.ErrNotFound)
	}

	return nil
}
