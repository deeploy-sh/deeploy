package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/deeploy-sh/deeploy/internal/server/auth"
	"github.com/deeploy-sh/deeploy/internal/server/forms"
	"github.com/deeploy-sh/deeploy/internal/server/service"
	"github.com/deeploy-sh/deeploy/internal/shared/model"
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

	pod := &model.Pod{
		ID:        uuid.New().String(),
		UserID:    auth.GetUser(r.Context()).ID,
		ProjectID: form.ProjectID,
		Title:     form.Title,
	}

	_, err = h.service.Create(pod)
	if err != nil {
		writeError(w, err)
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
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pod)
}

func (h *PodHandler) PodsByProject(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("id")

	pods, err := h.service.PodsByProject(projectID)
	if err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pods)
}

func (h *PodHandler) PodsByUser(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUser(r.Context()).ID

	pods, err := h.service.PodsByUser(userID)
	if err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pods)
}

func (h *PodHandler) Update(w http.ResponseWriter, r *http.Request) {
	var pod model.Pod

	err := json.NewDecoder(r.Body).Decode(&pod)
	if err != nil {
		writeError(w, err)
		return
	}

	err = h.service.Update(pod)
	if err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pod)
}

func (h *PodHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := h.service.Delete(id)
	if err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 - Standard for successful DELETE
}
