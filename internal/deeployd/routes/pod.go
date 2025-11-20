package routes

import (
	apihandler "github.com/deeploy-sh/deeploy/internal/deeploy/handler"
	"github.com/deeploy-sh/deeploy/internal/deeployd/app"
	mw "github.com/deeploy-sh/deeploy/internal/deeployd/middleware"
	"github.com/deeploy-sh/deeploy/internal/deeployd/repo"
	"github.com/deeploy-sh/deeploy/internal/deeployd/service"
)

func Pod(app app.App) {
	podRepo := repo.NewPodRepo(app.DB)
	podService := service.NewPodService(podRepo)
	apiPodHandler := apihandler.NewPodHandler(podService)

	userRepo := repo.NewUserRepo(app.DB)
	userService := service.NewUserService(userRepo)
	auth := mw.NewAuthMiddleware(userService)

	// API
	app.Router.HandleFunc("POST /api/pods", auth.Auth(apiPodHandler.Create))
	app.Router.HandleFunc("GET /api/pods/{id}", auth.Auth(apiPodHandler.Pod))
	app.Router.HandleFunc("GET /api/pods/project/{id}", auth.Auth(apiPodHandler.PodsByProject))
	app.Router.HandleFunc("PUT /api/pods", auth.Auth(apiPodHandler.Update))
	app.Router.HandleFunc("DELETE /api/pods/{id}", auth.Auth(apiPodHandler.Delete))

	// Web coming soon
}
