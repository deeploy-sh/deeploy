package service

import (
	"github.com/deeploy-sh/deeploy/internal/deeployd/repo"
)

type PodServiceInterface interface {
	Create(pod *repo.Pod) (*repo.Pod, error)
	Pod(id string) (*repo.Pod, error)
	PodsByProject(id string) ([]repo.Pod, error)
	Update(pod repo.Pod) error
	Delete(id string) error
}

type PodService struct {
	repo repo.PodRepoInterface
}

func NewPodService(repo *repo.PodRepo) *PodService {
	return &PodService{repo: repo}
}

func (s *PodService) Create(pod *repo.Pod) (*repo.Pod, error) {
	err := s.repo.Create(pod)
	if err != nil {
		return nil, err
	}
	return pod, nil
}

func (s *PodService) Pod(id string) (*repo.Pod, error) {
	pod, err := s.repo.Pod(id)
	if err != nil {
		return nil, err
	}
	return pod, nil
}

func (s *PodService) PodsByProject(id string) ([]repo.Pod, error) {
	pods, err := s.repo.PodsByProject(id)
	if err != nil {
		return nil, err
	}
	return pods, nil
}

func (s *PodService) Update(pod repo.Pod) error {
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
