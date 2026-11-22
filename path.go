package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func getExecDir() (string, error) {
	execDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("get current directory: %w", err)
	}

	execDir, err = filepath.EvalSymlinks(execDir)
	if err != nil {
		return "", fmt.Errorf("resolve current directory: %w", err)
	}

	info, err := os.Stat(execDir)
	if err != nil {
		return "", fmt.Errorf("inspect current directory: %w", err)
	}

	if !info.IsDir() {
		return "", fmt.Errorf("current path is not a directory: %s", execDir)
	}

	return execDir, nil
}
