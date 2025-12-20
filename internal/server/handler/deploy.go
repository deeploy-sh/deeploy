package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/deeploy-sh/deeploy/internal/server/service"
)

// Note: slog is kept for Deploy() goroutine logging

type DeployHandler struct {
	service *service.DeployService
}

func NewDeployHandler(service *service.DeployService) *DeployHandler {
	return &DeployHandler{service: service}
}

func (h *DeployHandler) Deploy(w http.ResponseWriter, r *http.Request) {
	podID := r.PathValue("id")
	slog.Info("Deploy request received", "podID", podID, "remoteAddr", r.RemoteAddr)

	// Start deploy in background
	go func() {
		err := h.service.Deploy(context.Background(), podID)
		if err != nil {
			slog.Error("deploy failed", "podID", podID, "error", err)
		}
	}()

	// Return immediately
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"status": "building"})
}


func (h *DeployHandler) Stop(w http.ResponseWriter, r *http.Request) {
	podID := r.PathValue("id")

	err := h.service.Stop(r.Context(), podID)
	if err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "stopped"})
}

func (h *DeployHandler) Restart(w http.ResponseWriter, r *http.Request) {
	podID := r.PathValue("id")

	err := h.service.Restart(r.Context(), podID)
	if err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "restarting"})
}

func (h *DeployHandler) Logs(w http.ResponseWriter, r *http.Request) {
	podID := r.PathValue("id")

	logs, status, err := h.service.GetLogs(r.Context(), podID, 100)
	if err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"logs":   logs,
		"status": status,
	})
}
