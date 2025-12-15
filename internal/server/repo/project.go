package repo

import (
	"database/sql"
	"errors"

	"github.com/deeploy-sh/deeploy/internal/shared/model"
	"github.com/jmoiron/sqlx"
)

type ProjectRepoInterface interface {
	Create(project *model.Project) error
	Project(id string) (*model.Project, error)
	ProjectsByUser(id string) ([]model.Project, error)
	Update(project model.Project) error
	Delete(id string) error
}

type ProjectRepo struct {
	db *sqlx.DB
}

func NewProjectRepo(db *sqlx.DB) *ProjectRepo {
	return &ProjectRepo{db: db}
}

func (r *ProjectRepo) Create(project *model.Project) error {
	query := `INSERT INTO projects (id, user_id, title) VALUES ($1, $2, $3)`

	_, err := r.db.Exec(query, project.ID, project.UserID, project.Title)
	if err != nil {
		return err
	}

	return nil
}

func (r *ProjectRepo) Project(id string) (*model.Project, error) {
	project := &model.Project{}
	query := `SELECT id, user_id, title, created_at, updated_at FROM projects WHERE id = $1`

	err := r.db.Get(project, query, id)
	if err == sql.ErrNoRows {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (r *ProjectRepo) ProjectsByUser(id string) ([]model.Project, error) {
	projects := []model.Project{}
	query := `SELECT id, user_id, title, created_at, updated_at FROM projects WHERE user_id = $1`

	err := r.db.Select(&projects, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (r *ProjectRepo) Update(project model.Project) error {
	query := `UPDATE projects SET title = $1 WHERE id = $2`

	result, err := r.db.Exec(query, project.Title, project.ID)
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
