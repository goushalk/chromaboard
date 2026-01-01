package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	// Text styles
	titleStyle = lipgloss.NewStyle().
			Bold(true)

	activeText = lipgloss.NewStyle().
			Foreground(lipgloss.Color("231")). // white
			Bold(true)

	inactiveText = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	// Borders
	inactiveBorder = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("250")) // white

	activeBorder = lipgloss.NewStyle().
			Border(lipgloss.ThickBorder()).
			BorderForeground(lipgloss.Color("81")) // cyan accent

	// Board active borders
	todoBorderActive = lipgloss.NewStyle().
				Border(lipgloss.ThickBorder()).
				BorderForeground(lipgloss.Color("196")) // red

	pendingBorderActive = lipgloss.NewStyle().
				Border(lipgloss.ThickBorder()).
				BorderForeground(lipgloss.Color("220")) // yellow

	doneBorderActive = lipgloss.NewStyle().
				Border(lipgloss.ThickBorder()).
				BorderForeground(lipgloss.Color("46")) // green
)

func (m Model) View() string {
	// ---------- INPUT MODE ----------
	if m.InputActive {
		label := "Input:"
		if m.InputType == InputNewProject {
			label = "New Project"
		}
		if m.InputType == InputNewTask {
			label = "New Task"
		}
		if m.InputType == InputRenameTask {
			label = "Rename Task"
		}

		return "\n" +
			titleStyle.Render(label) + "\n\n" +
			m.InputValue + "\n\n" +
			inactiveText.Render("Enter = save • Esc = cancel")
	}

	switch m.ActivePane {
	case PaneProjects:
		return renderProjectsPane(m)
	case PaneBoard:
		return renderBoard(m)
	default:
		return "unknown state"
	}
}

/* ================= PROJECT SELECTION ================= */

func renderProjectsPane(m Model) string {
	width := m.Width - 4
	if width < 40 {
		width = 40
	}

	height := m.Height - 4
	if height < 10 {
		height = 10
	}

	var b strings.Builder

	b.WriteString(titleStyle.Render(" Projects "))
	b.WriteString("\n\n")

	if len(m.Projects) == 0 {
		b.WriteString(inactiveText.Render("No projects yet\n"))
		b.WriteString(inactiveText.Render("Press 'n' to create one"))
	} else {
		for i, name := range m.Projects {
			line := "  " + name
			if i == m.ProjectIndex {
				line = "▶ " + name
				b.WriteString(activeText.Render(line))
			} else {
				b.WriteString(inactiveText.Render(line))
			}
			b.WriteString("\n")
		}
	}

	content := b.String()

	footer := inactiveText.Render(
		"\n\nn: new project  enter: open  q: quit",
	)

	return activeBorder.
		Width(width).
		Height(height).
		Padding(1, 2).
		Render(content + footer)
}

/* ================= BOARD ================= */

func renderBoard(m Model) string {
	if m.CurrentProject == nil {
		return "no project loaded"
	}

	colWidth := (m.Width - 6) / 3
	if colWidth < 30 {
		colWidth = 30
	}

	todo := renderColumn(m, ColumnTodo, "TODO", colWidth)
	pending := renderColumn(m, ColumnPending, "Pending", colWidth)
	done := renderColumn(m, ColumnDone, "Done", colWidth)

	board := lipgloss.JoinHorizontal(lipgloss.Top, todo, pending, done)

	footer := inactiveText.Render(
		"\nh/l: column  j/k: move  a: add  r: rename  m: move  d: delete  esc: back",
	)

	return titleStyle.Render(" "+m.CurrentProject.Name+" ") +
		"\n\n" + board + "\n" + footer
}

/* ================= COLUMN ================= */

func renderColumn(m Model, col Column, title string, width int) string {
	var b strings.Builder

	if m.ActiveColumn == col {
		b.WriteString(activeText.Render(title) + "\n\n")
	} else {
		b.WriteString(inactiveText.Render(title) + "\n\n")
	}

	index := 0
	for _, t := range m.CurrentProject.Tasks {
		if statusToColumn(t.Status) != col {
			continue
		}

		line := "• " + t.Title
		if m.ActiveColumn == col && index == m.TaskIndex {
			b.WriteString(activeText.Render(line))
		} else {
			b.WriteString(inactiveText.Render(line))
		}
		b.WriteString("\n")
		index++
	}

	if index == 0 {
		b.WriteString(inactiveText.Render("(empty)\n"))
	}

	minHeight := m.Height - 8
	if minHeight < 10 {
		minHeight = 10
	}

	var style lipgloss.Style
	if m.ActiveColumn == col {
		switch col {
		case ColumnTodo:
			style = todoBorderActive
		case ColumnPending:
			style = pendingBorderActive
		case ColumnDone:
			style = doneBorderActive
		}
	} else {
		style = inactiveBorder
	}

	return style.
		Width(width).
		Height(minHeight).
		Padding(1, 2).
		Render(b.String())
}
