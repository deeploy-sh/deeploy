package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/deeploy-sh/deeploy/internal/server/docker"
	"github.com/deeploy-sh/deeploy/internal/server/repo"
)

type DeployService struct {
	podRepo          repo.PodRepoInterface
	podDomainRepo    repo.PodDomainRepoInterface
	podEnvVarService PodEnvVarServiceInterface
	gitTokenService  GitTokenServiceInterface
	docker           *docker.DockerService

	// Build logs storage (simple)
	buildLogsMu sync.RWMutex
	buildLogs   map[string][]string // podID -> log lines

	// Prevent parallel deploys of the same pod
	deployingMu sync.Mutex
	deploying   map[string]bool
}

func NewDeployService(
	podRepo *repo.PodRepo,
	podDomainRepo *repo.PodDomainRepo,
	podEnvVarService PodEnvVarServiceInterface,
	gitTokenService *GitTokenService,
	docker *docker.DockerService,
) *DeployService {
	return &DeployService{
		podRepo:          podRepo,
		podDomainRepo:    podDomainRepo,
		podEnvVarService: podEnvVarService,
		gitTokenService:  gitTokenService,
		docker:           docker,
		buildLogs:        make(map[string][]string),
		deploying:        make(map[string]bool),
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
	// Prevent parallel deploys of the same pod
	s.deployingMu.Lock()
	if s.deploying[podID] {
		s.deployingMu.Unlock()
		return fmt.Errorf("deploy already in progress for pod %s", podID)
	}
	s.deploying[podID] = true
	s.deployingMu.Unlock()

	defer func() {
		s.deployingMu.Lock()
		delete(s.deploying, podID)
		s.deployingMu.Unlock()
	}()

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

	// 2. Check domain exists BEFORE starting build (fail fast)
	domains, _ := s.podDomainRepo.DomainsByPod(podID)
	if len(domains) == 0 {
		s.appendBuildLog(podID, "ERROR: no domain configured - add a domain first")
		return fmt.Errorf("no domain configured for pod")
	}

	// 3. Update status to building
	pod.Status = "building"
	err = s.podRepo.Update(*pod)
	if err != nil {
		return fmt.Errorf("failed to update pod status: %w", err)
	}

	s.appendBuildLog(podID, "=== Starting deployment ===")
	s.appendBuildLog(podID, fmt.Sprintf("Repo: %s @ %s", *pod.RepoURL, pod.Branch))

	// 4. Get git token if configured
	var gitToken string
	if pod.GitTokenID != nil {
		token, err := s.gitTokenService.GitToken(*pod.GitTokenID)
		if err != nil {
			s.appendBuildLog(podID, fmt.Sprintf("ERROR: failed to get git token: %v", err))
			return fmt.Errorf("failed to get git token: %w", err)
		}
		gitToken = token.Token
		s.appendBuildLog(podID, "Using configured git token for private repo")
	}

	// 5. Clone repo
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

	// 6. Build image
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

	// 7. Prepare domains and env vars
	var domainConfigs []docker.DomainConfig
	for _, d := range domains {
		domainConfigs = append(domainConfigs, docker.DomainConfig{
			Domain: d.Domain,
			Port:   d.Port,
		})
		s.appendBuildLog(podID, fmt.Sprintf("Domain: %s (port %d)", d.Domain, d.Port))
	}

	// Get env vars (decrypted via service)
	envVars, err := s.podEnvVarService.EnvVarsByPod(podID)
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

	// 8. Rename existing container to make room for new one (zero-downtime)
	oldContainerID := ""
	containerName := fmt.Sprintf("deeploy-%s", podID)
	if pod.ContainerID != nil && *pod.ContainerID != "" {
		oldContainerID = *pod.ContainerID
		s.appendBuildLog(podID, "Preparing zero-downtime deployment...")
		err := s.docker.RenameContainer(ctx, oldContainerID, fmt.Sprintf("deeploy-%s-old", podID))
		if err != nil {
			// Container ID is stale - cleanup DB and orphaned container
			s.appendBuildLog(podID, "Cleaning up stale container...")
			pod.ContainerID = nil
			s.podRepo.Update(*pod)
			s.docker.StopContainer(ctx, containerName)
			s.docker.RemoveContainer(ctx, containerName)
			oldContainerID = "" // no rollback needed
		}
	}

	// 9. Clean up any orphaned container with target name (from previous failed deploys)
	// This handles edge case where deploy failed after container creation but before DB update
	s.docker.StopContainer(ctx, containerName)
	s.docker.RemoveContainer(ctx, containerName)

	// 10. Run new container (old still running for zero-downtime)
	s.appendBuildLog(podID, "")
	s.appendBuildLog(podID, "=== Starting new container ===")
	containerID, err := s.docker.RunContainer(ctx, docker.RunContainerOptions{
		ImageName:     imageName,
		ContainerName: fmt.Sprintf("deeploy-%s", podID),
		PodID:         podID,
		Domains:       domainConfigs,
		EnvVars:       envMap,
	})
	if err != nil {
		// Rollback: rename old container back
		if oldContainerID != "" {
			s.docker.RenameContainer(ctx, oldContainerID, fmt.Sprintf("deeploy-%s", podID))
		}
		pod.Status = "failed"
		s.podRepo.Update(*pod)
		s.appendBuildLog(podID, fmt.Sprintf("ERROR: failed to run container: %v", err))
		return fmt.Errorf("failed to run container: %w", err)
	}

	// 11. Stop old container (Traefik handles health checks and routing)
	// Wait 5s before stopping old container for zero-downtime deployment:
	// - Traefik health check interval is 2s
	// - Gives new container time to be detected and marked healthy
	// - Old container handles in-flight requests during transition
	// TODO: For v0.2+, could query Traefik API to confirm routing instead of fixed delay
	if oldContainerID != "" {
		time.Sleep(5 * time.Second)
		s.appendBuildLog(podID, "Stopping old container...")
		s.docker.StopContainer(ctx, oldContainerID)
		s.docker.RemoveContainer(ctx, oldContainerID)
	}

	// 12. Update pod with container ID and status
	pod.ContainerID = &containerID
	pod.Status = "running"
	err = s.podRepo.Update(*pod)
	if err != nil {
		return fmt.Errorf("failed to update pod: %w", err)
	}

	s.appendBuildLog(podID, fmt.Sprintf("Container started: %s", containerID[:12]))
	s.appendBuildLog(podID, "")
	s.appendBuildLog(podID, "=== Deployment successful! ===")
	s.appendBuildLog(podID, fmt.Sprintf("Your app is available at %d domain(s)", len(domains)))

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

	err = s.docker.StopContainer(ctx, *pod.ContainerID)
	if err != nil {
		return fmt.Errorf("failed to stop container: %w", err)
	}

	pod.Status = "stopped"
	err = s.podRepo.Update(*pod)
	if err != nil {
		return fmt.Errorf("failed to update pod: %w", err)
	}

	return nil
}

// Restart restarts a running container with current config (zero-downtime).
// Unlike Deploy, this does not rebuild the image - it reuses the existing one.
func (s *DeployService) Restart(ctx context.Context, podID string) error {
	// 1. Load pod
	pod, err := s.podRepo.Pod(podID)
	if err != nil {
		return fmt.Errorf("pod not found: %w", err)
	}

	if pod.ContainerID == nil || *pod.ContainerID == "" {
		return fmt.Errorf("pod has no running container")
	}
	oldContainerID := *pod.ContainerID

	// 2. Get image from current container
	imageName, err := s.docker.GetContainerImage(ctx, oldContainerID)
	if err != nil {
		// Container ID is stale - clear from DB
		pod.ContainerID = nil
		pod.Status = "stopped"
		s.podRepo.Update(*pod)

		// Cleanup any orphaned container with this name
		containerName := fmt.Sprintf("deeploy-%s", podID)
		s.docker.StopContainer(ctx, containerName)
		s.docker.RemoveContainer(ctx, containerName)

		return fmt.Errorf("container not found - use deploy instead")
	}

	// 3. Get domains
	domains, _ := s.podDomainRepo.DomainsByPod(podID)
	if len(domains) == 0 {
		return fmt.Errorf("no domain configured for pod")
	}

	var domainConfigs []docker.DomainConfig
	for _, d := range domains {
		domainConfigs = append(domainConfigs, docker.DomainConfig{
			Domain: d.Domain,
			Port:   d.Port,
		})
	}

	// 4. Get env vars (decrypted via service)
	envVars, err := s.podEnvVarService.EnvVarsByPod(podID)
	if err != nil {
		return fmt.Errorf("failed to get env vars: %w", err)
	}

	envMap := make(map[string]string)
	for _, ev := range envVars {
		envMap[ev.Key] = ev.Value
	}

	// 5. Rename old container (zero-downtime: keep running)
	s.docker.RenameContainer(ctx, oldContainerID, fmt.Sprintf("deeploy-%s-old", podID))

	// 6. Start new container
	containerID, err := s.docker.RunContainer(ctx, docker.RunContainerOptions{
		ImageName:     imageName,
		ContainerName: fmt.Sprintf("deeploy-%s", podID),
		PodID:         podID,
		Domains:       domainConfigs,
		EnvVars:       envMap,
	})
	if err != nil {
		// Rollback: rename old container back
		s.docker.RenameContainer(ctx, oldContainerID, fmt.Sprintf("deeploy-%s", podID))
		return fmt.Errorf("failed to run container: %w", err)
	}

	// 7. Stop old container (Traefik handles health checks and routing)
	// Wait 5s for zero-downtime (see Deploy function for detailed comment)
	time.Sleep(5 * time.Second)
	s.docker.StopContainer(ctx, oldContainerID)
	s.docker.RemoveContainer(ctx, oldContainerID)

	// 8. Update pod
	pod.ContainerID = &containerID
	pod.Status = "running"
	err = s.podRepo.Update(*pod)
	if err != nil {
		return fmt.Errorf("failed to update pod: %w", err)
	}

	return nil
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

