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
	PodDomains(podID string) []model.PodDomain
	PodEnvVars(podID string) []model.PodEnvVar
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
	Projects   []model.Project
	Pods       []model.Pod
	GitTokens  []model.GitToken
	PodDomains []model.PodDomain
	PodEnvVars []model.PodEnvVar
}

type ProjectsLoaded struct{ Projects []model.Project }
type PodsLoaded struct{ Pods []model.Pod }

// --- CRUD Success (optimistic updates) ---

type ProjectCreated struct{ Project model.Project }
type ProjectUpdated struct{ Project model.Project }
type ProjectDeleted struct{ ProjectID string }
type PodCreated struct{ Pod model.Pod }
type PodUpdated struct{ Pod model.Pod }
type PodDeleted struct {
	PodID     string
	ProjectID string
}

// --- Pod Deploy ---

type PodLoaded struct{ Pod model.Pod }
type PodDeployed struct{}
type PodStopped struct{}
type PodRestarted struct{}
type PodLogsLoaded struct{ Logs []string }

// --- Git Tokens ---

type GitTokenCreated struct{ Token model.GitToken }
type GitTokenDeleted struct{ TokenID string }

// --- Pod Domains ---

type PodDomainsLoaded struct{ Domains []model.PodDomain }
type PodDomainCreated struct{ Domain model.PodDomain }
type PodDomainUpdated struct{ Domain model.PodDomain }
type PodDomainDeleted struct {
	DomainID string
	PodID    string
}

// --- Pod Env Vars ---

type PodEnvVarsLoaded struct{ EnvVars []model.PodEnvVar }
type PodEnvVarsUpdated struct {
	PodID   string
	EnvVars []model.PodEnvVar
}

// --- Server Settings ---

type ServerDomainLoaded struct{ Domain string }
type ServerDomainSet struct{}
type ServerDomainDeleted struct{}

// --- Errors ---

type Error struct{ Err error }

// --- Loading State ---

type StartLoading struct {
	Text string
}

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
