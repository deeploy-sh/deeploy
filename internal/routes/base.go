package routes

import (
	"github.com/axadrn/deeploy/internal/data"
	"github.com/axadrn/deeploy/internal/deeploy"
	handlers "github.com/axadrn/deeploy/internal/handlers/web"
	"github.com/axadrn/deeploy/internal/services"

	mw "github.com/axadrn/deeploy/internal/middleware"
)

func Base(app deeploy.App) {
	dashboardHandler := handlers.NewDashboardHandler()
	userRepo := data.NewUserRepo(app.DB)
	userService := services.NewUserService(userRepo)

	auth := mw.NewAuthMiddleware(userService)

	app.Router.HandleFunc("GET /dashboard", mw.RequireAuth(auth.Auth(dashboardHandler.DashboardView)))
}
