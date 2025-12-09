package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/deeploy-sh/deeploy/internal/deeployd/service"
)

type DeployHandler struct {
	service *service.DeployService
}

func NewDeployHandler(service *service.DeployService) *DeployHandler {
	return &DeployHandler{service: service}
}

func (h *DeployHandler) Deploy(w http.ResponseWriter, r *http.Request) {
	podID := r.PathValue("id")

	// Start deploy in background
	go func() {
		if err := h.service.Deploy(context.Background(), podID); err != nil {
			log.Printf("Deploy failed for pod %s: %v", podID, err)
		}
	}()

	// Return immediately
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"status": "building"})
}


func (h *DeployHandler) Stop(w http.ResponseWriter, r *http.Request) {
	podID := r.PathValue("id")

	if err := h.service.Stop(r.Context(), podID); err != nil {
		log.Printf("Failed to stop pod %s: %v", podID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "stopped"})
}

func (h *DeployHandler) Restart(w http.ResponseWriter, r *http.Request) {
	podID := r.PathValue("id")

	if err := h.service.Restart(r.Context(), podID); err != nil {
		log.Printf("Failed to restart pod %s: %v", podID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "restarting"})
}

func (h *DeployHandler) Logs(w http.ResponseWriter, r *http.Request) {
	podID := r.PathValue("id")

	logs, status, err := h.service.GetLogs(r.Context(), podID, 100)
	if err != nil {
		log.Printf("Failed to get logs for pod %s: %v", podID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"logs":   logs,
		"status": status,
	})
}
