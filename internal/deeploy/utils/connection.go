package utils

import (
	"errors"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/deeploy-sh/deeploy/internal/deeploy/config"
	"github.com/deeploy-sh/deeploy/internal/shared/errs"
)

type ConnectionResultMsg struct {
	NeedsSetup bool
	Offline    bool
	NeedsAuth  bool
}

func CheckConnection() tea.Msg {
	// Config check
	cfg, err := config.Load()
	if err != nil || cfg == nil || cfg.Server == "" || cfg.Token == "" {
		return ConnectionResultMsg{NeedsSetup: true}
	}

	// Server check
	_, err = Request("GET", "/health", nil)
	if err != nil {
		if errors.Is(err, errs.ErrUnauthorized) {
			return ConnectionResultMsg{NeedsAuth: true}
		}
		return ConnectionResultMsg{Offline: true}
	}

	return ConnectionResultMsg{Offline: false}
}
