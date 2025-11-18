package messages

import (
	"github.com/deeploy-sh/deeploy/internal/shared/repo"
	tea "github.com/charmbracelet/bubbletea"
)

// Navigation Messages

type ChangePageMsg struct {
	Page tea.Model
}

// Auth Messages
type AuthErrorMsg struct {
	Err error
}
type AuthSuccessMsg struct {
}

// Project Messages
type ProjectCreatedMsg repo.ProjectDTO
type ProjectUpdatedMsg repo.ProjectDTO
type ProjectDeleteMsg *repo.ProjectDTO
type ProjectErrMsg error
type ProjectsInitDataMsg []repo.ProjectDTO
type ProjectPushPageMsg struct{ Page tea.Model }
type ProjectPopPageMsg struct{}
