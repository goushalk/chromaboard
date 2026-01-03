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

		// ================= HELP MODE =================
		if m.ShowHelp {
			switch msg.String() {
			case "?", "esc":
				m.ShowHelp = false
			}
			return m, nil
		}

		// Toggle help (global)
		if msg.String() == "?" {
			m.ShowHelp = true
			return m, nil
		}

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
				if ok {
					_ = m.CurrentProject.MoveTask(taskID, columnToNextStatus(m.ActiveColumn))
					_ = storage.SaveRegistry(*m.CurrentProject)
				}

			case "d":
				taskID, ok := selectedTaskID(m)
				if ok {
					_ = m.CurrentProject.DeleteTask(taskID)
					_ = storage.SaveRegistry(*m.CurrentProject)
					if m.TaskIndex > 0 {
						m.TaskIndex--
					}
				}
			}
		}
	}

	return m, nil
}
