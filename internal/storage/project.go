// safe

package storage

import (
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/goushalk/chromaboard/internal/domain"
)

type ProjectSummary struct {
	Name      string
	CreatedAt time.Time
}

func LoadRegistry(projectName string) (domain.Project, error) {
	var project domain.Project

	filePath, err := ProjectRegPath(projectName)
	if err != nil {
		return project, err
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return project, err
	}

	if err := json.Unmarshal(data, &project); err != nil {
		return project, err
	}

	return project, nil
}

func SaveRegistry(project domain.Project) error {
	if err := EnsureStorage(); err != nil {
		return err
	}

	filePath, err := ProjectRegPath(project.Name)
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(project, "", " ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(filePath, data, 0o664); err != nil {
		return err
	}
	return nil
}

func ListProjects() ([]string, error) {
	projectDir, err := ProjectDir()
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(projectDir)
	if err != nil {
		return nil, err
	}

	var projects []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if strings.HasSuffix(name, ".json") {
			projects = append(projects, strings.TrimSuffix(name, ".json"))
		}
	}

	return projects, nil
}

func DeleteProject(projectName string) error {
	filePath, err := ProjectRegPath(projectName)
	if err != nil {
		return err
	}

	if err := os.Remove(filePath); err != nil {
		return err
	}

	return nil
}

func ListProjectSummaries() ([]ProjectSummary, error) {
	projects, err := ListProjects()
	if err != nil {
		return nil, err
	}

	summaries := make([]ProjectSummary, 0, len(projects))
	for _, name := range projects {
		project, err := LoadRegistry(name)
		if err != nil {
			return nil, err
		}

		summaries = append(summaries, ProjectSummary{
			Name:      name,
			CreatedAt: project.CreatedAt,
		})
	}

	return summaries, nil
}
