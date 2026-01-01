package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true)

	activeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true)

	inactiveStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	columnStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			Padding(0, 1).
			Width(24)
)

func (m Model) View() string {
	// INPUT MODE VIEW
	if m.InputActive {
		label := "Input:"
		if m.InputType == InputNewProject {
			label = "New Project:"
		}
		if m.InputType == InputNewTask {
			label = "New Task:"
		}

		return "\n" +
			titleStyle.Render(label) + "\n\n" +
			m.InputValue + "\n\n" +
			inactiveStyle.Render("Enter = save | Esc = cancel")
	}

	switch m.ActivePane {
	case PaneProjects:
		return renderProjects(m)
	case PaneBoard:
		return renderBoard(m)
	default:
		return "unknown state"
	}
}

/*
PROJECTS VIEW
*/
func renderProjects(m Model) string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Projects\n\n"))

	for i, name := range m.Projects {
		if i == m.ProjectIndex {
			b.WriteString(activeStyle.Render("> " + name))
		} else {
			b.WriteString(inactiveStyle.Render("  " + name))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(inactiveStyle.Render("n: new  enter: open  q: quit"))

	return b.String()
}

/*
BOARD VIEW
*/
func renderBoard(m Model) string {
	if m.CurrentProject == nil {
		return "no project loaded"
	}

	todo := renderColumn(m, ColumnTodo, "TODO")
	pending := renderColumn(m, ColumnPending, "Pending")
	done := renderColumn(m, ColumnDone, "Done")

	board := lipgloss.JoinHorizontal(lipgloss.Top, todo, pending, done)

	footer := inactiveStyle.Render(
		"\nh/l: column  j/k: move  a: add  m: move  d: delete  esc: back",
	)

	return titleStyle.Render(m.CurrentProject.Name) + "\n\n" + board + footer
}

func renderColumn(m Model, col Column, title string) string {
	var b strings.Builder

	if m.ActiveColumn == col {
		b.WriteString(activeStyle.Render(title) + "\n")
	} else {
		b.WriteString(inactiveStyle.Render(title) + "\n")
	}

	index := 0
	for _, t := range m.CurrentProject.Tasks {
		if statusToColumn(t.Status) != col {
			continue
		}

		line := "• " + t.Title
		if m.ActiveColumn == col && index == m.TaskIndex {
			b.WriteString(activeStyle.Render(line))
		} else {
			b.WriteString(inactiveStyle.Render(line))
		}
		b.WriteString("\n")
		index++
	}

	if index == 0 {
		b.WriteString(inactiveStyle.Render("(empty)\n"))
	}

	return columnStyle.Render(b.String())
}
