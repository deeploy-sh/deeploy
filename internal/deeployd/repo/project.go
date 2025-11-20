package repo

import (
	"database/sql"
	"errors"
	"time"
)

type Project struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProjectDTO struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (u *Project) ToDTO() *ProjectDTO {
	return &ProjectDTO{
		ID:          u.ID,
		UserID:      u.UserID,
		Title:       u.Title,
		Description: u.Description,
	}
}

type ProjectRepoInterface interface {
	Create(project *Project) error
	Project(id string) (*Project, error)
	ProjectsByUser(id string) ([]Project, error)
	Update(project Project) error
	Delete(id string) error
}

type ProjectRepo struct {
	db *sql.DB
}

func NewProjectRepo(db *sql.DB) *ProjectRepo {
	return &ProjectRepo{db: db}
}

func (m *ProjectRepo) Create(project *Project) error {
	query := `
		INSERT INTO projects (id, user_id, title, description)
		VALUES(?, ?, ?, ?)`

	_, err := m.db.Exec(query, project.ID, project.UserID, project.Title, project.Description)
	if err != nil {
		return err
	}

	return nil
}

func (r *ProjectRepo) Project(id string) (*Project, error) {
	project := &Project{}

	query := `
		SELECT id, user_id, title, description, created_at, updated_at 
		FROM projects
		WHERE id = ?`

	err := r.db.QueryRow(query, id).Scan(
		&project.ID,
		&project.UserID,
		&project.Title,
		&project.Description,
		&project.CreatedAt,
		&project.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, err // INFO: Like project not found
	}
	if err != nil {
		return nil, err // INFO: real db errors
	}
	return project, nil
}

func (r *ProjectRepo) ProjectsByUser(id string) ([]Project, error) {
	projects := []Project{}

	query := `
		SELECT id, user_id, title, description, created_at, updated_at 
		FROM projects
		WHERE user_id = ?`

	rows, err := r.db.Query(query, id)
	if err == sql.ErrNoRows {
		return nil, nil // INFO: Like project not found
	}
	if err != nil {
		return nil, err // INFO: real db errors
	}
	defer rows.Close()

	for rows.Next() {
		p := &Project{}
		err := rows.Scan(
			&p.ID,
			&p.UserID,
			&p.Title,
			&p.Description,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		projects = append(projects, *p)
	}
	return projects, nil
}

func (r *ProjectRepo) Update(project Project) error {
	query := `
		UPDATE projects
		SET title = ?, description = ?
		WHERE id = ?`

	result, err := r.db.Exec(query, project.Title, project.Description, project.ID)
	if err != nil {
		return err // INFO: real db errors
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("Project not found")
	}

	return nil
}

func (r *ProjectRepo) Delete(id string) error {
	query := `
		DELETE FROM projects
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
		return errors.New("Project not found")
	}

	return nil
}
