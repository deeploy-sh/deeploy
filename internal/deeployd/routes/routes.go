package routes

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/deeploy-sh/deeploy/assets"
	"github.com/deeploy-sh/deeploy/internal/deeployd/app"
	apihandler "github.com/deeploy-sh/deeploy/internal/deeploy/handler"
	handlers "github.com/deeploy-sh/deeploy/internal/deeployd/handler"
	mw "github.com/deeploy-sh/deeploy/internal/deeployd/middleware"
)

func Setup(app *app.App) http.Handler {
	mux := http.NewServeMux()

	// Handlers
	auth := mw.NewAuthMiddleware(app.UserService)
	userHandler := handlers.NewUserHandler(app.UserService)
	dashboardHandler := handlers.NewDashboardHandler()
	projectHandler := apihandler.NewProjectHandler(app.ProjectService)
	podHandler := apihandler.NewPodHandler(app.PodService)

	// Assets
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", assetHandler(app.Cfg.IsDevelopment())))

	// Auth
	mux.HandleFunc("GET /", mw.RequireGuest(userHandler.AuthView))
	mux.HandleFunc("POST /login", userHandler.Login)
	mux.HandleFunc("POST /register", userHandler.Register)
	mux.HandleFunc("GET /logout", userHandler.Logout)

	// Dashboard
	mux.HandleFunc("GET /dashboard", mw.RequireAuth(auth.Auth(dashboardHandler.DashboardView)))
	mux.HandleFunc("GET /api/dashboard", auth.Auth(dashboardPlaceholder))

	// Projects
	mux.HandleFunc("POST /api/projects", auth.Auth(projectHandler.Create))
	mux.HandleFunc("GET /api/projects", auth.Auth(projectHandler.ProjectsByUser))
	mux.HandleFunc("GET /api/projects/{id}", auth.Auth(projectHandler.Project))
	mux.HandleFunc("PUT /api/projects", auth.Auth(projectHandler.Update))
	mux.HandleFunc("DELETE /api/projects/{id}", auth.Auth(projectHandler.Delete))

	// Pods
	mux.HandleFunc("POST /api/pods", auth.Auth(podHandler.Create))
	mux.HandleFunc("GET /api/pods/{id}", auth.Auth(podHandler.Pod))
	mux.HandleFunc("GET /api/pods/project/{id}", auth.Auth(podHandler.PodsByProject))
	mux.HandleFunc("PUT /api/pods", auth.Auth(podHandler.Update))
	mux.HandleFunc("DELETE /api/pods/{id}", auth.Auth(podHandler.Delete))

	// Health
	mux.HandleFunc("GET /api/health", auth.Auth(healthHandler))

	// Middleware
	return mw.Chain(mux,
		mw.RequestLogging,
	)
}

func assetHandler(isDevelopment bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var fs http.Handler
		if isDevelopment {
			w.Header().Set("Cache-Control", "no-store")
			fs = http.FileServer(http.Dir("./assets"))
		} else {
			fs = http.FileServer(http.FS(assets.Assets))
		}
		fs.ServeHTTP(w, r)
	})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Service string `json:"service"`
		Version string `json:"version"`
	}{
		Service: "deeploy",
		Version: "0.0.1",
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		slog.Error("health handler failed", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func dashboardPlaceholder(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from server!"))
}
