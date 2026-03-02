package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/goushalk/chromaboard/internal/domain"
)

const (
	MinWidth  = 80
	MinHeight = 24
)

var (
	titleStyle = lipgloss.NewStyle().Bold(true)

	activeText = lipgloss.NewStyle().
			Foreground(lipgloss.Color("231")).
			Bold(true)

	inactiveText = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	doneText = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Strikethrough(true)

	doneActiveText = lipgloss.NewStyle().
			Foreground(lipgloss.Color("231")).
			Bold(true).
			Strikethrough(true)

	appBorder = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color("81"))

	helpBorder = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("81")).
			Padding(1, 2)

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

func (m Model) View() string {
	// ---------- SIZE GUARD ----------
	if m.Width < MinWidth || m.Height < MinHeight {
		return renderTooSmall(m)
	}

	// ---------- HELP OVERLAY ----------
	if m.ShowHelp {
		return renderHelp(m)
	}

	var inner string

	// ---------- INPUT MODE ----------
	if m.InputActive {
		label := "Input"
		if m.InputType == InputNewProject {
			label = "New Project"
		}
		if m.InputType == InputNewTask {
			label = "New Task (title | description)"
		}
		if m.InputType == InputRenameTask {
			label = "Rename Task"
		}
		if m.InputType == InputTaskDescription {
			label = "Task Description"
		}
		if m.InputType == InputNewSubSection {
			label = "New Sub Section"
		}
		if m.InputType == InputNewSubTask {
			label = "New Sub Task (section | subtask)"
		}
		if m.InputType == InputToggleSubTask {
			label = "Toggle Sub Task (section | subtask)"
		}

		inner = titleStyle.Render(label) + "\n\n" +
			m.InputValue + "\n\n" +
			inactiveText.Render("Enter = save • Esc = cancel")
		if m.Status != "" {
			inner += "\n" + inactiveText.Render(m.Status)
		}

	} else {
		switch m.ActivePane {
		case PaneProjects:
			inner = renderProjectsPane(m)
		case PaneBoard:
			inner = renderBoard(m)
		case PaneTaskDetail:
			inner = renderTaskDetail(m)
		}
	}

	framed := appBorder.
		Width(m.Width-2).
		Height(m.Height-2).
		Padding(1, 2).
		Render(inner)

	title := lipgloss.Place(
		m.Width-2,
		1,
		lipgloss.Center,
		lipgloss.Top,
		titleStyle.Render(" Chromaboard "),
	)

	return lipgloss.JoinVertical(lipgloss.Top, title, framed)
}

/* ================= HELP ================= */

func renderHelp(m Model) string {
	content := `
Navigation
  j / k        Move up / down
  h / l        Switch columns
  tab          Switch pane

Actions
  n            New project
  d            Delete project (in Projects pane)
  a            Add task (title | description)
  r            Rename task
  d            Delete task (in Board pane)
  m / M        Move task right / left
  enter        Open selected task
  e            Edit task description (Task Detail)
  s            Add sub section (Task Detail)
  t            Add sub task (Task Detail)
  x            Toggle sub task done (Task Detail)

General
  ?            Toggle help
  esc          Back / Close
  q            Quit
`

	box := helpBorder.Render(
		titleStyle.Render(" Help ") + "\n" +
			inactiveText.Render(strings.TrimSpace(content)),
	)

	return lipgloss.Place(
		m.Width,
		m.Height,
		lipgloss.Center,
		lipgloss.Center,
		box,
	)
}

/* ================= TOO SMALL ================= */

func renderTooSmall(m Model) string {
	msg := fmt.Sprintf(
		"Terminal too small\n\nRequired:\n  width  ≥ %d\n  height ≥ %d\n\nCurrent:\n  width  = %d\n  height = %d",
		MinWidth, MinHeight, m.Width, m.Height,
	)

	box := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("196")).
		Padding(1, 2).
		Render(msg)

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

	for i, name := range m.Projects {
		line := "  " + name
		if createdAt, ok := m.ProjectDates[name]; ok {
			line += "  (" + formatCreatedAt(createdAt) + ")"
		}

		if i == m.ProjectIndex {
			b.WriteString(activeText.Render("▶ " + strings.TrimPrefix(line, "  ")))
		} else {
			b.WriteString(inactiveText.Render(line))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(inactiveText.Render("n: new • d: delete • enter: open • ?: help"))

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
		"\nh/l: column • j/k: move • enter: open task • a: add • m/M: move task • ?: help",
	)

	return board + "\n" + footer
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
		if t.Description != "" || len(t.SubSections) > 0 {
			line += " *"
		}

		if t.Status == domain.StatusDone {
			if m.ActiveColumn == col && index == m.TaskIndex {
				b.WriteString(doneActiveText.Render(line))
			} else {
				b.WriteString(doneText.Render(line))
			}
		} else {
			if m.ActiveColumn == col && index == m.TaskIndex {
				b.WriteString(activeText.Render(line))
			} else {
				b.WriteString(inactiveText.Render(line))
			}
		}
		b.WriteString("\n")
		index++
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
		Height(m.Height-12).
		Padding(1, 2).
		Render(b.String())
}

func renderTaskDetail(m Model) string {
	task, ok := selectedTask(m)
	if !ok {
		return inactiveText.Render("Task not found")
	}

	var b strings.Builder
	b.WriteString(titleStyle.Render("Task Detail") + "\n\n")
	if task.Status == domain.StatusDone {
		b.WriteString(doneActiveText.Render(task.Title) + "\n")
	} else {
		b.WriteString(activeText.Render(task.Title) + "\n")
	}
	b.WriteString(inactiveText.Render("Status: "+string(task.Status)) + "\n\n")
	b.WriteString(inactiveText.Render("Created: "+formatCreatedAt(task.CreatedAt)) + "\n\n")

	if task.Description == "" {
		b.WriteString(inactiveText.Render("Description: (empty)") + "\n\n")
	} else {
		b.WriteString(inactiveText.Render("Description:") + "\n")
		b.WriteString(task.Description + "\n\n")
	}

	b.WriteString(inactiveText.Render("Sub Sections") + "\n")
	if len(task.SubSections) == 0 {
		b.WriteString(inactiveText.Render("  (none)") + "\n")
	} else {
		for _, section := range task.SubSections {
			b.WriteString("• " + section.Title + "\n")
			if len(section.SubTasks) == 0 {
				b.WriteString(inactiveText.Render("  - (no sub tasks)") + "\n")
				continue
			}
			for _, subTask := range section.SubTasks {
				mark := "[ ]"
				if subTask.Done {
					mark = "[x]"
				}
				b.WriteString("  - " + mark + " " + subTask.Title + "\n")
			}
		}
	}

	b.WriteString("\n")
	b.WriteString(inactiveText.Render("e: description • s: sub section • t: sub task • x: toggle sub task • esc: back"))
	if m.Status != "" {
		b.WriteString("\n" + inactiveText.Render("Status: "+m.Status))
	}

	return inactivePaneBorder.
		Width(m.Width-8).
		Height(m.Height-8).
		Padding(1, 2).
		Render(b.String())
}

func formatCreatedAt(t time.Time) string {
	if t.IsZero() {
		return "unknown"
	}
	return t.Local().Format("2006-01-02 15:04:05")
}
