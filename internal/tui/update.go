package tui

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/goushalk/chromaboard/internal/domain"
	"github.com/goushalk/chromaboard/internal/storage"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil

	case tea.KeyMsg:

		// ================= INPUT MODE =================
		if m.InputActive {
			switch msg.String() {

			case "esc":
				m.InputActive = false
				m.InputValue = ""
				m.InputType = InputNone
				return m, nil

			case "enter":
				switch m.InputType {

				case InputNewProject:
					project := domain.NewProject(m.InputValue)
					_ = storage.SaveRegistry(project)

					projects, _ := storage.ListProjects()
					m.Projects = projects
					m.ProjectIndex = len(projects) - 1

				case InputNewTask:
					if m.CurrentProject != nil {
						m.CurrentProject.AddTask(m.InputValue)
						_ = storage.SaveRegistry(*m.CurrentProject)
					}

				case InputRenameTask:
					if m.CurrentProject != nil {
						taskID, ok := selectedTaskID(m)
						if ok {
							_ = m.CurrentProject.RenameTask(taskID, m.InputValue)
							_ = storage.SaveRegistry(*m.CurrentProject)
						}
					}
				}

				m.InputActive = false
				m.InputValue = ""
				m.InputType = InputNone
				return m, nil

			case "backspace":
				if len(m.InputValue) > 0 {
					m.InputValue = m.InputValue[:len(m.InputValue)-1]
				}
				return m, nil

			default:
				if len(msg.String()) == 1 {
					m.InputValue += msg.String()
				}
				return m, nil
			}
		}

		// ================= GLOBAL =================
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}

		// ================= PANES =================
		switch m.ActivePane {

		// ---------- PROJECTS ----------
		case PaneProjects:
			switch msg.String() {

			case "j", "down":
				if m.ProjectIndex < len(m.Projects)-1 {
					m.ProjectIndex++
				}

			case "k", "up":
				if m.ProjectIndex > 0 {
					m.ProjectIndex--
				}

			case "enter":
				if len(m.Projects) == 0 {
					return m, nil
				}

				name := m.Projects[m.ProjectIndex]
				project, err := storage.LoadRegistry(name)
				if err != nil {
					m.Error = err
					return m, nil
				}

				m.CurrentProject = &project
				m.ActivePane = PaneBoard
				m.ActiveColumn = ColumnTodo
				m.TaskIndex = 0

			case "n":
				m.InputActive = true
				m.InputType = InputNewProject
				m.InputValue = ""
			}

		// ---------- BOARD ----------
		case PaneBoard:
			if m.CurrentProject == nil {
				m.ActivePane = PaneProjects
				return m, nil
			}

			switch msg.String() {

			case "esc":
				m.ActivePane = PaneProjects
				m.CurrentProject = nil
				m.TaskIndex = 0

			case "h", "left":
				if m.ActiveColumn > ColumnTodo {
					m.ActiveColumn--
					m.TaskIndex = 0
				}

			case "l", "right":
				if m.ActiveColumn < ColumnDone {
					m.ActiveColumn++
					m.TaskIndex = 0
				}

			case "j", "down":
				if m.TaskIndex < countTasksInColumn(m)-1 {
					m.TaskIndex++
				}

			case "k", "up":
				if m.TaskIndex > 0 {
					m.TaskIndex--
				}

			case "a":
				m.InputActive = true
				m.InputType = InputNewTask
				m.InputValue = ""

			case "r":
				m.InputActive = true
				m.InputType = InputRenameTask
				m.InputValue = ""

			case "m":
				taskID, ok := selectedTaskID(m)
				if !ok {
					return m, nil
				}

				next := columnToNextStatus(m.ActiveColumn)
				if err := m.CurrentProject.MoveTask(taskID, next); err != nil {
					m.Error = err
					return m, nil
				}

				_ = storage.SaveRegistry(*m.CurrentProject)

			case "d":
				taskID, ok := selectedTaskID(m)
				if !ok {
					return m, nil
				}

				if err := m.CurrentProject.DeleteTask(taskID); err != nil {
					m.Error = err
					return m, nil
				}

				_ = storage.SaveRegistry(*m.CurrentProject)

				if m.TaskIndex > 0 {
					m.TaskIndex--
				}
			}
		}
	}

	return m, nil
}

/* ---------- HELPERS ---------- */

func countTasksInColumn(m Model) int {
	count := 0
	for _, t := range m.CurrentProject.Tasks {
		if statusToColumn(t.Status) == m.ActiveColumn {
			count++
		}
	}
	return count
}

func selectedTaskID(m Model) (int, bool) {
	index := 0
	for _, t := range m.CurrentProject.Tasks {
		if statusToColumn(t.Status) == m.ActiveColumn {
			if index == m.TaskIndex {
				return t.ID, true
			}
			index++
		}
	}
	return 0, false
}

func statusToColumn(s domain.Status) Column {
	switch s {
	case domain.StatusTodo:
		return ColumnTodo
	case domain.StatusPending:
		return ColumnPending
	case domain.StatusDone:
		return ColumnDone
	default:
		return ColumnTodo
	}
}

func columnToNextStatus(c Column) domain.Status {
	switch c {
	case ColumnTodo:
		return domain.StatusPending
	case ColumnPending:
		return domain.StatusDone
	default:
		return domain.StatusDone
	}
}
