package utils

import (
	"os/exec"
	"runtime"

	tea "charm.land/bubbletea/v2"
)

// OpenBrowser opens the given URL in the system's default browser.
func OpenBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default:
		cmd = "xdg-open"
	}

	return exec.Command(cmd, append(args, url)...).Start()
}

// OpenBrowserCmd returns a tea.Cmd that opens the URL in the browser.
func OpenBrowserCmd(url string) tea.Cmd {
	return func() tea.Msg {
		OpenBrowser(url)
		return nil
	}
}
