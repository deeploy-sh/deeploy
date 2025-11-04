package messages

import (
	"github.com/axadrn/deeploy/internal/data"
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
	token string
}

// Project Messages
type ProjectCreatedMsg data.ProjectDTO
type ProjectUpdatedMsg data.ProjectDTO
type ProjectDeleteMsg *data.ProjectDTO
type ProjectErrMsg error
type ProjectsInitDataMsg []data.ProjectDTO
type ProjectPushPageMsg struct{ Page tea.Model }
type ProjectPopPageMsg struct{}
