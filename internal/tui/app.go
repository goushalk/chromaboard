package tui

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/goushalk/chromaboard/internal/domain"
)

/*
High-level panes
*/
type Pane int

const (
	PaneProjects Pane = iota
	PaneBoard
)

/*
Board columns (UI concept)
*/
type Column int

const (
	ColumnTodo Column = iota
	ColumnPending
	ColumnDone
)

/*
Input mode purpose
*/
type InputMode int

const (
	InputNone InputMode = iota
	InputNewProject
	InputNewTask
	InputRenameTask
)

/*
Model represents the entire UI state
*/
type Model struct {
	// Terminal size
	Width  int
	Height int

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
