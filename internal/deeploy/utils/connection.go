package utils

import (
	"errors"

	tea "charm.land/bubbletea/v2"
	"github.com/deeploy-sh/deeploy/internal/deeploy/config"
	"github.com/deeploy-sh/deeploy/internal/shared/errs"
)

type ConnectionResultMsg struct {
	NeedsSetup bool
	NeedsAuth  bool
	Offline    bool
	Online     bool
}

// TODO: Connection check improvements for post-MVP:
// - Validate that server is actually a Deeploy instance (check response body)
// - Differentiate "Offline" reasons (DNS error, server down, wrong URL)
func CheckConnection() tea.Msg {
	// Config check
	cfg, err := config.Load()
	if err != nil || cfg == nil || cfg.Server == "" {
		return ConnectionResultMsg{NeedsSetup: true}
	}
	if cfg.Token == "" {
		return ConnectionResultMsg{NeedsAuth: true}
	}

	// Server check
	_, err = Request("GET", "/health", nil)
	if err != nil {
		if errors.Is(err, errs.ErrUnauthorized) {
			return ConnectionResultMsg{NeedsAuth: true}
		}
		return ConnectionResultMsg{Offline: true}
	}

	return ConnectionResultMsg{Online: true}
}
