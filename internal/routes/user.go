package routes

import (
	"github.com/axadrn/deeploy/internal/data"
	"github.com/axadrn/deeploy/internal/deeploy"
	handlers "github.com/axadrn/deeploy/internal/handlers/web"
	mw "github.com/axadrn/deeploy/internal/middleware"
	"github.com/axadrn/deeploy/internal/services"
)

func User(app deeploy.App) {
	userRepo := data.NewUserRepo(app.DB)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	// Views
	app.Router.HandleFunc("GET /", mw.RequireGuest(userHandler.AuthView))

	// APIs
	app.Router.HandleFunc("POST /login", userHandler.Login)
	app.Router.HandleFunc("POST /register", userHandler.Register)
	app.Router.HandleFunc("GET /logout", userHandler.Logout)
}
