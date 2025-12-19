package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
	"strings"

	"github.com/a-h/templ"
	"github.com/deeploy-sh/deeploy/internal/docs/config"
	"github.com/deeploy-sh/deeploy/internal/docs/middleware"
	"github.com/deeploy-sh/deeploy/internal/docs/ui/pages"
	"github.com/deeploy-sh/deeploy/scripts"
	sharedAssets "github.com/deeploy-sh/deeploy/internal/shared/assets"
	"github.com/resend/resend-go/v2"
)

func main() {
	config.LoadConfig()
	mux := http.NewServeMux()
	setupAssets(mux)
	mux.Handle("GET /", templ.Handler(pages.Landing()))

	// Serve install scripts
	mux.Handle("GET /server.sh", serveScript("server.sh"))
	mux.Handle("GET /tui.sh", serveScript("tui.sh"))

	// Newsletter subscribe endpoint
	mux.HandleFunc("POST /api/subscribe", handleSubscribe)

	fmt.Println("Server is running on http://localhost:8090")
	http.ListenAndServe(":8090", middleware.GitHubStarsMiddleware(mux))
}

func serveScript(name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, err := scripts.Files.ReadFile(name)
		if err != nil {
			http.Error(w, "file not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Disposition", "attachment; filename="+name)
		w.Write(file)
	})
}

func handleSubscribe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		Email string `json:"email"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Invalid request",
		})
		return
	}

	email := strings.TrimSpace(strings.ToLower(req.Email))

	// Validate email
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Invalid email",
		})
		return
	}

	// Dev mode: just log
	if config.AppConfig.IsDevelopment() {
		slog.Info("newsletter subscription (dev mode)", "email", email)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
		})
		return
	}

	// Prod mode: add to Resend audience
	client := resend.NewClient(config.AppConfig.ResendAPIKey)
	_, err = client.Contacts.Create(&resend.CreateContactRequest{
		Email:      email,
		AudienceId: config.AppConfig.ResendAudienceID,
	})
	if err != nil {
		// Log error but return success to prevent email enumeration
		slog.Warn("newsletter subscription failed", "error", err, "email", email)
	} else {
		slog.Info("newsletter subscription successful", "email", email)
	}

	// Always return success to prevent email enumeration
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
	})
}

func setupAssets(mux *http.ServeMux) {
	isDev := config.AppConfig.IsDevelopment()

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
	mux.Handle("GET /assets/img/", http.StripPrefix("/assets/", serve("./internal/shared/assets", http.FS(sharedAssets.Assets))))
	mux.Handle("GET /assets/fonts/", http.StripPrefix("/assets/", serve("./internal/shared/assets", http.FS(sharedAssets.Assets))))
	mux.Handle("GET /assets/js/", http.StripPrefix("/assets/", serve("./internal/shared/assets", http.FS(sharedAssets.Assets))))
}
