package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	Tasks     []Task    `json:"tasks"`
}

type Task struct {
	ID          int          `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	CreatedAt   time.Time    `json:"created_at"`
	Status      Status       `json:"status"`
	SubSections []SubSection `json:"sub_sections"`
}

type SubSection struct {
	Title    string    `json:"title"`
	SubTasks []SubTask `json:"sub_tasks"`
}

type SubTask struct {
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

type Status string

const (
	StatusTodo    Status = "TODO"
	StatusPending Status = "Pending"
	StatusDone    Status = "Done"
)

func NewProject(name string) Project {
	return Project{
		ID:        uuid.NewString(),
		Name:      name,
		CreatedAt: time.Now(),
		Tasks:     []Task{},
	}
}

func (p *Project) AddTask(title, description string) {
	task := Task{
		ID:          len(p.Tasks) + 1,
		Title:       title,
		Description: description,
		CreatedAt:   time.Now(),
		Status:      StatusTodo,
		SubSections: []SubSection{},
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

func (p *Project) SetTaskDescription(taskID int, description string) error {
	for i := range p.Tasks {
		if p.Tasks[i].ID == taskID {
			p.Tasks[i].Description = description
			return nil
		}
	}
	return errors.New("task not found")
}

func (p *Project) AddSubSection(taskID int, title string) error {
	for i := range p.Tasks {
		if p.Tasks[i].ID == taskID {
			p.Tasks[i].SubSections = append(p.Tasks[i].SubSections, SubSection{
				Title:    title,
				SubTasks: []SubTask{},
			})
			return nil
		}
	}
	return errors.New("task not found")
}

func (p *Project) AddSubTask(taskID int, sectionTitle, subTaskTitle string) error {
	for i := range p.Tasks {
		if p.Tasks[i].ID != taskID {
			continue
		}
		for j := range p.Tasks[i].SubSections {
			if p.Tasks[i].SubSections[j].Title == sectionTitle {
				p.Tasks[i].SubSections[j].SubTasks = append(
					p.Tasks[i].SubSections[j].SubTasks,
					SubTask{Title: subTaskTitle},
				)
				return nil
			}
		}
		return errors.New("sub section not found")
	}
	return errors.New("task not found")
}

func (p *Project) ToggleSubTask(taskID int, sectionTitle, subTaskTitle string) error {
	for i := range p.Tasks {
		if p.Tasks[i].ID != taskID {
			continue
		}
		for j := range p.Tasks[i].SubSections {
			if p.Tasks[i].SubSections[j].Title != sectionTitle {
				continue
			}
			for k := range p.Tasks[i].SubSections[j].SubTasks {
				if p.Tasks[i].SubSections[j].SubTasks[k].Title == subTaskTitle {
					p.Tasks[i].SubSections[j].SubTasks[k].Done = !p.Tasks[i].SubSections[j].SubTasks[k].Done
					return nil
				}
			}
			return errors.New("sub task not found")
		}
		return errors.New("sub section not found")
	}
	return errors.New("task not found")
}

func (p *Project) TaskByID(taskID int) (*Task, error) {
	for i := range p.Tasks {
		if p.Tasks[i].ID == taskID {
			return &p.Tasks[i], nil
		}
	}
	return nil, errors.New("task not found")
}
