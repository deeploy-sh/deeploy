package repo

import (
	"database/sql"
	"errors"

	"github.com/deeploy-sh/deeploy/internal/shared/model"
	"github.com/jmoiron/sqlx"
)

type GitTokenRepoInterface interface {
	Create(token *model.GitToken) error
	GitToken(id string) (*model.GitToken, error)
	GitTokensByUser(userID string) ([]model.GitToken, error)
	Update(token model.GitToken) error
	Delete(id string) error
}

type GitTokenRepo struct {
	db *sqlx.DB
}

func NewGitTokenRepo(db *sqlx.DB) *GitTokenRepo {
	return &GitTokenRepo{db: db}
}

func (r *GitTokenRepo) Create(token *model.GitToken) error {
	query := `INSERT INTO git_tokens (id, user_id, name, provider, token) VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.Exec(query, token.ID, token.UserID, token.Name, token.Provider, token.Token)
	if err != nil {
		return err
	}

	return nil
}

func (r *GitTokenRepo) GitToken(id string) (*model.GitToken, error) {
	token := &model.GitToken{}
	query := `SELECT id, user_id, name, provider, token, created_at, updated_at FROM git_tokens WHERE id = $1`

	err := r.db.Get(token, query, id)
	if err == sql.ErrNoRows {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (r *GitTokenRepo) GitTokensByUser(userID string) ([]model.GitToken, error) {
	tokens := []model.GitToken{}
	query := `SELECT id, user_id, name, provider, token, created_at, updated_at FROM git_tokens WHERE user_id = $1`

	err := r.db.Select(&tokens, query, userID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (r *GitTokenRepo) Update(token model.GitToken) error {
	query := `UPDATE git_tokens SET name = $1, provider = $2, token = $3 WHERE id = $4`

	result, err := r.db.Exec(query, token.Name, token.Provider, token.Token, token.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("git token not found")
	}

	return nil
}

func (r *GitTokenRepo) Delete(id string) error {
	query := `DELETE FROM git_tokens WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("git token not found")
	}

	return nil
}
