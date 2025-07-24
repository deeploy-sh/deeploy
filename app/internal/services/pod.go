package services

import (
	"github.com/axadrn/deeploy/internal/data"
)

type PodServiceInterface interface {
	Create(pod *data.Pod) (*data.Pod, error)
	Pod(id string) (*data.Pod, error)
	PodsByProject(id string) ([]data.Pod, error)
	Update(pod data.Pod) error
	Delete(id string) error
}

type PodService struct {
	repo data.PodRepoInterface
}

func NewPodService(repo *data.PodRepo) *PodService {
	return &PodService{repo: repo}
}

func (s *PodService) Create(pod *data.Pod) (*data.Pod, error) {
	err := s.repo.Create(pod)
	if err != nil {
		return nil, err
	}
	return pod, nil
}

func (s *PodService) Pod(id string) (*data.Pod, error) {
	pod, err := s.repo.Pod(id)
	if err != nil {
		return nil, err
	}
	return pod, nil
}

func (s *PodService) PodsByProject(id string) ([]data.Pod, error) {
	pods, err := s.repo.PodsByProject(id)
	if err != nil {
		return nil, err
	}
	return pods, nil
}

func (s *PodService) Update(pod data.Pod) error {
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
