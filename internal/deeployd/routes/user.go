package routes

import (
	"github.com/deeploy-sh/deeploy/internal/deeployd/app"
	handlers "github.com/deeploy-sh/deeploy/internal/deeployd/handler"
	mw "github.com/deeploy-sh/deeploy/internal/deeployd/middleware"
	"github.com/deeploy-sh/deeploy/internal/deeployd/repo"
	"github.com/deeploy-sh/deeploy/internal/deeployd/service"
)

func User(app app.App) {
	userRepo := repo.NewUserRepo(app.DB)
	userService := service.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	// Views
	app.Router.HandleFunc("GET /", mw.RequireGuest(userHandler.AuthView))

	// APIs
	app.Router.HandleFunc("POST /login", userHandler.Login)
	app.Router.HandleFunc("POST /register", userHandler.Register)
	app.Router.HandleFunc("GET /logout", userHandler.Logout)
}
