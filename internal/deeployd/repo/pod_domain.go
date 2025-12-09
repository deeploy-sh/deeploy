package repo

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

type PodDomain struct {
	ID         string    `json:"id" db:"id"`
	PodID      string    `json:"pod_id" db:"pod_id"`
	Domain     string    `json:"domain" db:"domain"`
	IsPrimary  bool      `json:"is_primary" db:"is_primary"`
	SSLEnabled bool      `json:"ssl_enabled" db:"ssl_enabled"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type PodDomainRepoInterface interface {
	Create(domain *PodDomain) error
	Domain(id string) (*PodDomain, error)
	DomainByName(domain string) (*PodDomain, error)
	DomainsByPod(podID string) ([]PodDomain, error)
	Update(domain PodDomain) error
	SetPrimary(id string, podID string) error
	Delete(id string) error
	DeleteByPod(podID string) error
}

type PodDomainRepo struct {
	db *sqlx.DB
}

func NewPodDomainRepo(db *sqlx.DB) *PodDomainRepo {
	return &PodDomainRepo{db: db}
}

func (r *PodDomainRepo) Create(domain *PodDomain) error {
	query := `INSERT INTO pod_domains (id, pod_id, domain, is_primary, ssl_enabled) VALUES (?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query, domain.ID, domain.PodID, domain.Domain, domain.IsPrimary, domain.SSLEnabled)
	if err != nil {
		return err
	}

	return nil
}

func (r *PodDomainRepo) Domain(id string) (*PodDomain, error) {
	domain := &PodDomain{}
	query := `SELECT id, pod_id, domain, is_primary, ssl_enabled, created_at, updated_at FROM pod_domains WHERE id = ?`

	err := r.db.Get(domain, query, id)
	if err == sql.ErrNoRows {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	return domain, nil
}

func (r *PodDomainRepo) DomainByName(domainName string) (*PodDomain, error) {
	domain := &PodDomain{}
	query := `SELECT id, pod_id, domain, is_primary, ssl_enabled, created_at, updated_at FROM pod_domains WHERE domain = ?`

	err := r.db.Get(domain, query, domainName)
	if err == sql.ErrNoRows {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	return domain, nil
}

func (r *PodDomainRepo) DomainsByPod(podID string) ([]PodDomain, error) {
	domains := []PodDomain{}
	query := `SELECT id, pod_id, domain, is_primary, ssl_enabled, created_at, updated_at FROM pod_domains WHERE pod_id = ?`

	err := r.db.Select(&domains, query, podID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return domains, nil
}

func (r *PodDomainRepo) Update(domain PodDomain) error {
	query := `UPDATE pod_domains SET domain = ?, ssl_enabled = ? WHERE id = ?`

	result, err := r.db.Exec(query, domain.Domain, domain.SSLEnabled, domain.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("domain not found")
	}

	return nil
}

func (r *PodDomainRepo) SetPrimary(id string, podID string) error {
	// First, unset all other primaries for this pod
	_, err := r.db.Exec(`UPDATE pod_domains SET is_primary = 0 WHERE pod_id = ?`, podID)
	if err != nil {
		return err
	}

	// Then set the new primary
	result, err := r.db.Exec(`UPDATE pod_domains SET is_primary = 1 WHERE id = ?`, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("domain not found")
	}

	return nil
}

func (r *PodDomainRepo) Delete(id string) error {
	query := `DELETE FROM pod_domains WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("domain not found")
	}

	return nil
}

func (r *PodDomainRepo) DeleteByPod(podID string) error {
	query := `DELETE FROM pod_domains WHERE pod_id = ?`

	_, err := r.db.Exec(query, podID)
	if err != nil {
		return err
	}

	return nil
}
