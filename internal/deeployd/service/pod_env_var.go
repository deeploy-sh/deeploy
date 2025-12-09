package service

import (
	"github.com/deeploy-sh/deeploy/internal/deeployd/repo"
)

type PodEnvVarServiceInterface interface {
	Create(envVar *repo.PodEnvVar) (*repo.PodEnvVar, error)
	EnvVar(id string) (*repo.PodEnvVar, error)
	EnvVarsByPod(podID string) ([]repo.PodEnvVar, error)
	Update(envVar repo.PodEnvVar) error
	Delete(id string) error
	DeleteByPod(podID string) error
}

type PodEnvVarService struct {
	repo repo.PodEnvVarRepoInterface
}

func NewPodEnvVarService(repo *repo.PodEnvVarRepo) *PodEnvVarService {
	return &PodEnvVarService{repo: repo}
}

func (s *PodEnvVarService) Create(envVar *repo.PodEnvVar) (*repo.PodEnvVar, error) {
	err := s.repo.Create(envVar)
	if err != nil {
		return nil, err
	}
	return envVar, nil
}

func (s *PodEnvVarService) EnvVar(id string) (*repo.PodEnvVar, error) {
	envVar, err := s.repo.EnvVar(id)
	if err != nil {
		return nil, err
	}
	return envVar, nil
}

func (s *PodEnvVarService) EnvVarsByPod(podID string) ([]repo.PodEnvVar, error) {
	envVars, err := s.repo.EnvVarsByPod(podID)
	if err != nil {
		return nil, err
	}
	return envVars, nil
}

func (s *PodEnvVarService) Update(envVar repo.PodEnvVar) error {
	err := s.repo.Update(envVar)
	if err != nil {
		return err
	}
	return nil
}

func (s *PodEnvVarService) Delete(id string) error {
	err := s.repo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (s *PodEnvVarService) DeleteByPod(podID string) error {
	err := s.repo.DeleteByPod(podID)
	if err != nil {
		return err
	}
	return nil
}
