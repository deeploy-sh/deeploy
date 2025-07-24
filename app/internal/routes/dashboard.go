package routes

import (
	"fmt"
	"net/http"

	"github.com/axadrn/deeploy/internal/data"
	"github.com/axadrn/deeploy/internal/deeploy"
	"github.com/axadrn/deeploy/internal/services"

	mw "github.com/axadrn/deeploy/internal/middleware"
)

func Dashboard(app deeploy.App) {
	userRepo := data.NewUserRepo(app.DB)
	userService := services.NewUserService(userRepo)

	auth := mw.NewAuthMiddleware(userService)

	app.Router.HandleFunc("GET /api/dashboard", auth.Auth(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello from server!")
	}))
}
