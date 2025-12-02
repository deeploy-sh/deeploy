package messages

import (
	"github.com/deeploy-sh/deeploy/internal/deeployd/repo"
	tea "charm.land/bubbletea/v2"
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
type ProjectCreatedMsg repo.ProjectDTO
type ProjectUpdatedMsg repo.ProjectDTO
type ProjectDeleteMsg *repo.ProjectDTO
type ProjectErrMsg error
type ProjectsInitDataMsg []repo.ProjectDTO

// Pod Messages
type PodCreatedMsg repo.PodDTO
type PodUpdatedMsg repo.PodDTO
type PodDeleteMsg *repo.PodDTO
type PodErrMsg error
type PodsInitDataMsg []repo.PodDTO
