package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/deeploy-sh/deeploy/internal/deeployd/repo"
	"github.com/deeploy-sh/deeploy/internal/deeployd/service"
	"github.com/google/uuid"
)

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

func (h *PodDomainHandler) Delete(w http.ResponseWriter, r *http.Request) {
	domainID := r.PathValue("domainId")

	if err := h.service.Delete(domainID); err != nil {
		log.Printf("Failed to delete domain: %v", err)
		http.Error(w, "Failed to delete domain", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type generateDomainRequest struct {
	Port       int  `json:"port"`
	SSLEnabled bool `json:"ssl_enabled"`
}

func (h *PodDomainHandler) Generate(w http.ResponseWriter, r *http.Request) {
	podID := r.PathValue("id")

	var req generateDomainRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Port == 0 {
		req.Port = 8080
	}

	// Get pod to use title for subdomain
	pod, err := h.podService.Pod(podID)
	if err != nil {
		log.Printf("Failed to get pod: %v", err)
		http.Error(w, "Pod not found", http.StatusNotFound)
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

	domain := &repo.PodDomain{
		ID:         uuid.New().String(),
		PodID:      podID,
		Domain:     domainName,
		Type:       "auto",
		Port:       req.Port,
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

// getPublicIP fetches the public IP via ipify.org (cached)
func (h *PodDomainHandler) getPublicIP() string {
	h.publicIPOnce.Do(func() {
		resp, err := http.Get("https://api.ipify.org")
		if err != nil {
			log.Printf("Failed to get public IP: %v", err)
			h.publicIP = "127-0-0-1" // fallback
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Failed to read public IP response: %v", err)
			h.publicIP = "127-0-0-1" // fallback
			return
		}

		h.publicIP = strings.TrimSpace(string(body))
		log.Printf("Detected public IP: %s", h.publicIP)
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
