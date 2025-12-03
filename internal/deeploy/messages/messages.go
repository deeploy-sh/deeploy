package messages

import (
	tea "charm.land/bubbletea/v2"
	"github.com/deeploy-sh/deeploy/internal/deeployd/repo"
)

// Navigation Messages

type ChangePageMsg struct {
	Page tea.Model
}

// Auth Messages
type AuthErrorMsg struct {
	Err error
}
type AuthSuccessMsg struct{}

// Project Messages
type ProjectCreatedMsg repo.Project
type ProjectUpdatedMsg repo.Project
type ProjectDeleteMsg *repo.Project
type ProjectErrMsg error
type ProjectsInitDataMsg []repo.Project

// Pod Messages
type PodCreatedMsg repo.Pod
type PodUpdatedMsg repo.Pod
type PodDeleteMsg *repo.Pod
type PodErrMsg error
type PodsInitDataMsg []repo.Pod
