package domain

import (
	"errors"

	"github.com/google/uuid"
)

type Project struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Tasks []Task `json:"tasks"`
}

type Task struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status Status `json:"status"`
}

type Status string

const (
	StatusTodo    Status = "TODO"
	StatusPending Status = "Pending"
	StatusDone    Status = "Done"
)

func NewProject(name string) Project {
	return Project{
		ID:    uuid.NewString(),
		Name:  name,
		Tasks: []Task{},
	}
}

func (p *Project) AddTask(title string) {
	task := Task{
		ID:     len(p.Tasks) + 1,
		Title:  title,
		Status: StatusTodo,
	}
	p.Tasks = append(p.Tasks, task)
}

func (p *Project) MoveTask(taskID int, status Status) error {
	for i := range p.Tasks {
		if p.Tasks[i].ID == taskID {
			p.Tasks[i].Status = status
			return nil
		}
	}
	return errors.New("task not found")
}

func (p *Project) DeleteTask(taskID int) error {
	for i := range p.Tasks {
		if p.Tasks[i].ID == taskID {
			p.Tasks = append(p.Tasks[:i], p.Tasks[i+1:]...)
			return nil
		}
	}
	return errors.New("task not found")
}

func (p *Project) RenameTask(taskID int, title string) error {
	for i := range p.Tasks {
		if p.Tasks[i].ID == taskID {
			p.Tasks[i].Title = title
			return nil
		}
	}
	return errors.New("task not found")
}
