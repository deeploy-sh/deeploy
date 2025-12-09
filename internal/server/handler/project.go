package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/deeploy-sh/deeploy/internal/server/auth"
	"github.com/deeploy-sh/deeploy/internal/server/forms"
	"github.com/deeploy-sh/deeploy/internal/server/repo"
	"github.com/deeploy-sh/deeploy/internal/server/service"
	"github.com/google/uuid"
)

type ProjectHandler struct {
	service service.ProjectServiceInterface
}

func NewProjectHandler(service *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{service: service}
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

	project := &repo.Project{
		ID:          uuid.New().String(),
		UserID:      auth.GetUser(r.Context()).ID,
		Title:       form.Title,
		Description: form.Description,
	}

	_, err = h.service.Create(project)
	if err != nil {
		slog.Error("failed to create project", "error", err)
		http.Error(w, "Failed to create project", http.StatusInternalServerError)
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
		slog.Warn("failed to get project", "projectID", id, "error", err)
		http.Error(w, "Failed to get project", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
}

func (h *ProjectHandler) ProjectsByUser(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUser(r.Context()).ID

	projects, err := h.service.ProjectsByUser(userID)
	if err != nil {
		slog.Error("failed to get projects by user", "userID", userID, "error", err)
		http.Error(w, "Failed to get projects", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

func (h *ProjectHandler) Update(w http.ResponseWriter, r *http.Request) {
	var project repo.Project

	err := json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		slog.Warn("failed to decode json", "error", err)
		http.Error(w, "Failed to decode json", http.StatusBadRequest)
		return
	}

	err = h.service.Update(project)
	if err != nil {
		slog.Error("failed to update project", "error", err)
		http.Error(w, "Failed to update project", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
}

func (h *ProjectHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := h.service.Delete(id)
	if err != nil {
		slog.Error("failed to delete project", "projectID", id, "error", err)
		http.Error(w, "Could not delete project", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 - Standard for successful DELETE
}
