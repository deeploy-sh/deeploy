package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/axadrn/deeploy/internal/auth"
	"github.com/axadrn/deeploy/internal/data"
	"github.com/axadrn/deeploy/internal/forms"
	"github.com/axadrn/deeploy/internal/services"
	"github.com/google/uuid"
)

type PodHandler struct {
	service services.PodServiceInterface
}

func NewPodHandler(service *services.PodService) *PodHandler {
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

	pod := &data.Pod{
		ID:          uuid.New().String(),
		UserID:      auth.GetUser(r.Context()).ID,
		ProjectID:   form.ProjectID,
		Title:       form.Title,
		Description: form.Description,
	}

	_, err = h.service.Create(pod)
	if err != nil {
		log.Printf("Failed to create pod: %v", err)
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
		log.Printf("Failed to get pod: %v", err)
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
		log.Printf("Failed to get pods: %v", err)
		http.Error(w, "Failed to get pods", http.StatusInternalServerError)
		return
	}

	dto := make([]data.PodDTO, len(pods))
	for i, pod := range pods {
		dto[i] = *pod.ToDTO()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dto)
}

func (h *PodHandler) Update(w http.ResponseWriter, r *http.Request) {
	var pod data.Pod

	err := json.NewDecoder(r.Body).Decode(&pod)
	if err != nil {
		log.Printf("Failed to decode json: %v", err)
		http.Error(w, "Failed to decode json", http.StatusInternalServerError)
		return
	}

	err = h.service.Update(pod)
	if err != nil {
		log.Printf("Failed to update pod: %v", err)
		http.Error(w, "Failed to update pod", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pod.ToDTO())
}

func (h *PodHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := h.service.Delete(id)
	if err != nil {
		log.Printf("Failed to delete pod: %v", err)
		http.Error(w, "Could not delete pod", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusNoContent) // 204 - Standard for successful DELETE
}
