package app

import (
	"github.com/deeploy-sh/deeploy/internal/server/config"
	"github.com/deeploy-sh/deeploy/internal/server/crypto"
	"github.com/deeploy-sh/deeploy/internal/server/db"
	"github.com/deeploy-sh/deeploy/internal/server/docker"
	"github.com/deeploy-sh/deeploy/internal/server/repo"
	"github.com/deeploy-sh/deeploy/internal/server/service"
	"github.com/jmoiron/sqlx"
)

type App struct {
	Cfg              *config.Config
	DB               *sqlx.DB
	Docker           *docker.DockerService
	UserService      *service.UserService
	ProjectService   *service.ProjectService
	PodService       *service.PodService
	PodEnvVarService *service.PodEnvVarService
	PodDomainService *service.PodDomainService
	GitTokenService  *service.GitTokenService
	DeployService    *service.DeployService
	TraefikService   *service.TraefikService
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

	// Docker service
	// isDevelopment determines if we use HTTP (dev) or HTTPS (prod) for Traefik routing
	dockerService, err := docker.NewDockerService(cfg.BuildDir, cfg.IsDevelopment())
	if err != nil {
		return nil, err
	}

	// Repositories
	userRepo := repo.NewUserRepo(database)
	projectRepo := repo.NewProjectRepo(database)
	podRepo := repo.NewPodRepo(database)
	podEnvVarRepo := repo.NewPodEnvVarRepo(database)
	podDomainRepo := repo.NewPodDomainRepo(database)
	gitTokenRepo := repo.NewGitTokenRepo(database)
	serverSettingsRepo := repo.NewServerSettingsRepo(database)

	// Services
	userService := service.NewUserService(userRepo)
	projectService := service.NewProjectService(projectRepo)
	podService := service.NewPodService(podRepo, dockerService)
	podEnvVarService := service.NewPodEnvVarService(podEnvVarRepo, encryptor)
	podDomainService := service.NewPodDomainService(podDomainRepo)
	gitTokenService := service.NewGitTokenService(gitTokenRepo, encryptor)
	deployService := service.NewDeployService(podRepo, podDomainRepo, podEnvVarRepo, gitTokenRepo, dockerService)
	traefikService := service.NewTraefikService(serverSettingsRepo, cfg.TraefikConfigDir, cfg.IsDevelopment())

	return &App{
		Cfg:              cfg,
		DB:               database,
		Docker:           dockerService,
		UserService:      userService,
		ProjectService:   projectService,
		PodService:       podService,
		PodEnvVarService: podEnvVarService,
		PodDomainService: podDomainService,
		GitTokenService:  gitTokenService,
		DeployService:    deployService,
		TraefikService:   traefikService,
	}, nil
}

func (a *App) Close() error {
	return a.DB.Close()
}
