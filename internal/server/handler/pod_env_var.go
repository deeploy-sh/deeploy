package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/deeploy-sh/deeploy/internal/server/service"
	"github.com/deeploy-sh/deeploy/internal/shared/model"
	"github.com/google/uuid"
)

type PodEnvVarHandler struct {
	service    *service.PodEnvVarService
	podService *service.PodService
}

func NewPodEnvVarHandler(service *service.PodEnvVarService, podService *service.PodService) *PodEnvVarHandler {
	return &PodEnvVarHandler{
		service:    service,
		podService: podService,
	}
}
func (h *PodEnvVarHandler) List(w http.ResponseWriter, r *http.Request) {
	podID := r.PathValue("id")

	envVars, err := h.service.EnvVarsByPod(podID)
	if err != nil {
		slog.Error("failed to get env vars", "podID", podID, "error", err)
		http.Error(w, "Failed to get env vars", http.StatusInternalServerError)
		return
	}

	response := make([]model.PodEnvVar, len(envVars))
	for i, v := range envVars {
		response[i] = model.PodEnvVar{
			ID:    v.ID,
			Key:   v.Key,
			Value: v.Value,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

type bulkUpdateRequest struct {
	Vars []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"vars"`
}

// BulkUpdate replaces all env vars for a pod (delete all + create new)
func (h *PodEnvVarHandler) BulkUpdate(w http.ResponseWriter, r *http.Request) {
	podID := r.PathValue("id")

	var req bulkUpdateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Delete all existing vars for this pod
	err = h.service.DeleteByPod(podID)
	if err != nil {
		slog.Error("failed to delete env vars", "podID", podID, "error", err)
		http.Error(w, "Failed to update env vars", http.StatusInternalServerError)
		return
	}

	// Create new vars
	for _, v := range req.Vars {
		if v.Key == "" {
			continue // skip empty keys
		}

		envVar := &model.PodEnvVar{
			ID:    uuid.New().String(),
			PodID: podID,
			Key:   v.Key,
			Value: v.Value,
		}

		_, err := h.service.Create(envVar)
		if err != nil {
			slog.Error("failed to create env var", "podID", podID, "key", v.Key, "error", err)
			http.Error(w, "Failed to create env var", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
