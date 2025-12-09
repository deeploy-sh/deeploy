package service

import (
	"github.com/deeploy-sh/deeploy/internal/deeployd/repo"
)

type PodDomainServiceInterface interface {
	Create(domain *repo.PodDomain) (*repo.PodDomain, error)
	Domain(id string) (*repo.PodDomain, error)
	DomainByName(domain string) (*repo.PodDomain, error)
	DomainsByPod(podID string) ([]repo.PodDomain, error)
	Update(domain repo.PodDomain) error
	Delete(id string) error
	DeleteByPod(podID string) error
}

type PodDomainService struct {
	repo repo.PodDomainRepoInterface
}

func NewPodDomainService(repo *repo.PodDomainRepo) *PodDomainService {
	return &PodDomainService{repo: repo}
}

func (s *PodDomainService) Create(domain *repo.PodDomain) (*repo.PodDomain, error) {
	err := s.repo.Create(domain)
	if err != nil {
		return nil, err
	}
	return domain, nil
}

func (s *PodDomainService) Domain(id string) (*repo.PodDomain, error) {
	domain, err := s.repo.Domain(id)
	if err != nil {
		return nil, err
	}
	return domain, nil
}

func (s *PodDomainService) DomainByName(domainName string) (*repo.PodDomain, error) {
	domain, err := s.repo.DomainByName(domainName)
	if err != nil {
		return nil, err
	}
	return domain, nil
}

func (s *PodDomainService) DomainsByPod(podID string) ([]repo.PodDomain, error) {
	domains, err := s.repo.DomainsByPod(podID)
	if err != nil {
		return nil, err
	}
	return domains, nil
}

func (s *PodDomainService) Update(domain repo.PodDomain) error {
	err := s.repo.Update(domain)
	if err != nil {
		return err
	}
	return nil
}

func (s *PodDomainService) Delete(id string) error {
	err := s.repo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (s *PodDomainService) DeleteByPod(podID string) error {
	err := s.repo.DeleteByPod(podID)
	if err != nil {
		return err
	}
	return nil
}
