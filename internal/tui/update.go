package tui

import (
	"strings"
	"time"

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
				shouldCloseInput := true

				switch m.InputType {

				case InputNewProject:
					project := domain.NewProject(m.InputValue)
					_ = storage.SaveRegistry(project)

					refreshProjects(&m)
					if len(m.Projects) > 0 {
						m.ProjectIndex = len(m.Projects) - 1
					} else {
						m.ProjectIndex = 0
					}

				case InputNewTask:
					if m.CurrentProject != nil {
						title := strings.TrimSpace(m.InputValue)
						description := ""
						if left, right, ok := parseTwoPartInput(m.InputValue); ok {
							title = left
							description = right
						}
						if title != "" {
							m.CurrentProject.AddTask(title, description)
							_ = storage.SaveRegistry(*m.CurrentProject)
						}
					}

				case InputRenameTask:
					if m.CurrentProject != nil {
						taskID, ok := selectedTaskID(m)
						if ok {
							_ = m.CurrentProject.RenameTask(taskID, m.InputValue)
							_ = storage.SaveRegistry(*m.CurrentProject)
						}
					}

				case InputTaskDescription:
					if m.CurrentProject != nil && m.SelectedTaskID != 0 {
						_ = m.CurrentProject.SetTaskDescription(m.SelectedTaskID, m.InputValue)
						_ = storage.SaveRegistry(*m.CurrentProject)
					}

				case InputNewSubSection:
					if m.CurrentProject != nil && m.SelectedTaskID != 0 {
						title := strings.TrimSpace(m.InputValue)
						if title != "" {
							_ = m.CurrentProject.AddSubSection(m.SelectedTaskID, title)
							_ = storage.SaveRegistry(*m.CurrentProject)
						}
					}

				case InputNewSubTask:
					if m.CurrentProject != nil && m.SelectedTaskID != 0 {
						section, subTask, ok := parseTwoPartInput(m.InputValue)
						if ok {
							if err := m.CurrentProject.AddSubTask(m.SelectedTaskID, section, subTask); err != nil {
								m.Error = err
								m.Status = "Sub task failed: create the sub section first with 's'."
								shouldCloseInput = false
								return m, nil
							}
							_ = storage.SaveRegistry(*m.CurrentProject)
							m.Status = "Sub task added."
						} else {
							m.Status = "Use format: section | subtask"
							shouldCloseInput = false
							return m, nil
						}
					}

				case InputToggleSubTask:
					if m.CurrentProject != nil && m.SelectedTaskID != 0 {
						section, subTask, ok := parseTwoPartInput(m.InputValue)
						if ok {
							_ = m.CurrentProject.ToggleSubTask(m.SelectedTaskID, section, subTask)
							_ = storage.SaveRegistry(*m.CurrentProject)
						}
					}
				}

				if shouldCloseInput {
					m.InputActive = false
					m.InputValue = ""
					m.InputType = InputNone
				}
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
				m.SelectedTaskID = 0

			case "n":
				m.InputActive = true
				m.InputType = InputNewProject
				m.InputValue = ""

			case "d":
				if len(m.Projects) == 0 {
					return m, nil
				}

				name := m.Projects[m.ProjectIndex]
				if err := storage.DeleteProject(name); err != nil {
					m.Error = err
					return m, nil
				}

				refreshProjects(&m)
				if m.Error != nil {
					return m, nil
				}
				if len(m.Projects) == 0 {
					m.ProjectIndex = 0
				} else if m.ProjectIndex >= len(m.Projects) {
					m.ProjectIndex = len(m.Projects) - 1
				}
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
				m.SelectedTaskID = 0

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

			case "enter":
				taskID, ok := selectedTaskID(m)
				if ok {
					m.SelectedTaskID = taskID
					m.ActivePane = PaneTaskDetail
				}

			case "m":
				taskID, ok := selectedTaskID(m)
				if ok {
					_ = m.CurrentProject.MoveTask(taskID, columnToNextStatus(m.ActiveColumn))
					_ = storage.SaveRegistry(*m.CurrentProject)
				}

			case "M":
				taskID, ok := selectedTaskID(m)
				if ok {
					_ = m.CurrentProject.MoveTask(taskID, columnToPrevStatus(m.ActiveColumn))
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

		case PaneTaskDetail:
			if m.CurrentProject == nil {
				m.ActivePane = PaneProjects
				return m, nil
			}
			if _, ok := selectedTask(m); !ok {
				m.ActivePane = PaneBoard
				return m, nil
			}

			switch msg.String() {
			case "esc":
				m.ActivePane = PaneBoard

			case "e":
				m.InputActive = true
				m.InputType = InputTaskDescription
				m.InputValue = ""

			case "s":
				m.InputActive = true
				m.InputType = InputNewSubSection
				m.InputValue = ""

			case "t":
				m.InputActive = true
				m.InputType = InputNewSubTask
				m.InputValue = ""

			case "x":
				m.InputActive = true
				m.InputType = InputToggleSubTask
				m.InputValue = ""
			}
		}
	}

	return m, nil
}

func refreshProjects(m *Model) {
	summaries, err := storage.ListProjectSummaries()
	if err != nil {
		m.Error = err
		return
	}

	projects := make([]string, 0, len(summaries))
	projectDates := make(map[string]time.Time, len(summaries))
	for _, summary := range summaries {
		projects = append(projects, summary.Name)
		projectDates[summary.Name] = summary.CreatedAt
	}

	m.Projects = projects
	m.ProjectDates = projectDates
}
