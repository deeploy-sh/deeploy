package service

import (
	"github.com/deeploy-sh/deeploy/internal/server/crypto"
	"github.com/deeploy-sh/deeploy/internal/server/repo"
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
	repo      repo.PodEnvVarRepoInterface
	encryptor *crypto.Encryptor
}

func NewPodEnvVarService(repo *repo.PodEnvVarRepo, encryptor *crypto.Encryptor) *PodEnvVarService {
	return &PodEnvVarService{repo: repo, encryptor: encryptor}
}

func (s *PodEnvVarService) Create(envVar *repo.PodEnvVar) (*repo.PodEnvVar, error) {
	if s.encryptor != nil {
		encrypted, err := s.encryptor.Encrypt(envVar.Value)
		if err != nil {
			return nil, err
		}
		envVar.Value = encrypted
	}

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

	if s.encryptor != nil {
		decrypted, err := s.encryptor.Decrypt(envVar.Value)
		if err != nil {
			return nil, err
		}
		envVar.Value = decrypted
	}

	return envVar, nil
}

func (s *PodEnvVarService) EnvVarsByPod(podID string) ([]repo.PodEnvVar, error) {
	envVars, err := s.repo.EnvVarsByPod(podID)
	if err != nil {
		return nil, err
	}

	if s.encryptor != nil {
		for i := range envVars {
			decrypted, err := s.encryptor.Decrypt(envVars[i].Value)
			if err != nil {
				return nil, err
			}
			envVars[i].Value = decrypted
		}
	}

	return envVars, nil
}

func (s *PodEnvVarService) Update(envVar repo.PodEnvVar) error {
	if s.encryptor != nil {
		encrypted, err := s.encryptor.Encrypt(envVar.Value)
		if err != nil {
			return err
		}
		envVar.Value = encrypted
	}

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
