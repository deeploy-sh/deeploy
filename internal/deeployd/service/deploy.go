package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"

	"github.com/deeploy-sh/deeploy/internal/deeployd/docker"
	"github.com/deeploy-sh/deeploy/internal/deeployd/repo"
	"github.com/google/uuid"
)

type DeployService struct {
	podRepo       repo.PodRepoInterface
	podDomainRepo repo.PodDomainRepoInterface
	podEnvVarRepo repo.PodEnvVarRepoInterface
	gitTokenRepo  repo.GitTokenRepoInterface
	docker        *docker.DockerService
	baseDomain    string

	// Build logs storage (simple)
	buildLogsMu sync.RWMutex
	buildLogs   map[string][]string // podID -> log lines
}

func NewDeployService(
	podRepo *repo.PodRepo,
	podDomainRepo *repo.PodDomainRepo,
	podEnvVarRepo *repo.PodEnvVarRepo,
	gitTokenRepo *repo.GitTokenRepo,
	docker *docker.DockerService,
	baseDomain string,
) *DeployService {
	return &DeployService{
		podRepo:       podRepo,
		podDomainRepo: podDomainRepo,
		podEnvVarRepo: podEnvVarRepo,
		gitTokenRepo:  gitTokenRepo,
		docker:        docker,
		baseDomain:    baseDomain,
		buildLogs:     make(map[string][]string),
	}
}

// appendBuildLog adds a log line to the build logs for a pod.
func (s *DeployService) appendBuildLog(podID, line string) {
	s.buildLogsMu.Lock()
	defer s.buildLogsMu.Unlock()
	s.buildLogs[podID] = append(s.buildLogs[podID], line)
}

// GetBuildLogs returns the current build logs for a pod.
func (s *DeployService) GetBuildLogs(podID string) []string {
	s.buildLogsMu.RLock()
	defer s.buildLogsMu.RUnlock()
	return s.buildLogs[podID]
}

// clearBuildLogs removes build logs for a pod (called at start of new build).
func (s *DeployService) clearBuildLogs(podID string) {
	s.buildLogsMu.Lock()
	defer s.buildLogsMu.Unlock()
	s.buildLogs[podID] = nil
}

// Deploy builds and runs a container for a pod.
func (s *DeployService) Deploy(ctx context.Context, podID string) error {
	// Clear old build logs
	s.clearBuildLogs(podID)

	// 1. Get pod
	pod, err := s.podRepo.Pod(podID)
	if err != nil {
		s.appendBuildLog(podID, fmt.Sprintf("ERROR: pod not found: %v", err))
		return fmt.Errorf("pod not found: %w", err)
	}

	if pod.RepoURL == nil || *pod.RepoURL == "" {
		s.appendBuildLog(podID, "ERROR: pod has no repo URL configured")
		return fmt.Errorf("pod has no repo URL configured")
	}

	// 2. Update status to building
	pod.Status = "building"
	if err := s.podRepo.Update(*pod); err != nil {
		return fmt.Errorf("failed to update pod status: %w", err)
	}

	s.appendBuildLog(podID, "=== Starting deployment ===")
	s.appendBuildLog(podID, fmt.Sprintf("Repo: %s @ %s", *pod.RepoURL, pod.Branch))

	// 3. Get git token if configured
	var gitToken string
	if pod.GitTokenID != nil {
		token, err := s.gitTokenRepo.GitToken(*pod.GitTokenID)
		if err != nil {
			s.appendBuildLog(podID, fmt.Sprintf("ERROR: failed to get git token: %v", err))
			return fmt.Errorf("failed to get git token: %w", err)
		}
		gitToken = token.Token
		s.appendBuildLog(podID, "Using configured git token for private repo")
	}

	// 4. Clone repo
	s.appendBuildLog(podID, "Cloning repository...")
	clonePath, err := s.docker.CloneRepo(*pod.RepoURL, pod.Branch, gitToken)
	if err != nil {
		pod.Status = "failed"
		s.podRepo.Update(*pod)
		s.appendBuildLog(podID, fmt.Sprintf("ERROR: failed to clone repo: %v", err))
		return fmt.Errorf("failed to clone repo: %w", err)
	}
	defer s.docker.Cleanup(clonePath)
	s.appendBuildLog(podID, "Repository cloned successfully")

	// 5. Build image
	s.appendBuildLog(podID, "")
	s.appendBuildLog(podID, "=== Building Docker image ===")
	imageName := fmt.Sprintf("deeploy-%s:latest", podID)

	logCallback := func(line string) {
		s.appendBuildLog(podID, line)
	}
	_, err = s.docker.BuildImage(ctx, clonePath, pod.DockerfilePath, imageName, logCallback)
	if err != nil {
		pod.Status = "failed"
		s.podRepo.Update(*pod)
		s.appendBuildLog(podID, fmt.Sprintf("ERROR: failed to build image: %v", err))
		return fmt.Errorf("failed to build image: %w", err)
	}
	s.appendBuildLog(podID, "")
	s.appendBuildLog(podID, "=== Docker image built successfully ===")

	// 6. Create auto-domain if not exists
	subdomain := s.generateSubdomain(pod.Title)
	domain := fmt.Sprintf("%s.%s", subdomain, s.baseDomain)

	domains, _ := s.podDomainRepo.DomainsByPod(podID)
	hasAutoDomain := false
	for _, d := range domains {
		if d.Type == "auto" {
			hasAutoDomain = true
			domain = d.Domain
			break
		}
	}

	if !hasAutoDomain {
		s.appendBuildLog(podID, fmt.Sprintf("Creating auto-domain: %s", domain))
		autoDomain := &repo.PodDomain{
			ID:         uuid.New().String(),
			PodID:      podID,
			Domain:     domain,
			Type:       "auto",
			Port:       80,
			IsPrimary:  true,
			SSLEnabled: true,
		}
		if err := s.podDomainRepo.Create(autoDomain); err != nil {
			pod.Status = "failed"
			s.podRepo.Update(*pod)
			s.appendBuildLog(podID, fmt.Sprintf("ERROR: failed to create auto-domain: %v", err))
			return fmt.Errorf("failed to create auto-domain: %w", err)
		}
	} else {
		s.appendBuildLog(podID, fmt.Sprintf("Using existing domain: %s", domain))
	}

	// 7. Get env vars
	envVars, err := s.podEnvVarRepo.EnvVarsByPod(podID)
	if err != nil {
		pod.Status = "failed"
		s.podRepo.Update(*pod)
		return fmt.Errorf("failed to get env vars: %w", err)
	}

	envMap := make(map[string]string)
	for _, ev := range envVars {
		envMap[ev.Key] = ev.Value
	}
	if len(envMap) > 0 {
		s.appendBuildLog(podID, fmt.Sprintf("Loaded %d environment variables", len(envMap)))
	}

	// 8. Stop existing container if running
	if pod.ContainerID != nil && *pod.ContainerID != "" {
		s.appendBuildLog(podID, "Stopping existing container...")
		s.docker.StopContainer(ctx, *pod.ContainerID)
		s.docker.RemoveContainer(ctx, *pod.ContainerID)
	}

	// 9. Run container
	s.appendBuildLog(podID, "")
	s.appendBuildLog(podID, "=== Starting container ===")
	containerID, err := s.docker.RunContainer(ctx, docker.RunContainerOptions{
		ImageName:     imageName,
		ContainerName: fmt.Sprintf("deeploy-%s", podID),
		PodID:         podID,
		Domain:        domain,
		Port:          80,
		EnvVars:       envMap,
	})
	if err != nil {
		pod.Status = "failed"
		s.podRepo.Update(*pod)
		s.appendBuildLog(podID, fmt.Sprintf("ERROR: failed to run container: %v", err))
		return fmt.Errorf("failed to run container: %w", err)
	}

	// 10. Update pod with container ID and status
	pod.ContainerID = &containerID
	pod.Status = "running"
	if err := s.podRepo.Update(*pod); err != nil {
		return fmt.Errorf("failed to update pod: %w", err)
	}

	s.appendBuildLog(podID, fmt.Sprintf("Container started: %s", containerID[:12]))
	s.appendBuildLog(podID, "")
	s.appendBuildLog(podID, "=== Deployment successful! ===")
	s.appendBuildLog(podID, fmt.Sprintf("Your app is available at: https://%s", domain))

	return nil
}

// Stop stops a running pod.
func (s *DeployService) Stop(ctx context.Context, podID string) error {
	pod, err := s.podRepo.Pod(podID)
	if err != nil {
		return fmt.Errorf("pod not found: %w", err)
	}

	if pod.ContainerID == nil || *pod.ContainerID == "" {
		return fmt.Errorf("pod has no running container")
	}

	if err := s.docker.StopContainer(ctx, *pod.ContainerID); err != nil {
		return fmt.Errorf("failed to stop container: %w", err)
	}

	pod.Status = "stopped"
	if err := s.podRepo.Update(*pod); err != nil {
		return fmt.Errorf("failed to update pod: %w", err)
	}

	return nil
}

// Restart restarts a pod (stop + deploy).
func (s *DeployService) Restart(ctx context.Context, podID string) error {
	// Stop first (ignore error if not running)
	s.Stop(ctx, podID)

	// Deploy
	return s.Deploy(ctx, podID)
}

// GetLogs returns build logs (if building) or container logs (if running).
func (s *DeployService) GetLogs(ctx context.Context, podID string, lines int) ([]string, string, error) {
	pod, err := s.podRepo.Pod(podID)
	if err != nil {
		return nil, "", fmt.Errorf("pod not found: %w", err)
	}

	// Return build logs if building or no container yet
	if pod.Status == "building" || pod.ContainerID == nil || *pod.ContainerID == "" {
		return s.GetBuildLogs(podID), pod.Status, nil
	}

	// Return container logs
	logs, err := s.docker.GetLogsLines(ctx, *pod.ContainerID, lines)
	return logs, pod.Status, err
}

// generateSubdomain creates a URL-safe subdomain from title + random suffix.
func (s *DeployService) generateSubdomain(title string) string {
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
