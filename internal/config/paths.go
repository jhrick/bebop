package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var (
	ErrGetWorkingDirectory     = errors.New("failed to get working directory")
	ErrResolveWorkingDirectory = errors.New("failed to resolve working directory")
	ErrInspectWorkingDirectory = errors.New("failed to inspect working directory")
)

const app = "bop"

func PreferencesPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, app), nil
}

func WorkingDirectory() (string, error) {
	workingDirectory, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrGetWorkingDirectory, err)
	}

	workingDirectory, err = filepath.EvalSymlinks(workingDirectory)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrResolveWorkingDirectory, err)
	}

	directoryInfo, err := os.Stat(workingDirectory)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrInspectWorkingDirectory, err)
	}

	if !directoryInfo.IsDir() {
		return "", fmt.Errorf("%w: %q", ErrInspectWorkingDirectory, workingDirectory)
	}

	return workingDirectory, nil
}
