package tui

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/goushalk/chromaboard/internal/domain"
)

type Pane int

const (
	PaneProjects Pane = iota
	PaneBoard
)

type Column int

const (
	ColumnTodo Column = iota
	ColumnPending
	ColumnDone
)

// INPUT PURPOSE
type InputMode int

const (
	InputNone InputMode = iota
	InputNewProject
	InputNewTask
)

type Model struct {
	// Navigation
	ActivePane Pane

	// Projects pane
	Projects     []string
	ProjectIndex int

	// Board
	CurrentProject *domain.Project
	ActiveColumn   Column
	TaskIndex      int

	// Input mode
	InputActive bool
	InputValue  string
	InputType   InputMode

	// UI feedback
	Error  error
	Status string
}

func (m Model) Init() tea.Cmd {
	return nil
}
