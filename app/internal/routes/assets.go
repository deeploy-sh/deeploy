package routes

import (
	"net/http"
	"os"

	"github.com/axadrn/deeploy/assets"
	"github.com/axadrn/deeploy/internal/deeploy"
)

func Assets(app deeploy.App) {
	var isDevelopment = os.Getenv("GO_ENV") != "prod"
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
