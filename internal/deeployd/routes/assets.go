package routes

import (
	"net/http"
	"os"

	"github.com/deeploy-sh/deeploy/assets"
	"github.com/deeploy-sh/deeploy/internal/deeployd/app"
)

func Assets(app app.App) {
	var isDevelopment = os.Getenv("GO_ENV") != "production"
	assetHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var fs http.Handler
		if isDevelopment {
			w.Header().Set("Cache-Control", "no-store")
			fs = http.FileServer(http.Dir("./assets"))
		} else {
			fs = http.FileServer(http.FS(assets.Assets))
		}
		fs.ServeHTTP(w, r)
	})
	app.Router.Handle("GET /assets/", http.StripPrefix("/assets/", assetHandler))
}
