package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/goushalk/chromaboard/internal/storage"
	"github.com/goushalk/chromaboard/internal/tui"
)

func main() {
	// Ensure storage directories exist
	if err := storage.EnsureStorage(); err != nil {
		fmt.Println("failed to initialize storage:", err)
		os.Exit(1)
	}

	// Load project list (filenames without .json)
	projects, err := storage.ListProjects()
	if err != nil {
		fmt.Println("failed to load projects:", err)
		os.Exit(1)
	}

	// Initialize TUI model
	model := tui.Model{
		ActivePane:   tui.PaneProjects,
		Projects:     projects,
		ProjectIndex: 0,
	}

	// Start Bubble Tea program
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("error running program:", err)
		os.Exit(1)
	}
}
