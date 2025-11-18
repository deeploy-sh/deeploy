package routes

import (
	"github.com/deeploy-sh/deeploy/internal/deeployd/app"
	handlers "github.com/deeploy-sh/deeploy/internal/deeployd/handler"
	"github.com/deeploy-sh/deeploy/internal/shared/repo"
	"github.com/deeploy-sh/deeploy/internal/shared/service"

	mw "github.com/deeploy-sh/deeploy/internal/shared/middleware"
)

func Base(app app.App) {
	dashboardHandler := handlers.NewDashboardHandler()
	userRepo := repo.NewUserRepo(app.DB)
	userService := service.NewUserService(userRepo)

	auth := mw.NewAuthMiddleware(userService)

	app.Router.HandleFunc("GET /dashboard", mw.RequireAuth(auth.Auth(dashboardHandler.DashboardView)))
}
