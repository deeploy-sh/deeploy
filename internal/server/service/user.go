package service

import (
	"github.com/deeploy-sh/deeploy/internal/server/auth"
	"github.com/deeploy-sh/deeploy/internal/server/forms"
	"github.com/deeploy-sh/deeploy/internal/server/jwt"
	"github.com/deeploy-sh/deeploy/internal/server/repo"
	"github.com/deeploy-sh/deeploy/internal/shared/errs"
	"github.com/deeploy-sh/deeploy/internal/shared/model"
	"github.com/google/uuid"
)

type UserServiceInterface interface {
	Register(form forms.RegisterForm) (string, error)
	Login(email, password string) (string, error)
	GetUserByID(id string) (*model.User, error)
	HasUser() (bool, error)
}

type UserService struct {
	repo repo.UserRepoInterface
}

func NewUserService(repo *repo.UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) HasUser() (bool, error) {
	count, err := s.repo.CountUsers()
	if err != nil {
		return false, err
	}
	hasUser := count > 0
	return hasUser, nil
}

func (s *UserService) Register(form forms.RegisterForm) (string, error) {
	foundUser, err := s.repo.GetUserByEmail(form.Email)
	if err != nil {
		return "", err
	}
	if foundUser != nil {
		return "", errs.ErrDuplicateEmail
	}
	hashedPwd, err := auth.HashPassword(form.Password)
	if err != nil {
		return "", err
	}
	user := &model.User{
		ID:       uuid.New().String(),
		Email:    form.Email,
		Password: hashedPwd,
	}
	err = s.repo.CreateUser(user)
	if err != nil {
		return "", err
	}
	token, err := jwt.CreateToken(user.ID)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *UserService) Login(email, password string) (string, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errs.ErrInvalidCredentials
	}
	if !auth.ComparePassword(user.Password, password) {
		return "", errs.ErrInvalidCredentials
	}
	token, err := jwt.CreateToken(user.ID)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *UserService) GetUserByID(id string) (*model.User, error) {
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
