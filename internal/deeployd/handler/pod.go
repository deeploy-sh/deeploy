package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/deeploy-sh/deeploy/internal/deeployd/auth"
	"github.com/deeploy-sh/deeploy/internal/deeployd/forms"
	"github.com/deeploy-sh/deeploy/internal/deeployd/repo"
	"github.com/deeploy-sh/deeploy/internal/deeployd/service"
	"github.com/google/uuid"
)

type PodHandler struct {
	service service.PodServiceInterface
}

func NewPodHandler(service *service.PodService) *PodHandler {
	return &PodHandler{service: service}
}

func (h *PodHandler) Create(w http.ResponseWriter, r *http.Request) {
	var form forms.PodForm

	err := json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		http.Error(w, "Failed to decode json", http.StatusInternalServerError)
		return
	}

	errors := form.Validate()
	if errors.HasErrors() {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors)
		return
	}

	pod := &repo.Pod{
		ID:          uuid.New().String(),
		UserID:      auth.GetUser(r.Context()).ID,
		ProjectID:   form.ProjectID,
		Title:       form.Title,
		Description: form.Description,
	}

	_, err = h.service.Create(pod)
	if err != nil {
		slog.Error("failed to create pod", "error", err)
		http.Error(w, "Failed to create pod", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pod)
	return
}

func (h *PodHandler) Pod(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	pod, err := h.service.Pod(id)

	if err != nil {
		slog.Warn("failed to get pod", "podID", id, "error", err)
		http.Error(w, "Failed to get pod", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pod)
}

func (h *PodHandler) PodsByProject(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("id")

	pods, err := h.service.PodsByProject(projectID)
	if err != nil {
		slog.Error("failed to get pods by project", "projectID", projectID, "error", err)
		http.Error(w, "Failed to get pods", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pods)
}

func (h *PodHandler) PodsByUser(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUser(r.Context()).ID

	pods, err := h.service.PodsByUser(userID)
	if err != nil {
		slog.Error("failed to get pods by user", "userID", userID, "error", err)
		http.Error(w, "Failed to get pods", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pods)
}

func (h *PodHandler) Update(w http.ResponseWriter, r *http.Request) {
	var pod repo.Pod

	err := json.NewDecoder(r.Body).Decode(&pod)
	if err != nil {
		slog.Warn("failed to decode json", "error", err)
		http.Error(w, "Failed to decode json", http.StatusBadRequest)
		return
	}

	err = h.service.Update(pod)
	if err != nil {
		slog.Error("failed to update pod", "error", err)
		http.Error(w, "Failed to update pod", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pod)
}

func (h *PodHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := h.service.Delete(id)
	if err != nil {
		slog.Error("failed to delete pod", "podID", id, "error", err)
		http.Error(w, "Could not delete pod", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 - Standard for successful DELETE
}
