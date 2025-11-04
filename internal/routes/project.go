package routes

import (
	"github.com/axadrn/deeploy/internal/data"
	"github.com/axadrn/deeploy/internal/deeploy"
	apihandler "github.com/axadrn/deeploy/internal/handlers/api"
	mw "github.com/axadrn/deeploy/internal/middleware"
	"github.com/axadrn/deeploy/internal/services"
)

func Project(app deeploy.App) {
	projectRepo := data.NewProjectRepo(app.DB)
	projectService := services.NewProjectService(projectRepo)
	apiProjectHandler := apihandler.NewProjectHandler(projectService)

	userRepo := data.NewUserRepo(app.DB)
	userService := services.NewUserService(userRepo)
	auth := mw.NewAuthMiddleware(userService)

	// API
	app.Router.HandleFunc("POST /api/projects", auth.Auth(apiProjectHandler.Create))
	app.Router.HandleFunc("GET /api/projects/{id}", auth.Auth(apiProjectHandler.Project))
	app.Router.HandleFunc("GET /api/projects", auth.Auth(apiProjectHandler.ProjectsByUser))
	app.Router.HandleFunc("PUT /api/projects", auth.Auth(apiProjectHandler.Update))
	app.Router.HandleFunc("DELETE /api/projects/{id}", auth.Auth(apiProjectHandler.Delete))

	// Web coming soon
}
