package messages

import (
	"github.com/deeploy-sh/deeploy/internal/deeployd/repo"
)

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
