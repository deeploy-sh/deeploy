package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"sync"

	"github.com/deeploy-sh/deeploy/internal/server/service"
	"github.com/deeploy-sh/deeploy/internal/shared/errs"
	"github.com/deeploy-sh/deeploy/internal/shared/model"
	"github.com/google/uuid"
)

// Note: slog is kept for getPublicIP() logging

type PodDomainHandler struct {
	service       *service.PodDomainService
	podService    *service.PodService
	isDevelopment bool
	publicIP      string
	publicIPOnce  sync.Once
}

func NewPodDomainHandler(service *service.PodDomainService, podService *service.PodService, isDevelopment bool) *PodDomainHandler {
	return &PodDomainHandler{
		service:       service,
		podService:    podService,
		isDevelopment: isDevelopment,
	}
}

// setDomainURL computes the full URL based on environment
func (h *PodDomainHandler) setDomainURL(d *model.PodDomain) {
	scheme := "https://"
	if h.isDevelopment {
		scheme = "http://"
	}
	d.URL = scheme + d.Domain
}

type createDomainRequest struct {
	Domain     string `json:"domain"`
	Port       int    `json:"port"`
	SSLEnabled bool   `json:"ssl_enabled"`
}

func (h *PodDomainHandler) Create(w http.ResponseWriter, r *http.Request) {
	podID := r.PathValue("id")

	var req createDomainRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
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

	// SSL is always enabled in production (automatic via Let's Encrypt)
	// The actual SSL handling is done by Traefik based on isDevelopment flag
	domain := &model.PodDomain{
		ID:         uuid.New().String(),
		PodID:      podID,
		Domain:     req.Domain,
		Type:       "custom",
		Port:       req.Port,
		SSLEnabled: true, // Always true - SSL is automatic in production
	}

	_, err = h.service.Create(domain)
	if err != nil {
		writeError(w, err)
		return
	}

	h.setDomainURL(domain)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(domain)
}

func (h *PodDomainHandler) List(w http.ResponseWriter, r *http.Request) {
	podID := r.PathValue("id")

	domains, err := h.service.DomainsByPod(podID)
	if err != nil {
		writeError(w, err)
		return
	}

	for i := range domains {
		h.setDomainURL(&domains[i])
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(domains)
}

func (h *PodDomainHandler) Delete(w http.ResponseWriter, r *http.Request) {
	domainID := r.PathValue("domainId")

	err := h.service.Delete(domainID)
	if err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *PodDomainHandler) Update(w http.ResponseWriter, r *http.Request) {
	domainID := r.PathValue("domainId")

	var req model.PodDomain
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
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

	// Get existing domain to preserve type
	existing, err := h.service.Domain(domainID)
	if err != nil {
		writeError(w, fmt.Errorf("domain %s: %w", domainID, errs.ErrNotFound))
		return
	}

	domain := model.PodDomain{
		ID:         domainID,
		PodID:      existing.PodID,
		Domain:     req.Domain,
		Type:       existing.Type, // Preserve original type
		Port:       req.Port,
		SSLEnabled: req.SSLEnabled,
	}

	err = h.service.Update(domain)
	if err != nil {
		writeError(w, err)
		return
	}

	h.setDomainURL(&domain)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(domain)
}

type generateDomainRequest struct {
	Port       int  `json:"port"`
	SSLEnabled bool `json:"ssl_enabled"`
}

func (h *PodDomainHandler) Generate(w http.ResponseWriter, r *http.Request) {
	podID := r.PathValue("id")

	var req generateDomainRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Port == 0 {
		req.Port = 8080
	}

	// Get pod to use title for subdomain
	pod, err := h.podService.Pod(podID)
	if err != nil {
		writeError(w, fmt.Errorf("pod %s: %w", podID, errs.ErrNotFound))
		return
	}

	// Generate subdomain from pod title
	subdomain := generateSubdomain(pod.Title)

	// Build sslip.io domain (wildcard DNS that resolves to embedded IP)
	// Format: subdomain.IP.sslip.io -> resolves to IP
	var ip string
	if h.isDevelopment {
		ip = "127.0.0.1"
	} else {
		ip = h.getPublicIP()
	}
	domainName := fmt.Sprintf("%s.%s.sslip.io", subdomain, ip)

	// SSL is always enabled in production (automatic via Let's Encrypt)
	// sslip.io domains work with HTTP challenge just like custom domains
	domain := &model.PodDomain{
		ID:         uuid.New().String(),
		PodID:      podID,
		Domain:     domainName,
		Type:       "auto",
		Port:       req.Port,
		SSLEnabled: true, // Always true - SSL is automatic in production
	}

	_, err = h.service.Create(domain)
	if err != nil {
		writeError(w, err)
		return
	}

	h.setDomainURL(domain)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(domain)
}

// getPublicIP fetches the public IP via ipify.org (cached)
func (h *PodDomainHandler) getPublicIP() string {
	h.publicIPOnce.Do(func() {
		resp, err := http.Get("https://api.ipify.org")
		if err != nil {
			slog.Warn("failed to get public IP", "error", err)
			h.publicIP = "127-0-0-1" // fallback
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Warn("failed to read public IP response", "error", err)
			h.publicIP = "127-0-0-1" // fallback
			return
		}

		h.publicIP = strings.TrimSpace(string(body))
		slog.Info("detected public IP", "ip", h.publicIP)
	})
	return h.publicIP
}

// generateSubdomain creates a URL-safe subdomain from title + random suffix.
func generateSubdomain(title string) string {
	// Sanitize title
	subdomain := strings.ToLower(title)
	subdomain = strings.ReplaceAll(subdomain, " ", "-")

	// Keep only alphanumeric and hyphens
	var result strings.Builder
	for _, r := range subdomain {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}
	subdomain = result.String()

	// Trim to max 20 chars
	if len(subdomain) > 20 {
		subdomain = subdomain[:20]
	}

	// Add random suffix
	suffix := make([]byte, 4)
	rand.Read(suffix)
	subdomain = fmt.Sprintf("%s-%s", subdomain, hex.EncodeToString(suffix))

	return subdomain
}
