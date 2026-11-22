package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func installTemplates(configDir string) error {
	if err := os.MkdirAll(configDir, 0o700); err != nil {
		return fmt.Errorf("create configuration directory %s: %w", configDir, err)
	}

	files := map[string]string{
		"Dockerfile":  DockerfileTmpl,
		"compose.yml": ComposeTmpl,
	}

	for name, content := range files {
		path := filepath.Join(configDir, name)
		if _, err := os.Stat(path); err == nil {
			continue // Preserve user modifications to the installed templates.
		} else if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("inspect %s: %w", path, err)
		}

		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			return fmt.Errorf("install %s: %w", path, err)
		}
	}

	return nil
}

func checkDocker() error {
	if _, err := exec.LookPath("docker"); err != nil {
		return errors.New("Docker CLI not found in PATH")
	}

	cmd := exec.Command("docker", "compose", "version")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return errors.New("Docker Compose is unavailable; install or enable the Docker Compose plugin")
	}

	return nil
}
func compose(configDir, project string, env []string, args ...string) error {
	commandArgs := []string{
		"compose",
		"-f", filepath.Join(configDir, "compose.yml"),
		"-p", project,
	}
	commandArgs = append(commandArgs, args...)

	cmd := exec.Command("docker", commandArgs...)
	cmd.Env = env
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
