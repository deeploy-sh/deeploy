package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/deeploy-sh/deeploy/internal/server/service"
)

type ServerSettingsHandler struct {
	traefik *service.TraefikService
}

func NewServerSettingsHandler(traefik *service.TraefikService) *ServerSettingsHandler {
	return &ServerSettingsHandler{traefik: traefik}
}

type domainResponse struct {
	Domain string `json:"domain"`
}

type setDomainRequest struct {
	Domain string `json:"domain"`
}

// GetServerDomain returns the current server domain.
func (h *ServerSettingsHandler) GetServerDomain(w http.ResponseWriter, r *http.Request) {
	domain, err := h.traefik.GetServerDomain()
	if err != nil {
		slog.Error("failed to get server domain", "error", err)
		http.Error(w, "Failed to get domain", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(domainResponse{Domain: domain})
}

// SetServerDomain saves the server domain and writes Traefik config.
func (h *ServerSettingsHandler) SetServerDomain(w http.ResponseWriter, r *http.Request) {
	var req setDomainRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Clean domain (remove protocol prefix if present)
	domain := strings.TrimPrefix(req.Domain, "https://")
	domain = strings.TrimPrefix(domain, "http://")
	domain = strings.TrimSuffix(domain, "/")

	if domain == "" {
		http.Error(w, "Domain is required", http.StatusBadRequest)
		return
	}

	if err := h.traefik.SetServerDomain(domain); err != nil {
		slog.Error("failed to set server domain", "error", err)
		http.Error(w, "Failed to set domain", http.StatusInternalServerError)
		return
	}

	slog.Info("server domain configured", "domain", domain)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(domainResponse{Domain: domain})
}

// DeleteServerDomain removes the server domain and Traefik config.
func (h *ServerSettingsHandler) DeleteServerDomain(w http.ResponseWriter, r *http.Request) {
	if err := h.traefik.DeleteServerDomain(); err != nil {
		slog.Error("failed to delete server domain", "error", err)
		http.Error(w, "Failed to delete domain", http.StatusInternalServerError)
		return
	}

	slog.Info("server domain removed")
	w.WriteHeader(http.StatusNoContent)
}
