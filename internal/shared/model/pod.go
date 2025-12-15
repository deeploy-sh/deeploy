package model

import "time"

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
