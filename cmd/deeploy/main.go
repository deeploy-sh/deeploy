package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/deeploy-sh/deeploy/internal/deeploy/ui/pages"
)

func main() {
	// Logging Setup
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	// Start App
	m := pages.NewApp()
	p := tea.NewProgram(m)
	_, err := p.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
