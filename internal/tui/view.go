package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	/* ================= TEXT ================= */

	titleStyle = lipgloss.NewStyle().
			Bold(true)

	activeText = lipgloss.NewStyle().
			Foreground(lipgloss.Color("231")).
			Bold(true)

	inactiveText = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	/* ================= OUTER APP BORDER ================= */

	appBorder = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color("81")) // cyan

	appTitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("81")).
			Bold(true)

	/* ================= PANE BORDERS ================= */

	inactivePaneBorder = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("250")) // white

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
	var inner string

	// ---------- INPUT MODE ----------
	if m.InputActive {
		label := "Input"
		if m.InputType == InputNewProject {
			label = "New Project"
		}
		if m.InputType == InputNewTask {
			label = "New Task"
		}
		if m.InputType == InputRenameTask {
			label = "Rename Task"
		}

		inner = titleStyle.Render(label) + "\n\n" +
			m.InputValue + "\n\n" +
			inactiveText.Render("Enter = save • Esc = cancel")
	} else {
		switch m.ActivePane {
		case PaneProjects:
			inner = renderProjectsPane(m)
		case PaneBoard:
			inner = renderBoard(m)
		default:
			inner = "unknown state"
		}
	}

	// ---------- RESPONSIVE FRAME ----------
	width := m.Width
	height := m.Height
	if width < 60 {
		width = 60
	}
	if height < 20 {
		height = 20
	}

	// Render bordered app
	framed := appBorder.
		Width(width-2).
		Height(height-2).
		Padding(1, 2).
		Render(inner)

	// Overlay title ON the top border
	title := appTitleStyle.Render(" Chromaboard ")

	titleLine := lipgloss.Place(
		width-2,
		1,
		lipgloss.Center,
		lipgloss.Top,
		title,
	)

	return lipgloss.JoinVertical(
		lipgloss.Top,
		titleLine,
		framed,
	)
}

/* ================= PROJECT SELECTION ================= */

func renderProjectsPane(m Model) string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Projects") + "\n\n")

	if len(m.Projects) == 0 {
		b.WriteString(inactiveText.Render("No projects yet\n"))
		b.WriteString(inactiveText.Render("Press 'n' to create one"))
	} else {
		for i, name := range m.Projects {
			if i == m.ProjectIndex {
				b.WriteString(activeText.Render("▶ " + name))
			} else {
				b.WriteString(inactiveText.Render("  " + name))
			}
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(inactiveText.Render("n: new • enter: open • q: quit"))

	return inactivePaneBorder.
		Width(m.Width-8).
		Height(m.Height-8).
		Padding(1, 2).
		Render(b.String())
}

/* ================= BOARD ================= */

func renderBoard(m Model) string {
	colWidth := (m.Width - 16) / 3
	if colWidth < 28 {
		colWidth = 28
	}

	todo := renderColumn(m, ColumnTodo, "TODO", colWidth)
	pending := renderColumn(m, ColumnPending, "Pending", colWidth)
	done := renderColumn(m, ColumnDone, "Done", colWidth)

	board := lipgloss.JoinHorizontal(lipgloss.Top, todo, pending, done)

	footer := inactiveText.Render(
		"\nh/l: column • j/k: move • a: add • r: rename • m: move • d: delete • esc: back",
	)

	return titleStyle.Render(m.CurrentProject.Name) +
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

	minHeight := m.Height - 12
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
		style = inactivePaneBorder
	}

	return style.
		Width(width).
		Height(minHeight).
		Padding(1, 2).
		Render(b.String())
}
