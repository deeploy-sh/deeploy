package app

import (
	"github.com/deeploy-sh/deeploy/internal/deeployd/config"
	"github.com/deeploy-sh/deeploy/internal/deeployd/crypto"
	"github.com/deeploy-sh/deeploy/internal/deeployd/db"
	"github.com/deeploy-sh/deeploy/internal/deeployd/repo"
	"github.com/deeploy-sh/deeploy/internal/deeployd/service"
	"github.com/jmoiron/sqlx"
)

type App struct {
	Cfg              *config.Config
	DB               *sqlx.DB
	UserService      *service.UserService
	ProjectService   *service.ProjectService
	PodService       *service.PodService
	PodEnvVarService *service.PodEnvVarService
	PodDomainService *service.PodDomainService
	GitTokenService  *service.GitTokenService
}

func New(cfg *config.Config) (*App, error) {
	database, err := db.Init(cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	// Encryptor for env vars (nil in development if no key set)
	var encryptor *crypto.Encryptor
	if cfg.EncryptionKey != "" {
		encryptor, err = crypto.NewEncryptor(cfg.EncryptionKey)
		if err != nil {
			return nil, err
		}
	}

	// Repositories
	userRepo := repo.NewUserRepo(database)
	projectRepo := repo.NewProjectRepo(database)
	podRepo := repo.NewPodRepo(database)
	podEnvVarRepo := repo.NewPodEnvVarRepo(database)
	podDomainRepo := repo.NewPodDomainRepo(database)
	gitTokenRepo := repo.NewGitTokenRepo(database)

	// Services
	userService := service.NewUserService(userRepo)
	projectService := service.NewProjectService(projectRepo)
	podService := service.NewPodService(podRepo)
	podEnvVarService := service.NewPodEnvVarService(podEnvVarRepo, encryptor)
	podDomainService := service.NewPodDomainService(podDomainRepo)
	gitTokenService := service.NewGitTokenService(gitTokenRepo, encryptor)

	return &App{
		Cfg:              cfg,
		DB:               database,
		UserService:      userService,
		ProjectService:   projectService,
		PodService:       podService,
		PodEnvVarService: podEnvVarService,
		PodDomainService: podDomainService,
		GitTokenService:  gitTokenService,
	}, nil
}

func (a *App) Close() error {
	return a.DB.Close()
}
