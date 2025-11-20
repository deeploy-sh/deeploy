package repo

import (
	"database/sql"
	"time"
)

type User struct {
	ID        string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
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
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) CountUsers() (int, error) {
	var count int

	query := `
		SELECT count(*)
		from users`

	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (m *UserRepo) CreateUser(user *User) error {
	query := `
		INSERT INTO users (id, email, password)
		VALUES(?, ?, ?)`

	_, err := m.db.Exec(query, user.ID, user.Email, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepo) GetUserByEmail(email string) (*User, error) {
	user := &User{}

	query := `
		SELECT id, email, password, created_at, updated_at 
		FROM users
		WHERE email = ?`

	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil // INFO: Like user not found
	}
	if err != nil {
		return nil, err // INFO: real db errors
	}
	return user, nil
}

func (r *UserRepo) GetUserByID(id string) (*User, error) {
	user := &User{}

	query := `
		SELECT id, email, password, created_at, updated_at 
		FROM users
		WHERE id = ?`

	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil // INFO: Like user not found
	}
	if err != nil {
		return nil, err // INFO: real db errors
	}
	return user, nil
}
