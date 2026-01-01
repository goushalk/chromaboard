package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

/* ================= CONSTANTS ================= */

const (
	MinWidth  = 80
	MinHeight = 24
)

/* ================= STYLES ================= */

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true)

	activeText = lipgloss.NewStyle().
			Foreground(lipgloss.Color("231")).
			Bold(true)

	inactiveText = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	warnTitle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true)

	warnText = lipgloss.NewStyle().
			Foreground(lipgloss.Color("250"))

	appBorder = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color("81"))

	appTitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("81")).
			Bold(true)

	inactivePaneBorder = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("250"))

	todoBorderActive = lipgloss.NewStyle().
				Border(lipgloss.ThickBorder()).
				BorderForeground(lipgloss.Color("196"))

	pendingBorderActive = lipgloss.NewStyle().
				Border(lipgloss.ThickBorder()).
				BorderForeground(lipgloss.Color("220"))

	doneBorderActive = lipgloss.NewStyle().
				Border(lipgloss.ThickBorder()).
				BorderForeground(lipgloss.Color("46"))
)

/* ================= VIEW ================= */

func (m Model) View() string {
	// ---------- TERMINAL SIZE GUARD ----------
	if m.Width < MinWidth || m.Height < MinHeight {
		return renderTooSmall(m)
	}

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

	// ---------- APP FRAME ----------
	framed := appBorder.
		Width(m.Width-2).
		Height(m.Height-2).
		Padding(1, 2).
		Render(inner)

	title := appTitleStyle.Render(" Chromaboard ")

	titleLine := lipgloss.Place(
		m.Width-2,
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

/* ================= TOO SMALL VIEW ================= */

func renderTooSmall(m Model) string {
	message := fmt.Sprintf(
		"Terminal too small\n\nRequired:\n  width  ≥ %d\n  height ≥ %d\n\nCurrent:\n  width  = %d\n  height = %d",
		MinWidth,
		MinHeight,
		m.Width,
		m.Height,
	)

	box := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("196")).
		Padding(1, 2).
		Render(
			warnTitle.Render(" Resize Required ") + "\n\n" +
				warnText.Render(message),
		)

	return lipgloss.Place(
		m.Width,
		m.Height,
		lipgloss.Center,
		lipgloss.Center,
		box,
	)
}

/* ================= PROJECTS ================= */

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
