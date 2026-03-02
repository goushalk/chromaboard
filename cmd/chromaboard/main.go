package main

import (
	"fmt"
	"os"
	"time"

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

	// Load project list with metadata
	summaries, err := storage.ListProjectSummaries()
	if err != nil {
		fmt.Println("failed to load projects:", err)
		os.Exit(1)
	}

	projects := make([]string, 0, len(summaries))
	projectDates := make(map[string]time.Time, len(summaries))
	for _, summary := range summaries {
		projects = append(projects, summary.Name)
		projectDates[summary.Name] = summary.CreatedAt
	}

	// Initialize TUI model
	model := tui.Model{
		ActivePane:   tui.PaneProjects,
		Projects:     projects,
		ProjectIndex: 0,
		ProjectDates: projectDates,
	}

	// Start Bubble Tea program
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("error running program:", err)
		os.Exit(1)
	}
}
