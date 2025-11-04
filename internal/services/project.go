package services

import (
	"github.com/axadrn/deeploy/internal/data"
)

type ProjectServiceInterface interface {
	Create(project *data.Project) (*data.Project, error)
	Project(id string) (*data.Project, error)
	ProjectsByUser(id string) ([]data.Project, error)
	Update(project data.Project) error
	Delete(id string) error
}

type ProjectService struct {
	repo data.ProjectRepoInterface
}

func NewProjectService(repo *data.ProjectRepo) *ProjectService {
	return &ProjectService{repo: repo}
}

func (s *ProjectService) Create(project *data.Project) (*data.Project, error) {
	err := s.repo.Create(project)
	if err != nil {
		return nil, err
	}
	return project, nil
}

func (s *ProjectService) Project(id string) (*data.Project, error) {
	project, err := s.repo.Project(id)
	if err != nil {
		return nil, err
	}
	return project, nil
}

func (s *ProjectService) ProjectsByUser(id string) ([]data.Project, error) {
	projects, err := s.repo.ProjectsByUser(id)
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (s *ProjectService) Update(project data.Project) error {
	err := s.repo.Update(project)
	if err != nil {
		return err
	}
	return nil
}

func (s *ProjectService) Delete(id string) error {
	err := s.repo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
