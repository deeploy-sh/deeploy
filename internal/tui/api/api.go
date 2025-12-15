package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	tea "charm.land/bubbletea/v2"
	"github.com/deeploy-sh/deeploy/internal/server/repo"
	"github.com/deeploy-sh/deeploy/internal/shared/errs"
	"github.com/deeploy-sh/deeploy/internal/shared/model"
	"github.com/deeploy-sh/deeploy/internal/tui/config"
	"github.com/deeploy-sh/deeploy/internal/tui/msg"
)

// --- Helpers ---

func getConfig() (*config.Config, error) {
	return config.Load()
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
	if resp.StatusCode == http.StatusUnauthorized {
		return nil, errs.ErrUnauthorized
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
	if resp.StatusCode == http.StatusUnauthorized {
		return nil, errs.ErrUnauthorized
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
	if resp.StatusCode == http.StatusUnauthorized {
		return nil, errs.ErrUnauthorized
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
	if resp.StatusCode == http.StatusUnauthorized {
		return nil, errs.ErrUnauthorized
	}
	return resp, nil
}

// --- Load All Data ---

func LoadData() tea.Cmd {
	return func() tea.Msg {
		projects, errP := fetchProjects()
		pods, errPod := fetchPods()

		if errP != nil {
			return msg.Error{Err: errP}
		}
		if errPod != nil {
			return msg.Error{Err: errPod}
		}

		return msg.DataLoaded{
			Projects: projects,
			Pods:     pods,
		}
	}
}

// --- Projects ---

func fetchProjects() ([]repo.Project, error) {
	resp, err := get("/projects")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var projects []repo.Project
	if err := json.NewDecoder(resp.Body).Decode(&projects); err != nil {
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

		return msg.ProjectCreated{}
	}
}

func UpdateProject(project *repo.Project) tea.Cmd {
	return func() tea.Msg {
		resp, err := put("/projects", project)
		if err != nil {
			return msg.Error{Err: err}
		}
		defer resp.Body.Close()

		return msg.ProjectUpdated{}
	}
}

func DeleteProject(id string) tea.Cmd {
	return func() tea.Msg {
		resp, err := del("/projects/" + id)
		if err != nil {
			return msg.Error{Err: err}
		}
		defer resp.Body.Close()

		return msg.ProjectDeleted{}
	}
}

// --- Pods ---

func fetchPods() ([]repo.Pod, error) {
	resp, err := get("/pods")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var pods []repo.Pod
	if err := json.NewDecoder(resp.Body).Decode(&pods); err != nil {
		return nil, err
	}
	return pods, nil
}

func CreatePod(title, projectID string) tea.Cmd {
	return func() tea.Msg {
		data := struct {
			Title     string `json:"title"`
			ProjectID string `json:"project_id"`
		}{
			Title:     title,
			ProjectID: projectID,
		}

		resp, err := post("/pods", data)
		if err != nil {
			return msg.Error{Err: err}
		}
		defer resp.Body.Close()

		return msg.PodCreated{}
	}
}

func UpdatePod(pod *repo.Pod) tea.Cmd {
	return func() tea.Msg {
		resp, err := put("/pods", pod)
		if err != nil {
			return msg.Error{Err: err}
		}
		defer resp.Body.Close()

		return msg.PodUpdated{}
	}
}

func DeletePod(id string) tea.Cmd {
	return func() tea.Msg {
		resp, err := del("/pods/" + id)
		if err != nil {
			return msg.Error{Err: err}
		}
		defer resp.Body.Close()

		return msg.PodDeleted{}
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
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return msg.Error{Err: err}
		}

		return msg.PodLogsLoaded{Logs: result.Logs}
	}
}

// --- Git Tokens ---

type GitToken struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Provider  string `json:"provider"`
	CreatedAt string `json:"created_at"`
}

func FetchGitTokens() tea.Cmd {
	return func() tea.Msg {
		resp, err := get("/git-tokens")
		if err != nil {
			return msg.Error{Err: err}
		}
		defer resp.Body.Close()

		var tokens []GitToken
		if err := json.NewDecoder(resp.Body).Decode(&tokens); err != nil {
			return msg.Error{Err: err}
		}

		return msg.GitTokensLoaded{Tokens: tokens}
	}
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

		return msg.GitTokenCreated{}
	}
}

func DeleteGitToken(id string) tea.Cmd {
	return func() tea.Msg {
		resp, err := del("/git-tokens/" + id)
		if err != nil {
			return msg.Error{Err: err}
		}
		defer resp.Body.Close()

		return msg.GitTokenDeleted{}
	}
}

// --- Pod Domains ---

type PodDomain struct {
	ID         string `json:"id"`
	PodID      string `json:"pod_id"`
	Domain     string `json:"domain"`
	Type       string `json:"type"`
	Port       int    `json:"port"`
	SSLEnabled bool   `json:"ssl_enabled"`
}

func FetchPodDomains(podID string) tea.Cmd {
	return func() tea.Msg {
		resp, err := get("/pods/" + podID + "/domains")
		if err != nil {
			return msg.Error{Err: err}
		}
		defer resp.Body.Close()

		var domains []PodDomain
		if err := json.NewDecoder(resp.Body).Decode(&domains); err != nil {
			return msg.Error{Err: err}
		}

		return msg.PodDomainsLoaded{Domains: domains}
	}
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

		return msg.PodDomainCreated{}
	}
}

func DeletePodDomain(podID, domainID string) tea.Cmd {
	return func() tea.Msg {
		resp, err := del("/pods/" + podID + "/domains/" + domainID)
		if err != nil {
			return msg.Error{Err: err}
		}
		defer resp.Body.Close()

		return msg.PodDomainDeleted{}
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

		return msg.PodDomainUpdated{}
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

		return msg.PodDomainCreated{}
	}
}

// --- Pod Env Vars ---

func FetchPodEnvVars(podID string) tea.Cmd {
	return func() tea.Msg {
		resp, err := get("/pods/" + podID + "/vars")
		if err != nil {
			return msg.Error{Err: err}
		}
		defer resp.Body.Close()

		var envVars []model.PodEnvVar
		err = json.NewDecoder(resp.Body).Decode(&envVars)
		if err != nil {
			return msg.Error{Err: err}
		}

		return msg.PodEnvVarsLoaded{EnvVars: envVars}
	}
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

		return msg.PodEnvVarsUpdated{}
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

		_, err = get("/health")
		if err != nil {
			if err == errs.ErrUnauthorized {
				return msg.ConnectionResult{NeedsAuth: true}
			}
			return msg.ConnectionResult{Offline: true}
		}

		return msg.ConnectionResult{Online: true}
	}
}
