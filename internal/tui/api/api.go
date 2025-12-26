package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	tea "charm.land/bubbletea/v2"
	"github.com/deeploy-sh/deeploy/internal/shared/errs"
	"github.com/deeploy-sh/deeploy/internal/shared/model"
	"github.com/deeploy-sh/deeploy/internal/tui/config"
	"github.com/deeploy-sh/deeploy/internal/tui/msg"
)

// --- Helpers ---

func getConfig() (*config.Config, error) {
	return config.Load()
}

// checkResponse checks the HTTP response for errors and returns an error if the status code indicates failure.
func checkResponse(resp *http.Response) error {
	if resp.StatusCode == http.StatusUnauthorized {
		return errs.ErrUnauthorized
	}
	if resp.StatusCode >= 400 {
		// Try to parse JSON error response
		var apiErr struct {
			Error string `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err == nil && apiErr.Error != "" {
			return errors.New(apiErr.Error)
		}
		return fmt.Errorf("request failed with status %d", resp.StatusCode)
	}
	return nil
}

func get(path string) (*http.Response, error) {
	cfg, err := getConfig()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", cfg.Server+"/api"+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+cfg.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	err = checkResponse(resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func post(path string, data any) (*http.Response, error) {
	cfg, err := getConfig()
	if err != nil {
		return nil, err
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", cfg.Server+"/api"+path, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+cfg.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	err = checkResponse(resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func put(path string, data any) (*http.Response, error) {
	cfg, err := getConfig()
	if err != nil {
		return nil, err
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", cfg.Server+"/api"+path, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+cfg.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	err = checkResponse(resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func del(path string) (*http.Response, error) {
	cfg, err := getConfig()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("DELETE", cfg.Server+"/api"+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+cfg.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	err = checkResponse(resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// --- Load All Data ---

func LoadData() tea.Cmd {
	return func() tea.Msg {
		projects, errP := fetchProjects()
		pods, errPod := fetchPods()
		gitTokens, errT := fetchGitTokens()

		if errP != nil {
			return msg.Error{Err: errP}
		}
		if errPod != nil {
			return msg.Error{Err: errPod}
		}
		if errT != nil {
			return msg.Error{Err: errT}
		}

		// Load all domains and env vars for all pods
		var podDomains []model.PodDomain
		var podEnvVars []model.PodEnvVar
		for _, p := range pods {
			domains, _ := fetchPodDomains(p.ID)
			podDomains = append(podDomains, domains...)
			vars, _ := fetchPodEnvVars(p.ID)
			podEnvVars = append(podEnvVars, vars...)
		}

		return msg.DataLoaded{
			Projects:   projects,
			Pods:       pods,
			GitTokens:  gitTokens,
			PodDomains: podDomains,
			PodEnvVars: podEnvVars,
		}
	}
}

// --- Projects ---

func fetchProjects() ([]model.Project, error) {
	resp, err := get("/projects")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var projects []model.Project
	err = json.NewDecoder(resp.Body).Decode(&projects)
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func CreateProject(title string) tea.Cmd {
	return func() tea.Msg {
		data := struct {
			Title string `json:"title"`
		}{Title: title}

		resp, err := post("/projects", data)
		if err != nil {
			return msg.Error{Err: err}
		}
		defer resp.Body.Close()

		var created model.Project
		if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
			return msg.Error{Err: err}
		}

		return msg.ProjectCreated{Project: created}
	}
}

func UpdateProject(project *model.Project) tea.Cmd {
	return func() tea.Msg {
		resp, err := put("/projects", project)
		if err != nil {
			return msg.Error{Err: err}
		}
		defer resp.Body.Close()

		var updated model.Project
		if err := json.NewDecoder(resp.Body).Decode(&updated); err != nil {
			return msg.Error{Err: err}
		}

		return msg.ProjectUpdated{Project: updated}
	}
}

func DeleteProject(id string) tea.Cmd {
	return func() tea.Msg {
		resp, err := del("/projects/" + id)
		if err != nil {
			return msg.Error{Err: err}
		}
		defer resp.Body.Close()

		return msg.ProjectDeleted{ProjectID: id}
	}
}

// --- Pods ---

func fetchPods() ([]model.Pod, error) {
	resp, err := get("/pods")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var pods []model.Pod
	err = json.NewDecoder(resp.Body).Decode(&pods)
	if err != nil {
		return nil, err
	}
	return pods, nil
}

func CreatePod(pod *model.Pod) tea.Cmd {
	return func() tea.Msg {
		resp, err := post("/pods", pod)
		if err != nil {
			return msg.Error{Err: err}
		}
		defer resp.Body.Close()

		var created model.Pod
		if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
			return msg.Error{Err: err}
		}

		return msg.PodCreated{Pod: created}
	}
}

func UpdatePod(pod *model.Pod) tea.Cmd {
	return func() tea.Msg {
		resp, err := put("/pods", pod)
		if err != nil {
			return msg.Error{Err: err}
		}
		defer resp.Body.Close()

		var updated model.Pod
		if err := json.NewDecoder(resp.Body).Decode(&updated); err != nil {
			return msg.Error{Err: err}
		}

		return msg.PodUpdated{Pod: updated}
	}
}

func DeletePod(id, projectID string) tea.Cmd {
	return func() tea.Msg {
		resp, err := del("/pods/" + id)
		if err != nil {
			return msg.Error{Err: err}
		}
		defer resp.Body.Close()

		return msg.PodDeleted{PodID: id, ProjectID: projectID}
	}
}

// --- Pod Deploy ---

func DeployPod(id string) tea.Cmd {
	return func() tea.Msg {
		resp, err := post("/pods/"+id+"/deploy", nil)
		if err != nil {
			return msg.Error{Err: err}
		}
		defer resp.Body.Close()

		return msg.PodDeployed{}
	}
}

func StopPod(id string) tea.Cmd {
	return func() tea.Msg {
		resp, err := post("/pods/"+id+"/stop", nil)
		if err != nil {
			return msg.Error{Err: err}
		}
		defer resp.Body.Close()

		return msg.PodStopped{}
	}
}

func RestartPod(id string) tea.Cmd {
	return func() tea.Msg {
		resp, err := post("/pods/"+id+"/restart", nil)
		if err != nil {
			return msg.Error{Err: err}
		}
		defer resp.Body.Close()

		return msg.PodRestarted{}
	}
}

func FetchPodLogs(id string) tea.Cmd {
	return func() tea.Msg {
		resp, err := get("/pods/" + id + "/logs")
		if err != nil {
			return msg.Error{Err: err}
		}
		defer resp.Body.Close()

		var result struct {
			Logs []string `json:"logs"`
		}
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			return msg.Error{Err: err}
		}

		return msg.PodLogsLoaded{Logs: result.Logs}
	}
}

// --- Git Tokens ---

func fetchGitTokens() ([]model.GitToken, error) {
	resp, err := get("/git-tokens")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tokens []model.GitToken
	err = json.NewDecoder(resp.Body).Decode(&tokens)
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

func CreateGitToken(name, provider, token string) tea.Cmd {
	return func() tea.Msg {
		data := struct {
			Name     string `json:"name"`
			Provider string `json:"provider"`
			Token    string `json:"token"`
		}{
			Name:     name,
			Provider: provider,
			Token:    token,
		}

		resp, err := post("/git-tokens", data)
		if err != nil {
			return msg.Error{Err: err}
		}
		defer resp.Body.Close()

		var created model.GitToken
		if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
			return msg.Error{Err: err}
		}

		return msg.GitTokenCreated{Token: created}
	}
}

func DeleteGitToken(id string) tea.Cmd {
	return func() tea.Msg {
		resp, err := del("/git-tokens/" + id)
		if err != nil {
			return msg.Error{Err: err}
		}
		defer resp.Body.Close()

		return msg.GitTokenDeleted{TokenID: id}
	}
}

// --- Pod Domains ---

func fetchPodDomains(podID string) ([]model.PodDomain, error) {
	resp, err := get("/pods/" + podID + "/domains")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var domains []model.PodDomain
	err = json.NewDecoder(resp.Body).Decode(&domains)
	if err != nil {
		return nil, err
	}
	return domains, nil
}

func CreatePodDomain(podID, domain string, port int, sslEnabled bool) tea.Cmd {
	return func() tea.Msg {
		data := struct {
			Domain     string `json:"domain"`
			Port       int    `json:"port"`
			SSLEnabled bool   `json:"ssl_enabled"`
		}{
			Domain:     domain,
			Port:       port,
			SSLEnabled: sslEnabled,
		}

		resp, err := post("/pods/"+podID+"/domains", data)
		if err != nil {
			return msg.Error{Err: err}
		}
		defer resp.Body.Close()

		var created model.PodDomain
		if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
			return msg.Error{Err: err}
		}

		return msg.PodDomainCreated{Domain: created}
	}
}

func DeletePodDomain(podID, domainID string) tea.Cmd {
	return func() tea.Msg {
		resp, err := del("/pods/" + podID + "/domains/" + domainID)
		if err != nil {
			return msg.Error{Err: err}
		}
		defer resp.Body.Close()

		return msg.PodDomainDeleted{DomainID: domainID, PodID: podID}
	}
}

func UpdatePodDomain(podID, domainID, domain string, port int, sslEnabled bool) tea.Cmd {
	return func() tea.Msg {
		data := model.PodDomain{
			Domain:     domain,
			Port:       port,
			SSLEnabled: sslEnabled,
		}

		resp, err := put("/pods/"+podID+"/domains/"+domainID, data)
		if err != nil {
			return msg.Error{Err: err}
		}
		defer resp.Body.Close()

		var updated model.PodDomain
		if err := json.NewDecoder(resp.Body).Decode(&updated); err != nil {
			return msg.Error{Err: err}
		}

		return msg.PodDomainUpdated{Domain: updated}
	}
}

func GenerateAutoDomain(podID string, port int, sslEnabled bool) tea.Cmd {
	return func() tea.Msg {
		data := struct {
			Port       int  `json:"port"`
			SSLEnabled bool `json:"ssl_enabled"`
		}{
			Port:       port,
			SSLEnabled: sslEnabled,
		}

		resp, err := post("/pods/"+podID+"/domains/generate", data)
		if err != nil {
			return msg.Error{Err: err}
		}
		defer resp.Body.Close()

		var created model.PodDomain
		if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
			return msg.Error{Err: err}
		}

		return msg.PodDomainCreated{Domain: created}
	}
}

// --- Pod Env Vars ---

func fetchPodEnvVars(podID string) ([]model.PodEnvVar, error) {
	resp, err := get("/pods/" + podID + "/vars")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var envVars []model.PodEnvVar
	err = json.NewDecoder(resp.Body).Decode(&envVars)
	if err != nil {
		return nil, err
	}
	return envVars, nil
}

func UpdatePodEnvVars(podID string, vars []model.PodEnvVar) tea.Cmd {
	return func() tea.Msg {
		data := struct {
			Vars []model.PodEnvVar `json:"vars"`
		}{Vars: vars}

		resp, err := put("/pods/"+podID+"/vars", data)
		if err != nil {
			return msg.Error{Err: err}
		}
		defer resp.Body.Close()

		var updated []model.PodEnvVar
		err = json.NewDecoder(resp.Body).Decode(&updated)
		if err != nil {
			return msg.Error{Err: err}
		}

		return msg.PodEnvVarsUpdated{PodID: podID, EnvVars: updated}
	}
}

// --- Connection Check ---

func CheckConnection() tea.Cmd {
	return func() tea.Msg {
		cfg, err := config.Load()
		if err != nil || cfg == nil || cfg.Server == "" {
			return msg.ConnectionResult{NeedsSetup: true}
		}
		if cfg.Token == "" {
			return msg.ConnectionResult{NeedsAuth: true}
		}

		resp, err := get("/health")
		if err != nil {
			return msg.ConnectionResult{Offline: true}
		}
		defer resp.Body.Close()

		var health struct {
			Version string `json:"version"`
		}
		json.NewDecoder(resp.Body).Decode(&health)

		return msg.ConnectionResult{
			Online:        true,
			ServerVersion: health.Version,
		}
	}
}

// --- Server Settings ---

func GetServerDomain() tea.Cmd {
	return func() tea.Msg {
		resp, err := get("/settings/domain")
		if err != nil {
			return msg.Error{Err: err}
		}
		defer resp.Body.Close()

		var result struct {
			Domain string `json:"domain"`
		}
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			return msg.Error{Err: err}
		}

		return msg.ServerDomainLoaded{Domain: result.Domain}
	}
}

func SetServerDomain(domain string) tea.Cmd {
	return func() tea.Msg {
		data := struct {
			Domain string `json:"domain"`
		}{Domain: domain}

		resp, err := put("/settings/domain", data)
		if err != nil {
			return msg.Error{Err: err}
		}
		defer resp.Body.Close()

		return msg.ServerDomainSet{}
	}
}

func DeleteServerDomain() tea.Cmd {
	return func() tea.Msg {
		resp, err := del("/settings/domain")
		if err != nil {
			return msg.Error{Err: err}
		}
		defer resp.Body.Close()

		return msg.ServerDomainDeleted{}
	}
}

// --- Version Check ---

func CheckLatestVersion() tea.Cmd {
	return func() tea.Msg {
		resp, err := http.Get("https://api.github.com/repos/deeploy-sh/deeploy/releases/latest")
		if err != nil {
			return msg.LatestVersionResult{Error: err}
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return msg.LatestVersionResult{Error: fmt.Errorf("github api returned %d", resp.StatusCode)}
		}

		var release struct {
			TagName string `json:"tag_name"`
		}
		err = json.NewDecoder(resp.Body).Decode(&release)
		if err != nil {
			return msg.LatestVersionResult{Error: err}
		}

		return msg.LatestVersionResult{Version: release.TagName}
	}
}
