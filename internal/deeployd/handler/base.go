package handlers

import (
	"net/http"

	"github.com/deeploy-sh/deeploy/internal/deeployd/ui/pages"
)

type DashboardHandler struct{}

func NewDashboardHandler() DashboardHandler {
	return DashboardHandler{}
}

func (*DashboardHandler) DashboardView(w http.ResponseWriter, r *http.Request) {
	pages.Dashboard().Render(r.Context(), w)
}

func (*DashboardHandler) LandingView(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
