package model

import "time"

type GitToken struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id,omitempty" db:"user_id"`
	Name      string    `json:"name" db:"name"`
	Provider  string    `json:"provider" db:"provider"`
	Token     string    `json:"-" db:"token"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
