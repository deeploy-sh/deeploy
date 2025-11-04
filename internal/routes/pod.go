package routes

import (
	"github.com/axadrn/deeploy/internal/data"
	"github.com/axadrn/deeploy/internal/deeploy"
	apihandler "github.com/axadrn/deeploy/internal/handlers/api"
	mw "github.com/axadrn/deeploy/internal/middleware"
	"github.com/axadrn/deeploy/internal/services"
)

func Pod(app deeploy.App) {
	podRepo := data.NewPodRepo(app.DB)
	podService := services.NewPodService(podRepo)
	apiPodHandler := apihandler.NewPodHandler(podService)

	userRepo := data.NewUserRepo(app.DB)
	userService := services.NewUserService(userRepo)
	auth := mw.NewAuthMiddleware(userService)

	// API
	app.Router.HandleFunc("POST /api/pods", auth.Auth(apiPodHandler.Create))
	app.Router.HandleFunc("GET /api/pods/{id}", auth.Auth(apiPodHandler.Pod))
	app.Router.HandleFunc("GET /api/pods/project/{id}", auth.Auth(apiPodHandler.PodsByProject))
	app.Router.HandleFunc("PUT /api/pods", auth.Auth(apiPodHandler.Update))
	app.Router.HandleFunc("DELETE /api/pods/{id}", auth.Auth(apiPodHandler.Delete))

	// Web coming soon
}
