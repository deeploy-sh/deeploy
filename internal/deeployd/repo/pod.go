package repo

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

type Pod struct {
	ID             string    `json:"id" db:"id"`
	UserID         string    `json:"user_id" db:"user_id"`
	ProjectID      string    `json:"project_id" db:"project_id"`
	Title          string    `json:"title" db:"title"`
	Description    string    `json:"description" db:"description"`
	RepoURL        *string   `json:"repo_url" db:"repo_url"`
	Branch         string    `json:"branch" db:"branch"`
	DockerfilePath string    `json:"dockerfile_path" db:"dockerfile_path"`
	ContainerID    *string   `json:"container_id" db:"container_id"`
	Status         string    `json:"status" db:"status"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

type PodRepoInterface interface {
	Create(pod *Pod) error
	Pod(id string) (*Pod, error)
	PodsByProject(id string) ([]Pod, error)
	PodsByUser(id string) ([]Pod, error)
	Update(pod Pod) error
	Delete(id string) error
}

type PodRepo struct {
	db *sqlx.DB
}

func NewPodRepo(db *sqlx.DB) *PodRepo {
	return &PodRepo{db: db}
}

func (r *PodRepo) Create(pod *Pod) error {
	query := `INSERT INTO pods (id, user_id, project_id, title, description, repo_url, branch, dockerfile_path, status) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query, pod.ID, pod.UserID, pod.ProjectID, pod.Title, pod.Description, pod.RepoURL, pod.Branch, pod.DockerfilePath, pod.Status)
	if err != nil {
		return err
	}

	return nil
}

func (r *PodRepo) Pod(id string) (*Pod, error) {
	pod := &Pod{}
	query := `SELECT id, user_id, project_id, title, description, repo_url, branch, dockerfile_path, container_id, status, created_at, updated_at FROM pods WHERE id = ?`

	err := r.db.Get(pod, query, id)
	if err == sql.ErrNoRows {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	return pod, nil
}

func (r *PodRepo) PodsByProject(id string) ([]Pod, error) {
	pods := []Pod{}
	query := `SELECT id, user_id, project_id, title, description, repo_url, branch, dockerfile_path, container_id, status, created_at, updated_at FROM pods WHERE project_id = ?`

	err := r.db.Select(&pods, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return pods, nil
}

func (r *PodRepo) PodsByUser(id string) ([]Pod, error) {
	pods := []Pod{}
	query := `SELECT id, user_id, project_id, title, description, repo_url, branch, dockerfile_path, container_id, status, created_at, updated_at FROM pods WHERE user_id = ?`

	err := r.db.Select(&pods, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return pods, nil
}

func (r *PodRepo) Update(pod Pod) error {
	query := `UPDATE pods SET title = ?, description = ? WHERE id = ?`

	result, err := r.db.Exec(query, pod.Title, pod.Description, pod.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("Pod not found")
	}

	return nil
}

func (r *PodRepo) Delete(id string) error {
	query := `DELETE FROM pods WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("Pod not found")
	}

	return nil
}
