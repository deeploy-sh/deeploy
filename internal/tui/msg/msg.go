package msg

import (
	tea "charm.land/bubbletea/v2"
	"github.com/deeploy-sh/deeploy/internal/shared/model"
)

// --- Navigation ---

type ChangePage struct {
	PageFactory func(s Store) tea.Model
}

// Store interface for page factories
type Store interface {
	Projects() []model.Project
	Pods() []model.Pod
	GitTokens() []model.GitToken
}

// --- Connection ---

type ConnectionResult struct {
	NeedsSetup    bool
	NeedsAuth     bool
	Offline       bool
	Online        bool
	ServerVersion string // Server version from /health
}

// --- Version Check ---

type LatestVersionResult struct {
	Version string
	Error   error
}

// --- Auth ---

type AuthError struct{ Err error }
type AuthSuccess struct{}

// --- Data Loaded ---

type DataLoaded struct {
	Projects  []model.Project
	Pods      []model.Pod
	GitTokens []model.GitToken
}

type ProjectsLoaded struct{ Projects []model.Project }
type PodsLoaded struct{ Pods []model.Pod }

// --- CRUD Success (trigger reload) ---

type ProjectCreated struct{}
type ProjectUpdated struct{}
type ProjectDeleted struct{}
type PodCreated struct{}
type PodUpdated struct{}
type PodDeleted struct{}

// --- Pod Deploy ---

type PodDeployed struct{}
type PodStopped struct{}
type PodRestarted struct{}
type PodLogsLoaded struct{ Logs []string }

// --- Git Tokens ---

type GitTokenCreated struct{}
type GitTokenDeleted struct{}

// --- Pod Domains ---

type PodDomainsLoaded struct{ Domains []model.PodDomain }
type PodDomainCreated struct{}
type PodDomainUpdated struct{}
type PodDomainDeleted struct{}

// --- Pod Env Vars ---

type PodEnvVarsLoaded struct{ EnvVars []model.PodEnvVar }
type PodEnvVarsUpdated struct{}

// --- Errors ---

type Error struct{ Err error }

// --- Status Line ---

type StatusType int

const (
	StatusSuccess StatusType = iota
	StatusError
	StatusInfo
)

type ShowStatus struct {
	Text string
	Type StatusType
}

type ClearStatus struct{}

// --- Theme ---

type ThemeSwitcherClose struct{ Theme string }
type OpenThemeSwitcher struct{}
