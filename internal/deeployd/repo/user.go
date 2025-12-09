package repo

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

type User struct {
	ID        string    `db:"id"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type UserDTO struct {
	ID    string
	Email string
}

func (u *User) ToDTO() *UserDTO {
	return &UserDTO{
		ID:    u.ID,
		Email: u.Email,
	}
}

type UserRepoInterface interface {
	CountUsers() (int, error)
	CreateUser(user *User) error
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id string) (*User, error)
}

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) CountUsers() (int, error) {
	var count int
	query := `SELECT count(*) FROM users`

	err := r.db.Get(&count, query)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *UserRepo) CreateUser(user *User) error {
	query := `INSERT INTO users (id, email, password) VALUES ($1, $2, $3)`

	_, err := r.db.Exec(query, user.ID, user.Email, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepo) GetUserByEmail(email string) (*User, error) {
	user := &User{}
	query := `SELECT id, email, password, created_at, updated_at FROM users WHERE email = $1`

	err := r.db.Get(user, query, email)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepo) GetUserByID(id string) (*User, error) {
	user := &User{}
	query := `SELECT id, email, password, created_at, updated_at FROM users WHERE id = $1`

	err := r.db.Get(user, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}
