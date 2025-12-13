package repo

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

type Project struct {
	ID          string    `json:"id" db:"id"`
	UserID      string    `json:"user_id" db:"user_id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type ProjectRepoInterface interface {
	Create(project *Project) error
	Project(id string) (*Project, error)
	ProjectsByUser(id string) ([]Project, error)
	Update(project Project) error
	Delete(id string) error
}

type ProjectRepo struct {
	db *sqlx.DB
}

func NewProjectRepo(db *sqlx.DB) *ProjectRepo {
	return &ProjectRepo{db: db}
}

func (r *ProjectRepo) Create(project *Project) error {
	query := `INSERT INTO projects (id, user_id, title, description) VALUES ($1, $2, $3, $4)`

	_, err := r.db.Exec(query, project.ID, project.UserID, project.Title, project.Description)
	if err != nil {
		return err
	}

	return nil
}

func (r *ProjectRepo) Project(id string) (*Project, error) {
	project := &Project{}
	query := `SELECT id, user_id, title, description, created_at, updated_at FROM projects WHERE id = $1`

	err := r.db.Get(project, query, id)
	if err == sql.ErrNoRows {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (r *ProjectRepo) ProjectsByUser(id string) ([]Project, error) {
	projects := []Project{}
	query := `SELECT id, user_id, title, description, created_at, updated_at FROM projects WHERE user_id = $1`

	err := r.db.Select(&projects, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (r *ProjectRepo) Update(project Project) error {
	query := `UPDATE projects SET title = $1, description = $2 WHERE id = $3`

	result, err := r.db.Exec(query, project.Title, project.Description, project.ID)
	if err != nil {
		return err
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
	query := `DELETE FROM projects WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
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
