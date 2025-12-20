package repo

import (
	"database/sql"
	"fmt"

	"github.com/deeploy-sh/deeploy/internal/shared/errs"
	"github.com/deeploy-sh/deeploy/internal/shared/model"
	"github.com/jmoiron/sqlx"
)

type PodDomainRepoInterface interface {
	Create(domain *model.PodDomain) error
	Domain(id string) (*model.PodDomain, error)
	DomainByName(domain string) (*model.PodDomain, error)
	DomainsByPod(podID string) ([]model.PodDomain, error)
	Update(domain model.PodDomain) error
	Delete(id string) error
	DeleteByPod(podID string) error
}

type PodDomainRepo struct {
	db *sqlx.DB
}

func NewPodDomainRepo(db *sqlx.DB) *PodDomainRepo {
	return &PodDomainRepo{db: db}
}

func (r *PodDomainRepo) Create(domain *model.PodDomain) error {
	query := `INSERT INTO pod_domains (id, pod_id, domain, type, port, ssl_enabled) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.Exec(query, domain.ID, domain.PodID, domain.Domain, domain.Type, domain.Port, domain.SSLEnabled)
	if err != nil {
		return err
	}

	return nil
}

func (r *PodDomainRepo) Domain(id string) (*model.PodDomain, error) {
	domain := &model.PodDomain{}
	query := `SELECT id, pod_id, domain, type, port, ssl_enabled, created_at, updated_at FROM pod_domains WHERE id = $1`

	err := r.db.Get(domain, query, id)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("domain %s: %w", id, errs.ErrNotFound)
	}
	if err != nil {
		return nil, err
	}

	return domain, nil
}

func (r *PodDomainRepo) DomainByName(domainName string) (*model.PodDomain, error) {
	domain := &model.PodDomain{}
	query := `SELECT id, pod_id, domain, type, port, ssl_enabled, created_at, updated_at FROM pod_domains WHERE domain = $1`

	err := r.db.Get(domain, query, domainName)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("domain %s: %w", domainName, errs.ErrNotFound)
	}
	if err != nil {
		return nil, err
	}

	return domain, nil
}

func (r *PodDomainRepo) DomainsByPod(podID string) ([]model.PodDomain, error) {
	domains := []model.PodDomain{}
	query := `SELECT id, pod_id, domain, type, port, ssl_enabled, created_at, updated_at FROM pod_domains WHERE pod_id = $1`

	err := r.db.Select(&domains, query, podID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return domains, nil
}

func (r *PodDomainRepo) Update(domain model.PodDomain) error {
	query := `UPDATE pod_domains SET domain = $1, type = $2, port = $3, ssl_enabled = $4 WHERE id = $5`

	result, err := r.db.Exec(query, domain.Domain, domain.Type, domain.Port, domain.SSLEnabled, domain.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("domain %s: %w", domain.ID, errs.ErrNotFound)
	}

	return nil
}

func (r *PodDomainRepo) Delete(id string) error {
	query := `DELETE FROM pod_domains WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("domain %s: %w", id, errs.ErrNotFound)
	}

	return nil
}

func (r *PodDomainRepo) DeleteByPod(podID string) error {
	query := `DELETE FROM pod_domains WHERE pod_id = $1`

	_, err := r.db.Exec(query, podID)
	if err != nil {
		return err
	}

	return nil
}
