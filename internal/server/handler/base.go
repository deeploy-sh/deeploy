package handlers

import (
	"net/http"

	"github.com/deeploy-sh/deeploy/internal/server/ui/pages"
)

type DashboardHandler struct{}

func NewDashboardHandler() DashboardHandler {
	return DashboardHandler{}
}

func (*DashboardHandler) DashboardView(w http.ResponseWriter, r *http.Request) {
	pages.Dashboard().Render(r.Context(), w)
}
