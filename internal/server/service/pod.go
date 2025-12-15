package service

import (
	"github.com/deeploy-sh/deeploy/internal/server/repo"
	"github.com/deeploy-sh/deeploy/internal/shared/model"
)

type PodServiceInterface interface {
	Create(pod *model.Pod) (*model.Pod, error)
	Pod(id string) (*model.Pod, error)
	PodsByProject(id string) ([]model.Pod, error)
	PodsByUser(id string) ([]model.Pod, error)
	CountByProject(id string) (int, error)
	Update(pod model.Pod) error
	Delete(id string) error
}

type PodService struct {
	repo repo.PodRepoInterface
}

func NewPodService(repo *repo.PodRepo) *PodService {
	return &PodService{repo: repo}
}

func (s *PodService) Create(pod *model.Pod) (*model.Pod, error) {
	err := s.repo.Create(pod)
	if err != nil {
		return nil, err
	}
	return pod, nil
}

func (s *PodService) Pod(id string) (*model.Pod, error) {
	pod, err := s.repo.Pod(id)
	if err != nil {
		return nil, err
	}
	return pod, nil
}

func (s *PodService) PodsByProject(id string) ([]model.Pod, error) {
	pods, err := s.repo.PodsByProject(id)
	if err != nil {
		return nil, err
	}
	return pods, nil
}

func (s *PodService) PodsByUser(id string) ([]model.Pod, error) {
	pods, err := s.repo.PodsByUser(id)
	if err != nil {
		return nil, err
	}

	return pods, nil
}

func (s *PodService) CountByProject(id string) (int, error) {
	return s.repo.CountByProject(id)
}

func (s *PodService) Update(pod model.Pod) error {
	err := s.repo.Update(pod)
	if err != nil {
		return err
	}
	return nil
}

func (s *PodService) Delete(id string) error {
	err := s.repo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
