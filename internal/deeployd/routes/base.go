package routes

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/deeploy-sh/deeploy/internal/deeployd/app"
	handlers "github.com/deeploy-sh/deeploy/internal/deeployd/handler"
	"github.com/deeploy-sh/deeploy/internal/deeployd/repo"
	"github.com/deeploy-sh/deeploy/internal/deeployd/service"

	mw "github.com/deeploy-sh/deeploy/internal/deeployd/middleware"
)

func Base(app app.App) {
	dashboardHandler := handlers.NewDashboardHandler()
	userRepo := repo.NewUserRepo(app.DB)
	userService := service.NewUserService(userRepo)

	auth := mw.NewAuthMiddleware(userService)

	app.Router.HandleFunc("GET /dashboard", mw.RequireAuth(auth.Auth(dashboardHandler.DashboardView)))
	app.Router.HandleFunc("GET /api/health", auth.Auth(healthHandler))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	healthCheck := struct {
		Service string
		Version string
	}{
		Service: "deeploy",
		Version: "0.0.1",
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(healthCheck)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
