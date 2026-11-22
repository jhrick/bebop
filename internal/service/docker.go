package service

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	ErrDockerNotFound     = errors.New("Docker CLI not found in PATH")
	ErrComposeUnavailable = errors.New("Docker Compose plugin is unavailable")
	ErrDaemonUnavailable  = errors.New("Docker daemon is unavailable")
	ErrCreateConfigDir    = errors.New("failed to create configuration directory")
	ErrInspectTemplate    = errors.New("failed to inspect template")
	ErrInstallTemplate    = errors.New("failed to install template")
)

type DockerEnvironment struct {
	ProjectName   string
	ContainerName string
}

func CheckDocker() error {
	if _, err := exec.LookPath("docker"); err != nil {
		return fmt.Errorf("%w: %w", ErrDockerNotFound, err)
	}

	if err := exec.Command("docker", "compose", "version").Run(); err != nil {
		return fmt.Errorf("%w: %w", ErrComposeUnavailable, err)
	}

	if err := exec.Command("docker", "info").Run(); err != nil {
		return fmt.Errorf("%w: %w", ErrDaemonUnavailable, err)
	}

	return nil
}

func IdentifyDockerEnvironment(name, projectDirectory string) DockerEnvironment {
	projectHash := sha256.Sum256([]byte(projectDirectory))
	projectSuffix := fmt.Sprintf("%x", projectHash[:3])

	projectSlug := strings.ToLower(name)
	projectSlug = strings.Trim(projectSlug, "-_")

	if projectSlug == "" {
		projectSlug = "project"
	}

	projectName := "bop-" + projectSlug + "-" + projectSuffix

	return DockerEnvironment{
		ProjectName:   projectName,
		ContainerName: projectName + "-harness",
	}
}

func StartEnvironment(
	configurationDirectory string,
	projectName string,
	environmentVariables []string,
) error {
	composeFilePath := filepath.Join(configurationDirectory, "compose.yml")

	command := exec.Command(
		"docker",
		"compose",
		"-f", composeFilePath,
		"-p", projectName,
		"up", "-d", "--build",
	)

	command.Env = environmentVariables
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		return fmt.Errorf("execute Docker Compose: %w", err)
	}

	return nil
}
