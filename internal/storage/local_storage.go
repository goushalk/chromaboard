package storage

import (
	"errors"
	"os"
	"path/filepath"
)

const (
	AppName = "chromaboard"
)

func DataDir() (string, error) { // defines the directory where the project DIR should be.

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, ".local", "share", AppName), nil
}

func ProjectDir() (string, error) { // defines where the projects files should be.

	dataDir, err := DataDir()

	if err != nil {
		return "", err
	}

	return filepath.Join(dataDir, "projects"), nil
}

func ProjectRegPath(projectName string) (string, error) { // defines where the projects registry file should be.

	dataDir, err := ProjectDir()

	if err != nil {
		return "", err
	}
	fileName := projectName + ".json"
	return filepath.Join(dataDir, fileName), nil
}

func EnsureStorage() error { // Ensures that all the above directories exist.

	dataDir, err := DataDir()

	if err != nil {
		return err
	}

	projectDir, err := ProjectDir()
	if err != nil {

		return err
	}

	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return err
	}

	if err := os.MkdirAll(projectDir, 0o755); err != nil {
		return err
	}

	return nil
}

// checks for the file.

func ExistsPath(path string) (bool, error) {
	if path == "" {
		return false, errors.New("Empty path")
	}

	_, err := os.Stat(path)

	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err

}
