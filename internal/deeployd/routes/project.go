package routes

import (
	"github.com/deeploy-sh/deeploy/internal/shared/repo"
	"github.com/deeploy-sh/deeploy/internal/deeployd/app"
	apihandler "github.com/deeploy-sh/deeploy/internal/deeploy/handler"
	mw "github.com/deeploy-sh/deeploy/internal/shared/middleware"
	"github.com/deeploy-sh/deeploy/internal/shared/service"
)

func Project(app app.App) {
	projectRepo := repo.NewProjectRepo(app.DB)
	projectService := service.NewProjectService(projectRepo)
	apiProjectHandler := apihandler.NewProjectHandler(projectService)

	userRepo := repo.NewUserRepo(app.DB)
	userService := service.NewUserService(userRepo)
	auth := mw.NewAuthMiddleware(userService)

	// API
	app.Router.HandleFunc("POST /api/projects", auth.Auth(apiProjectHandler.Create))
	app.Router.HandleFunc("GET /api/projects/{id}", auth.Auth(apiProjectHandler.Project))
	app.Router.HandleFunc("GET /api/projects", auth.Auth(apiProjectHandler.ProjectsByUser))
	app.Router.HandleFunc("PUT /api/projects", auth.Auth(apiProjectHandler.Update))
	app.Router.HandleFunc("DELETE /api/projects/{id}", auth.Auth(apiProjectHandler.Delete))

	// Web coming soon
}
