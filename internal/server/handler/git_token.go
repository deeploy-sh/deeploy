package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/deeploy-sh/deeploy/internal/server/auth"
	"github.com/deeploy-sh/deeploy/internal/server/repo"
	"github.com/deeploy-sh/deeploy/internal/server/service"
	"github.com/google/uuid"
)

type GitTokenHandler struct {
	service *service.GitTokenService
}

func NewGitTokenHandler(service *service.GitTokenService) *GitTokenHandler {
	return &GitTokenHandler{service: service}
}

type createGitTokenRequest struct {
	Name     string `json:"name"`
	Provider string `json:"provider"`
	Token    string `json:"token"`
}

type gitTokenResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Provider  string `json:"provider"`
	CreatedAt string `json:"created_at"`
}

func (h *GitTokenHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createGitTokenRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Provider == "" || req.Token == "" {
		http.Error(w, "Name, provider and token are required", http.StatusBadRequest)
		return
	}

	userID := auth.GetUser(r.Context()).ID

	gitToken := &repo.GitToken{
		ID:       uuid.New().String(),
		UserID:   userID,
		Name:     req.Name,
		Provider: req.Provider,
		Token:    req.Token,
	}

	token, err := h.service.Create(gitToken)
	if err != nil {
		slog.Error("failed to create git token", "error", err)
		http.Error(w, "Failed to create git token", http.StatusInternalServerError)
		return
	}

	resp := gitTokenResponse{
		ID:        token.ID,
		Name:      token.Name,
		Provider:  token.Provider,
		CreatedAt: token.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *GitTokenHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUser(r.Context()).ID

	tokens, err := h.service.GitTokensByUser(userID)
	if err != nil {
		slog.Error("failed to get git tokens", "userID", userID, "error", err)
		http.Error(w, "Failed to get git tokens", http.StatusInternalServerError)
		return
	}

	resp := make([]gitTokenResponse, len(tokens))
	for i, t := range tokens {
		resp[i] = gitTokenResponse{
			ID:        t.ID,
			Name:      t.Name,
			Provider:  t.Provider,
			CreatedAt: t.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *GitTokenHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.service.Delete(id); err != nil {
		slog.Error("failed to delete git token", "tokenID", id, "error", err)
		http.Error(w, "Failed to delete git token", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
