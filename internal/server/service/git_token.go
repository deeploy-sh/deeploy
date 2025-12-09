package service

import (
	"github.com/deeploy-sh/deeploy/internal/server/crypto"
	"github.com/deeploy-sh/deeploy/internal/server/repo"
)

type GitTokenServiceInterface interface {
	Create(token *repo.GitToken) (*repo.GitToken, error)
	GitToken(id string) (*repo.GitToken, error)
	GitTokensByUser(userID string) ([]repo.GitToken, error)
	Update(token repo.GitToken) error
	Delete(id string) error
}

type GitTokenService struct {
	repo      repo.GitTokenRepoInterface
	encryptor *crypto.Encryptor
}

func NewGitTokenService(repo *repo.GitTokenRepo, encryptor *crypto.Encryptor) *GitTokenService {
	return &GitTokenService{repo: repo, encryptor: encryptor}
}

func (s *GitTokenService) Create(token *repo.GitToken) (*repo.GitToken, error) {
	if s.encryptor != nil {
		encrypted, err := s.encryptor.Encrypt(token.Token)
		if err != nil {
			return nil, err
		}
		token.Token = encrypted
	}

	err := s.repo.Create(token)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (s *GitTokenService) GitToken(id string) (*repo.GitToken, error) {
	token, err := s.repo.GitToken(id)
	if err != nil {
		return nil, err
	}

	if s.encryptor != nil {
		decrypted, err := s.encryptor.Decrypt(token.Token)
		if err != nil {
			return nil, err
		}
		token.Token = decrypted
	}

	return token, nil
}

func (s *GitTokenService) GitTokensByUser(userID string) ([]repo.GitToken, error) {
	tokens, err := s.repo.GitTokensByUser(userID)
	if err != nil {
		return nil, err
	}

	if s.encryptor != nil {
		for i := range tokens {
			decrypted, err := s.encryptor.Decrypt(tokens[i].Token)
			if err != nil {
				return nil, err
			}
			tokens[i].Token = decrypted
		}
	}

	return tokens, nil
}

func (s *GitTokenService) Update(token repo.GitToken) error {
	if s.encryptor != nil {
		encrypted, err := s.encryptor.Encrypt(token.Token)
		if err != nil {
			return err
		}
		token.Token = encrypted
	}

	err := s.repo.Update(token)
	if err != nil {
		return err
	}
	return nil
}

func (s *GitTokenService) Delete(id string) error {
	err := s.repo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
