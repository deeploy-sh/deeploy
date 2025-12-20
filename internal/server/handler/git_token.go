package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/deeploy-sh/deeploy/internal/server/auth"
	"github.com/deeploy-sh/deeploy/internal/server/service"
	"github.com/deeploy-sh/deeploy/internal/shared/model"
	"github.com/google/uuid"
)

type GitTokenHandler struct {
	service *service.GitTokenService
}

func NewGitTokenHandler(service *service.GitTokenService) *GitTokenHandler {
	return &GitTokenHandler{service: service}
}

func (h *GitTokenHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.GitTokenCreate

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Provider == "" || req.Token == "" {
		http.Error(w, "Name, provider and token are required", http.StatusBadRequest)
		return
	}

	userID := auth.GetUser(r.Context()).ID

	gitToken := &model.GitToken{
		ID:       uuid.New().String(),
		UserID:   userID,
		Name:     req.Name,
		Provider: req.Provider,
		Token:    req.Token,
	}

	token, err := h.service.Create(gitToken)
	if err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(token)
}

func (h *GitTokenHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUser(r.Context()).ID

	tokens, err := h.service.GitTokensByUser(userID)
	if err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokens)
}

func (h *GitTokenHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := h.service.Delete(id)
	if err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
