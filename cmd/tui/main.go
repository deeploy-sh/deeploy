package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/deeploy-sh/deeploy/internal/shared/version"
	"github.com/deeploy-sh/deeploy/internal/tui/ui/page"
)

func main() {
	// Handle --version flag
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Printf("deeploy %s\n", version.Version)
		os.Exit(0)
	}

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
	m := page.NewApp()
	p := tea.NewProgram(m)
	_, err := p.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
