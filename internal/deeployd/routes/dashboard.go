package routes

import (
	"fmt"
	"net/http"

	"github.com/deeploy-sh/deeploy/internal/deeployd/app"
	mw "github.com/deeploy-sh/deeploy/internal/shared/middleware"
	"github.com/deeploy-sh/deeploy/internal/shared/repo"
	"github.com/deeploy-sh/deeploy/internal/shared/service"
)

func Dashboard(app app.App) {
	userRepo := repo.NewUserRepo(app.DB)
	userService := service.NewUserService(userRepo)

	auth := mw.NewAuthMiddleware(userService)

	app.Router.HandleFunc("GET /api/dashboard", auth.Auth(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello from server!")
	}))
}
