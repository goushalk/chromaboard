package storage

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/goushalk/chromaboard/internal/domain"
)

// Load file.

func LoadRegistary(projectName string) (domain.Project, error) {

	var ProjectStruct domain.Project

	FilePath, err := ProjectRegPath(projectName)
	if err != nil {
		return ProjectStruct, err
	}

	data, err := os.ReadFile(FilePath)
	if err != nil {
		return ProjectStruct, err
	}

	jsonErr := json.Unmarshal(data, &ProjectStruct)
	if jsonErr != nil {
		return ProjectStruct, errors.New("curropted json")
	}

	return ProjectStruct, nil
}

// Create file if it dosent exist.

func CreateRegistry(projecFiletName string) (string, error) {

	EnsureStorage()
	projectName, err := ProjectRegPath(projecFiletName)

	if err != nil {
		return "", err
	}

	File, err := os.Create(projectName)

	if err != nil {
		return "file not created", err
	}
	defer File.Close()

	return "file successfuly created", nil
}
