package model

import "time"

type Pod struct {
	ID             string    `json:"id" db:"id"`
	UserID         string    `json:"user_id" db:"user_id"`
	ProjectID      string    `json:"project_id" db:"project_id"`
	Title          string    `json:"title" db:"title"`
	RepoURL        *string   `json:"repo_url" db:"repo_url"`
	Branch         string    `json:"branch" db:"branch"`
	DockerfilePath string    `json:"dockerfile_path" db:"dockerfile_path"`
	GitTokenID     *string   `json:"git_token_id" db:"git_token_id"`
	ContainerID    *string   `json:"container_id" db:"container_id"`
	Status         string    `json:"status" db:"status"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

type PodEnvVar struct {
	ID        string    `json:"id" db:"id"`
	PodID     string    `json:"pod_id" db:"pod_id"`
	Key       string    `json:"key" db:"key"`
	Value     string    `json:"value" db:"value"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type PodDomain struct {
	ID         string    `json:"id" db:"id"`
	PodID      string    `json:"pod_id" db:"pod_id"`
	Domain     string    `json:"domain" db:"domain"`
	Type       string    `json:"type" db:"type"`
	Port       int       `json:"port" db:"port"`
	SSLEnabled bool      `json:"ssl_enabled" db:"ssl_enabled"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}
