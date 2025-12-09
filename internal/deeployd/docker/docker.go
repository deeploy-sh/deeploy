package docker

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/go-connections/nat"
)

type DockerService struct {
	client   *client.Client
	buildDir string
}

func NewDockerService(buildDir string) (*DockerService, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	// Ensure build directory exists
	if err := os.MkdirAll(buildDir, 0755); err != nil {
		return nil, err
	}

	return &DockerService{
		client:   cli,
		buildDir: buildDir,
	}, nil
}

// CloneRepo clones a git repository. Token is optional (for private repos).
func (d *DockerService) CloneRepo(repoURL, branch, token string) (string, error) {
	// Parse URL and inject token if provided
	cloneURL := repoURL
	if token != "" {
		parsed, err := url.Parse(repoURL)
		if err != nil {
			return "", fmt.Errorf("invalid repo URL: %w", err)
		}
		parsed.User = url.User(token)
		cloneURL = parsed.String()
	}

	// Create unique directory for this clone
	repoName := filepath.Base(strings.TrimSuffix(repoURL, ".git"))
	cloneDir := filepath.Join(d.buildDir, fmt.Sprintf("%s-%d", repoName, os.Getpid()))

	// Remove if exists
	os.RemoveAll(cloneDir)

	// Clone
	cmd := exec.Command("git", "clone", "--depth", "1", "--branch", branch, cloneURL, cloneDir)
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git clone failed: %s - %w", string(output), err)
	}

	return cloneDir, nil
}

// BuildImage builds a Docker image from a directory with a Dockerfile.
// logCallback is called for each line of build output (can be nil).
func (d *DockerService) BuildImage(ctx context.Context, buildPath, dockerfilePath, imageName string, logCallback func(string)) (string, error) {
	// Create tar archive of build context
	tar, err := archive.TarWithOptions(buildPath, &archive.TarOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to create build context: %w", err)
	}
	defer tar.Close()

	// Build options
	opts := types.ImageBuildOptions{
		Tags:       []string{imageName},
		Dockerfile: dockerfilePath,
		Remove:     true,
	}

	// Build
	resp, err := d.client.ImageBuild(ctx, tar, opts)
	if err != nil {
		return "", fmt.Errorf("docker build failed: %w", err)
	}
	defer resp.Body.Close()

	// Stream build output
	scanner := bufio.NewScanner(resp.Body)
	var lastError string
	for scanner.Scan() {
		line := scanner.Text()

		// Docker build output is JSON
		var msg struct {
			Stream      string `json:"stream"`
			Error       string `json:"error"`
			ErrorDetail struct {
				Message string `json:"message"`
			} `json:"errorDetail"`
		}
		if err := json.Unmarshal([]byte(line), &msg); err == nil {
			if msg.Error != "" {
				lastError = msg.Error
				if logCallback != nil {
					logCallback("ERROR: " + msg.Error)
				}
			} else if msg.Stream != "" {
				// Remove trailing newline
				stream := strings.TrimSuffix(msg.Stream, "\n")
				if stream != "" && logCallback != nil {
					logCallback(stream)
				}
			}
		}
	}

	if lastError != "" {
		return "", fmt.Errorf("build failed: %s", lastError)
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("failed reading build output: %w", err)
	}

	return imageName, nil
}

// RunContainer starts a container with the given configuration.
func (d *DockerService) RunContainer(ctx context.Context, opts RunContainerOptions) (string, error) {
	// Container config
	config := &container.Config{
		Image: opts.ImageName,
		Env:   mapToEnvSlice(opts.EnvVars),
		Labels: map[string]string{
			"traefik.enable": "true",
			"traefik.http.routers." + opts.PodID + ".rule":                      fmt.Sprintf("Host(`%s`)", opts.Domain),
			"traefik.http.services." + opts.PodID + ".loadbalancer.server.port": fmt.Sprintf("%d", opts.Port),
			"deeploy.pod.id": opts.PodID,
		},
	}

	// Exposed port
	exposedPort := nat.Port(fmt.Sprintf("%d/tcp", opts.Port))
	config.ExposedPorts = nat.PortSet{exposedPort: struct{}{}}

	// Host config
	hostConfig := &container.HostConfig{
		RestartPolicy: container.RestartPolicy{Name: "unless-stopped"},
	}

	// Create container
	resp, err := d.client.ContainerCreate(ctx, config, hostConfig, nil, nil, opts.ContainerName)
	if err != nil {
		return "", fmt.Errorf("failed to create container: %w", err)
	}

	// Start container
	if err := d.client.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", fmt.Errorf("failed to start container: %w", err)
	}

	return resp.ID, nil
}

// StopContainer stops a running container.
func (d *DockerService) StopContainer(ctx context.Context, containerID string) error {
	timeout := 30
	return d.client.ContainerStop(ctx, containerID, container.StopOptions{Timeout: &timeout})
}

// RemoveContainer removes a container.
func (d *DockerService) RemoveContainer(ctx context.Context, containerID string) error {
	return d.client.ContainerRemove(ctx, containerID, container.RemoveOptions{Force: true})
}

// GetLogs returns a reader for container logs.
func (d *DockerService) GetLogs(ctx context.Context, containerID string, follow bool) (io.ReadCloser, error) {
	opts := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     follow,
		Tail:       "100",
	}

	return d.client.ContainerLogs(ctx, containerID, opts)
}

// GetLogsLines returns the last n lines of logs as a slice.
func (d *DockerService) GetLogsLines(ctx context.Context, containerID string, lines int) ([]string, error) {
	opts := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Tail:       fmt.Sprintf("%d", lines),
	}

	reader, err := d.client.ContainerLogs(ctx, containerID, opts)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	var result []string
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		// Docker logs have 8-byte header, skip it
		if len(line) > 8 {
			result = append(result, line[8:])
		}
	}

	return result, scanner.Err()
}

// Cleanup removes the cloned repo directory.
func (d *DockerService) Cleanup(path string) error {
	return os.RemoveAll(path)
}

// Close closes the Docker client.
func (d *DockerService) Close() error {
	return d.client.Close()
}

// RunContainerOptions holds options for running a container.
type RunContainerOptions struct {
	ImageName     string
	ContainerName string
	PodID         string
	Domain        string
	Port          int
	EnvVars       map[string]string
}

func mapToEnvSlice(m map[string]string) []string {
	var result []string
	for k, v := range m {
		result = append(result, fmt.Sprintf("%s=%s", k, v))
	}
	return result
}
