package repo

import (
	"database/sql"
	"errors"
	"time"
)

type Pod struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	ProjectID   string    `json:"project_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PodDTO struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	ProjectID   string `json:"project_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (u *Pod) ToDTO() *PodDTO {
	return &PodDTO{
		ID:          u.ID,
		UserID:      u.UserID,
		Title:       u.Title,
		Description: u.Description,
	}
}

type PodRepoInterface interface {
	Create(pod *Pod) error
	Pod(id string) (*Pod, error)
	PodsByProject(id string) ([]Pod, error)
	Update(pod Pod) error
	Delete(id string) error
}

type PodRepo struct {
	db *sql.DB
}

func NewPodRepo(db *sql.DB) *PodRepo {
	return &PodRepo{db: db}
}

func (m *PodRepo) Create(pod *Pod) error {
	query := `
		INSERT INTO pods (id, user_id, project_id, title, description)
		VALUES(?, ?, ?, ?, ?)`

	_, err := m.db.Exec(
		query,
		pod.ID,
		pod.UserID,
		pod.ProjectID,
		pod.Title,
		pod.Description,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *PodRepo) Pod(id string) (*Pod, error) {
	pod := &Pod{}

	query := `
		SELECT id, user_id, project_id, title, description, created_at, updated_at 
		FROM pods
		WHERE id = ?`

	err := r.db.QueryRow(query, id).Scan(
		&pod.ID,
		&pod.UserID,
		&pod.ProjectID,
		&pod.Title,
		&pod.Description,
		&pod.CreatedAt,
		&pod.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, err // INFO: Like pod not found
	}
	if err != nil {
		return nil, err // INFO: real db errors
	}
	return pod, nil
}

func (r *PodRepo) PodsByProject(id string) ([]Pod, error) {
	pods := []Pod{}

	query := `
		SELECT id, user_id, project_id, title, description, created_at, updated_at 
		FROM pods
		WHERE project_id = ?`

	rows, err := r.db.Query(query, id)
	if err == sql.ErrNoRows {
		return nil, nil // INFO: Like pod not found
	}
	if err != nil {
		return nil, err // INFO: real db errors
	}
	defer rows.Close()

	for rows.Next() {
		p := &Pod{}
		err := rows.Scan(
			&p.ID,
			&p.UserID,
			&p.ProjectID,
			&p.Title,
			&p.Description,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		pods = append(pods, *p)
	}
	return pods, nil
}

func (r *PodRepo) Update(pod Pod) error {
	query := `
		UPDATE pods
		SET title = ?, description = ?
		WHERE id = ?`

	result, err := r.db.Exec(query, pod.Title, pod.Description, pod.ID)
	if err != nil {
		return err // INFO: real db errors
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
	query := `
		DELETE FROM pods
		WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err // INFO: real db errors
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
