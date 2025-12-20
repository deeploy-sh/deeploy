package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/deeploy-sh/deeploy/internal/server/auth"
	"github.com/deeploy-sh/deeploy/internal/server/forms"
	"github.com/deeploy-sh/deeploy/internal/server/service"
	"github.com/deeploy-sh/deeploy/internal/shared/errs"
	"github.com/deeploy-sh/deeploy/internal/shared/model"
	"github.com/google/uuid"
)

type ProjectHandler struct {
	service    service.ProjectServiceInterface
	podService service.PodServiceInterface
}

func NewProjectHandler(service *service.ProjectService, podService *service.PodService) *ProjectHandler {
	return &ProjectHandler{service: service, podService: podService}
}

func (h *ProjectHandler) Create(w http.ResponseWriter, r *http.Request) {
	var form forms.ProjectForm

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

	project := &model.Project{
		ID:     uuid.New().String(),
		UserID: auth.GetUser(r.Context()).ID,
		Title:  form.Title,
	}

	_, err = h.service.Create(project)
	if err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
	return
}

func (h *ProjectHandler) Project(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	project, err := h.service.Project(id)

	if err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
}

func (h *ProjectHandler) ProjectsByUser(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUser(r.Context()).ID

	projects, err := h.service.ProjectsByUser(userID)
	if err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

func (h *ProjectHandler) Update(w http.ResponseWriter, r *http.Request) {
	var project model.Project

	err := json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		writeError(w, err)
		return
	}

	err = h.service.Update(project)
	if err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
}

func (h *ProjectHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	// Check if project has pods
	podCount, err := h.podService.CountByProject(id)
	if err != nil {
		writeError(w, err)
		return
	}
	if podCount > 0 {
		writeError(w, fmt.Errorf("cannot delete project with existing pods: %w", errs.ErrConflict))
		return
	}

	err = h.service.Delete(id)
	if err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 - Standard for successful DELETE
}
