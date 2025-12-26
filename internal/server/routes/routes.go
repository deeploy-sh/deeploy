package routes

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/deeploy-sh/deeploy/internal/server/app"
	handlers "github.com/deeploy-sh/deeploy/internal/server/handler"
	mw "github.com/deeploy-sh/deeploy/internal/server/middleware"
	sharedAssets "github.com/deeploy-sh/deeploy/internal/shared/assets"
	"github.com/deeploy-sh/deeploy/internal/shared/version"
)

func Setup(app *app.App) http.Handler {
	mux := http.NewServeMux()

	// Handlers
	auth := mw.NewAuthMiddleware(app.UserService)
	userHandler := handlers.NewUserHandler(app.UserService)
	projectHandler := handlers.NewProjectHandler(app.ProjectService, app.PodService)
	podHandler := handlers.NewPodHandler(app.PodService)
	gitTokenHandler := handlers.NewGitTokenHandler(app.GitTokenService)
	deployHandler := handlers.NewDeployHandler(app.DeployService)
	podDomainHandler := handlers.NewPodDomainHandler(app.PodDomainService, app.PodService, app.Cfg.IsDevelopment())
	podEnvVarHandler := handlers.NewPodEnvVarHandler(app.PodEnvVarService, app.PodService)
	serverSettingsHandler := handlers.NewServerSettingsHandler(app.TraefikService)

	// Assets
	setupAssets(mux, app.Cfg.IsDevelopment())

	// Landing
	mux.HandleFunc("GET /", userHandler.LandingView)

	// Auth
	mux.HandleFunc("GET /auth", mw.RequireGuest(userHandler.AuthView))
	mux.HandleFunc("POST /login", userHandler.Login)
	mux.HandleFunc("POST /register", userHandler.Register)
	mux.HandleFunc("GET /logout", userHandler.Logout)
	mux.HandleFunc("GET /api/auth/poll", userHandler.PollCLISession)

	// Projects
	mux.HandleFunc("POST /api/projects", auth.Auth(projectHandler.Create))
	mux.HandleFunc("GET /api/projects", auth.Auth(projectHandler.ProjectsByUser))
	mux.HandleFunc("GET /api/projects/{id}", auth.Auth(projectHandler.Project))
	mux.HandleFunc("PUT /api/projects", auth.Auth(projectHandler.Update))
	mux.HandleFunc("DELETE /api/projects/{id}", auth.Auth(projectHandler.Delete))

	// Pods
	mux.HandleFunc("POST /api/pods", auth.Auth(podHandler.Create))
	mux.HandleFunc("GET /api/pods", auth.Auth(podHandler.PodsByUser))
	mux.HandleFunc("GET /api/projects/{id}/pods", auth.Auth(podHandler.PodsByProject))
	mux.HandleFunc("GET /api/pods/{id}", auth.Auth(podHandler.Pod))
	mux.HandleFunc("PUT /api/pods", auth.Auth(podHandler.Update))
	mux.HandleFunc("DELETE /api/pods/{id}", auth.Auth(podHandler.Delete))

	// Pod Deploy
	mux.HandleFunc("POST /api/pods/{id}/deploy", auth.Auth(deployHandler.Deploy))
	mux.HandleFunc("POST /api/pods/{id}/stop", auth.Auth(deployHandler.Stop))
	mux.HandleFunc("POST /api/pods/{id}/restart", auth.Auth(deployHandler.Restart))
	mux.HandleFunc("GET /api/pods/{id}/logs", auth.Auth(deployHandler.Logs))

	// Pod Domains
	mux.HandleFunc("POST /api/pods/{id}/domains", auth.Auth(podDomainHandler.Create))
	mux.HandleFunc("POST /api/pods/{id}/domains/generate", auth.Auth(podDomainHandler.Generate))
	mux.HandleFunc("GET /api/pods/{id}/domains", auth.Auth(podDomainHandler.List))
	mux.HandleFunc("PUT /api/pods/{id}/domains/{domainId}", auth.Auth(podDomainHandler.Update))
	mux.HandleFunc("DELETE /api/pods/{id}/domains/{domainId}", auth.Auth(podDomainHandler.Delete))

	// Pod Env Vars
	mux.HandleFunc("GET /api/pods/{id}/vars", auth.Auth(podEnvVarHandler.List))
	mux.HandleFunc("PUT /api/pods/{id}/vars", auth.Auth(podEnvVarHandler.BulkUpdate))

	// Git Tokens
	mux.HandleFunc("POST /api/git-tokens", auth.Auth(gitTokenHandler.Create))
	mux.HandleFunc("GET /api/git-tokens", auth.Auth(gitTokenHandler.List))
	mux.HandleFunc("DELETE /api/git-tokens/{id}", auth.Auth(gitTokenHandler.Delete))

	// Server Settings
	mux.HandleFunc("GET /api/settings/domain", auth.Auth(serverSettingsHandler.GetServerDomain))
	mux.HandleFunc("PUT /api/settings/domain", auth.Auth(serverSettingsHandler.SetServerDomain))
	mux.HandleFunc("DELETE /api/settings/domain", auth.Auth(serverSettingsHandler.DeleteServerDomain))

	// Health (public - used by TUI for connection check + heartbeat)
	mux.HandleFunc("GET /api/health", healthHandler)

	return mux
}

func setupAssets(mux *http.ServeMux, isDev bool) {
	serve := func(devPath string, fs http.FileSystem) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if isDev {
				w.Header().Set("Cache-Control", "no-store")
				http.FileServer(http.Dir(devPath)).ServeHTTP(w, r)
			} else {
				http.FileServer(fs).ServeHTTP(w, r)
			}
		})
	}

	mux.Handle("GET /assets/css/", http.StripPrefix("/assets/", serve("./internal/shared/assets", http.FS(sharedAssets.Assets))))
	mux.Handle("GET /assets/fonts/", http.StripPrefix("/assets/", serve("./internal/shared/assets", http.FS(sharedAssets.Assets))))
	mux.Handle("GET /assets/js/", http.StripPrefix("/assets/", serve("./internal/shared/assets", http.FS(sharedAssets.Assets))))
	mux.Handle("GET /assets/img/", http.StripPrefix("/assets/", serve("./internal/shared/assets", http.FS(sharedAssets.Assets))))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Service string `json:"service"`
		Version string `json:"version"`
	}{
		Service: "deeploy",
		Version: version.Version,
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		slog.Error("health handler failed", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
