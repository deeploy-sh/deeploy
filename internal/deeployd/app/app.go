package app

import (
	"github.com/deeploy-sh/deeploy/internal/deeployd/config"
	"github.com/deeploy-sh/deeploy/internal/deeployd/db"
	"github.com/deeploy-sh/deeploy/internal/deeployd/repo"
	"github.com/deeploy-sh/deeploy/internal/deeployd/service"
	"github.com/jmoiron/sqlx"
)

type App struct {
	Cfg            *config.Config
	DB             *sqlx.DB
	UserService    *service.UserService
	ProjectService *service.ProjectService
	PodService     *service.PodService
}

func New(cfg *config.Config) (*App, error) {
	database, err := db.Init(cfg.DBConnection)
	if err != nil {
		return nil, err
	}

	// Repositories
	userRepo := repo.NewUserRepo(database)
	projectRepo := repo.NewProjectRepo(database)
	podRepo := repo.NewPodRepo(database)

	// Services
	userService := service.NewUserService(userRepo)
	projectService := service.NewProjectService(projectRepo)
	podService := service.NewPodService(podRepo)

	return &App{
		Cfg:            cfg,
		DB:             database,
		UserService:    userService,
		ProjectService: projectService,
		PodService:     podService,
	}, nil
}

func (a *App) Close() error {
	return a.DB.Close()
}
