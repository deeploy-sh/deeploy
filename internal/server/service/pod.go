package service

import (
	"context"
	"fmt"

	"github.com/deeploy-sh/deeploy/internal/server/docker"
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
	repo   repo.PodRepoInterface
	docker *docker.DockerService
}

func NewPodService(repo *repo.PodRepo, docker *docker.DockerService) *PodService {
	return &PodService{repo: repo, docker: docker}
}

// enrichWithContainerState fetches the live Docker container state for a pod.
// If the container doesn't exist or there's an error, ContainerState remains empty.
func (s *PodService) enrichWithContainerState(pod *model.Pod) {
	if pod.ContainerID == nil || *pod.ContainerID == "" {
		return
	}
	state, err := s.docker.GetContainerState(context.Background(), *pod.ContainerID)
	if err != nil {
		// Container not found or error - leave ContainerState empty
		return
	}
	pod.ContainerState = state
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
	s.enrichWithContainerState(pod)
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
	pod, err := s.repo.Pod(id)
	if err != nil {
		return err
	}

	s.cleanupDocker(id, pod.ContainerID)

	return s.repo.Delete(id)
}

func (s *PodService) cleanupDocker(podID string, containerID *string) {
	ctx := context.Background()

	if containerID != nil && *containerID != "" {
		s.docker.StopContainer(ctx, *containerID)
		s.docker.RemoveContainer(ctx, *containerID)
	}

	s.docker.RemoveImage(ctx, fmt.Sprintf("deeploy-%s:latest", podID))
}
