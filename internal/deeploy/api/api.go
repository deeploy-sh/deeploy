package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	tea "charm.land/bubbletea/v2"
	"github.com/deeploy-sh/deeploy/internal/deeploy/config"
	"github.com/deeploy-sh/deeploy/internal/deeploy/msg"
	"github.com/deeploy-sh/deeploy/internal/deeployd/repo"
	"github.com/deeploy-sh/deeploy/internal/shared/errs"
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
