package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/deeploy-sh/deeploy/internal/deeployd/repo"
	"github.com/deeploy-sh/deeploy/internal/deeployd/service"
	"github.com/google/uuid"
)

type PodDomainHandler struct {
	service *service.PodDomainService
}

func NewPodDomainHandler(service *service.PodDomainService) *PodDomainHandler {
	return &PodDomainHandler{service: service}
}

type createDomainRequest struct {
	Domain     string `json:"domain"`
	Port       int    `json:"port"`
	SSLEnabled bool   `json:"ssl_enabled"`
}

func (h *PodDomainHandler) Create(w http.ResponseWriter, r *http.Request) {
	podID := r.PathValue("id")

	var req createDomainRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Domain == "" {
		http.Error(w, "Domain is required", http.StatusBadRequest)
		return
	}

	if req.Port == 0 {
		req.Port = 80
	}

	domain := &repo.PodDomain{
		ID:         uuid.New().String(),
		PodID:      podID,
		Domain:     req.Domain,
		Type:       "custom",
		Port:       req.Port,
		IsPrimary:  false,
		SSLEnabled: req.SSLEnabled,
	}

	if _, err := h.service.Create(domain); err != nil {
		log.Printf("Failed to create domain: %v", err)
		http.Error(w, "Failed to create domain", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(domain)
}

func (h *PodDomainHandler) List(w http.ResponseWriter, r *http.Request) {
	podID := r.PathValue("id")

	domains, err := h.service.DomainsByPod(podID)
	if err != nil {
		log.Printf("Failed to get domains: %v", err)
		http.Error(w, "Failed to get domains", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(domains)
}

func (h *PodDomainHandler) SetPrimary(w http.ResponseWriter, r *http.Request) {
	podID := r.PathValue("id")
	domainID := r.PathValue("domainId")

	if err := h.service.SetPrimary(domainID, podID); err != nil {
		log.Printf("Failed to set primary domain: %v", err)
		http.Error(w, "Failed to set primary domain", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *PodDomainHandler) Delete(w http.ResponseWriter, r *http.Request) {
	domainID := r.PathValue("domainId")

	if err := h.service.Delete(domainID); err != nil {
		log.Printf("Failed to delete domain: %v", err)
		http.Error(w, "Failed to delete domain", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
