package tui

import "github.com/goushalk/chromaboard/internal/domain"

/*
Returns the ID of the currently selected task in the active column.
*/
func selectedTaskID(m Model) (int, bool) {
	if m.CurrentProject == nil {
		return 0, false
	}

	index := 0
	for _, t := range m.CurrentProject.Tasks {
		if statusToColumn(t.Status) != m.ActiveColumn {
			continue
		}

		if index == m.TaskIndex {
			return t.ID, true
		}
		index++
	}

	return 0, false
}

/*
Counts how many tasks exist in the currently active column.
*/
func countTasksInColumn(m Model) int {
	if m.CurrentProject == nil {
		return 0
	}

	count := 0
	for _, t := range m.CurrentProject.Tasks {
		if statusToColumn(t.Status) == m.ActiveColumn {
			count++
		}
	}

	return count
}

/*
Converts a board column to the next logical task status.
Used when pressing `m` (move task).
*/
func columnToNextStatus(col Column) domain.Status {
	switch col {
	case ColumnTodo:
		return domain.StatusPending
	case ColumnPending:
		return domain.StatusDone
	case ColumnDone:
		return domain.StatusDone
	default:
		return domain.StatusTodo
	}
}

/*
Maps a domain task status to a UI column.
*/
func statusToColumn(status domain.Status) Column {
	switch status {
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
