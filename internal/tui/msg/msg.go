package msg

import (
	tea "charm.land/bubbletea/v2"
	"github.com/deeploy-sh/deeploy/internal/server/repo"
)

// --- Navigation ---

type ChangePage struct {
	PageFactory func(s Store) tea.Model
}

// Store interface for page factories
type Store interface {
	Projects() []repo.Project
	Pods() []repo.Pod
}

// --- Connection ---

type ConnectionResult struct {
	NeedsSetup bool
	NeedsAuth  bool
	Offline    bool
	Online     bool
}

// --- Auth ---

type AuthError struct{ Err error }
type AuthSuccess struct{}

// --- Data Loaded ---

type DataLoaded struct {
	Projects []repo.Project
	Pods     []repo.Pod
}

type ProjectsLoaded struct{ Projects []repo.Project }
type PodsLoaded struct{ Pods []repo.Pod }

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

type GitTokensLoaded struct{ Tokens any }
type GitTokenCreated struct{}
type GitTokenDeleted struct{}

// --- Pod Domains ---

type PodDomainsLoaded struct{ Domains any }
type PodDomainCreated struct{}
type PodDomainDeleted struct{}

// --- Errors ---

type Error struct{ Err error }

// --- Theme ---

type ThemeSwitcherClose struct{ Theme string }
type OpenThemeSwitcher struct{}
