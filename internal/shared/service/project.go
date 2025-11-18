package service

import (
	"github.com/deeploy-sh/deeploy/internal/shared/repo"
)

type ProjectServiceInterface interface {
	Create(project *repo.Project) (*repo.Project, error)
	Project(id string) (*repo.Project, error)
	ProjectsByUser(id string) ([]repo.Project, error)
	Update(project repo.Project) error
	Delete(id string) error
}

type ProjectService struct {
	repo repo.ProjectRepoInterface
}

func NewProjectService(repo *repo.ProjectRepo) *ProjectService {
	return &ProjectService{repo: repo}
}

func (s *ProjectService) Create(project *repo.Project) (*repo.Project, error) {
	err := s.repo.Create(project)
	if err != nil {
		return nil, err
	}
	return project, nil
}

func (s *ProjectService) Project(id string) (*repo.Project, error) {
	project, err := s.repo.Project(id)
	if err != nil {
		return nil, err
	}
	return project, nil
}

func (s *ProjectService) ProjectsByUser(id string) ([]repo.Project, error) {
	projects, err := s.repo.ProjectsByUser(id)
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (s *ProjectService) Update(project repo.Project) error {
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
