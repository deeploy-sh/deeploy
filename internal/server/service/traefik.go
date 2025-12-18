package service

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/deeploy-sh/deeploy/internal/server/repo"
)

const settingKeyServerDomain = "server_domain"

type TraefikService struct {
	settingsRepo *repo.ServerSettingsRepo
	configDir    string
	isDev        bool
}

func NewTraefikService(settingsRepo *repo.ServerSettingsRepo, configDir string, isDev bool) *TraefikService {
	return &TraefikService{
		settingsRepo: settingsRepo,
		configDir:    configDir,
		isDev:        isDev,
	}
}

// GetServerDomain returns the current server domain.
func (s *TraefikService) GetServerDomain() (string, error) {
	return s.settingsRepo.Get(settingKeyServerDomain)
}

// SetServerDomain saves the domain and writes Traefik config.
func (s *TraefikService) SetServerDomain(domain string) error {
	if err := s.settingsRepo.Set(settingKeyServerDomain, domain); err != nil {
		return fmt.Errorf("failed to save domain: %w", err)
	}

	if err := s.writeServerConfig(domain); err != nil {
		return fmt.Errorf("failed to write traefik config: %w", err)
	}

	return nil
}

// DeleteServerDomain removes the domain and deletes Traefik config.
func (s *TraefikService) DeleteServerDomain() error {
	if err := s.settingsRepo.Delete(settingKeyServerDomain); err != nil {
		return fmt.Errorf("failed to delete domain: %w", err)
	}

	configPath := filepath.Join(s.configDir, "server.yml")
	if err := os.Remove(configPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete traefik config: %w", err)
	}

	return nil
}

// writeServerConfig writes the Traefik YAML config file.
func (s *TraefikService) writeServerConfig(domain string) error {
	// Ensure directory exists
	if err := os.MkdirAll(s.configDir, 0755); err != nil {
		return err
	}

	// Determine entrypoint based on environment
	entrypoint := "websecure"
	tlsConfig := `
      tls:
        certResolver: letsencrypt`
	if s.isDev {
		entrypoint = "web"
		tlsConfig = "" // No TLS in development
	}

	config := fmt.Sprintf(`# Deeploy Server Domain Configuration
# Auto-generated - do not edit manually
# Traefik watches this file and updates routing automatically

http:
  routers:
    # Route for the Deeploy server domain
    deeploy-server:
      rule: "Host(%s)"
      service: deeploy-server
      entryPoints:
        - %s%s

  services:
    # Load balancer to the deeploy-app container
    deeploy-server:
      loadBalancer:
        servers:
          # Container name from docker-compose.yml
          - url: "http://deeploy-app:8090"
`, "`"+domain+"`", entrypoint, tlsConfig)

	configPath := filepath.Join(s.configDir, "server.yml")
	return os.WriteFile(configPath, []byte(config), 0644)
}
