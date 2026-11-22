package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/jhrick/bop/internal/config"
	"github.com/jhrick/bop/internal/service"
)

const (
	appName       = "bop"
	containerName = "harness"
)

var (
	ErrMissingProjectName = errors.New("missing project name")
	ErrInvalidArguments   = errors.New("expected exactly one project name")
	ErrInvalidProjectName = errors.New("invalid project name")
)

var projectNameRegex = regexp.MustCompile(
	`^[a-zA-Z0-9][a-zA-Z0-9_-]{0,39}$`,
)

func Execute() error {
	return run(os.Args[1:])
}

func run(args []string) error {
	if len(args) == 0 {
		printUsage()
		return ErrMissingProjectName
	}

	switch args[0] {
	case "-h", "--help", "help":
		printUsage()
		return nil
	case "-v", "--version", "version":
		fmt.Println(appName, "dev")
		return nil
	}

	if len(args) != 1 {
		printUsage()
		return ErrInvalidArguments
	}

	projectName := args[0]
	if !projectNameRegex.MatchString(projectName) {
		return fmt.Errorf(
			"%w %q: use 1-40 letters, numbers, hyphens or underscores; the first character must be alphanumeric",
			ErrInvalidProjectName,
			projectName,
		)
	}

	preferencesDirectory, err := config.PreferencesPath()
	if err != nil {
		return fmt.Errorf("resolve preferences directory: %w", err)
	}

	if err := config.InstallTemplates(preferencesDirectory); err != nil {
		return fmt.Errorf("install configuration templates: %w", err)
	}

	if err := service.CheckDocker(); err != nil {
		return fmt.Errorf("check Docker environment: %w", err)
	}

	workingDirectory, err := config.WorkingDirectory()
	if err != nil {
		return fmt.Errorf("resolve project directory: %w", err)
	}

	environment := service.IdentifyDockerEnvironment(
		projectName,
		workingDirectory,
	)

	environmentVariables := append(os.Environ(),
		"BOP_PROJECT_DIR="+workingDirectory,
		"BOP_CONFIG_DIR="+preferencesDirectory,
		fmt.Sprintf("BOP_UID=%d", os.Getuid()),
		fmt.Sprintf("BOP_GID=%d", os.Getgid()),
		"BOP_CONTAINER_NAME="+environment.ContainerName,
	)

	fmt.Printf("Project:   %s\n", workingDirectory)
	fmt.Printf("Container: %s\n", environment.ContainerName)
	fmt.Println("Starting Docker environment...")

	if err := service.StartEnvironment(
		preferencesDirectory,
		environment.ProjectName,
		environmentVariables,
	); err != nil {
		return fmt.Errorf("start Docker environment: %w", err)
	}

	fmt.Println("Opening OpenCode. Exit OpenCode to return to your shell.")

	command := exec.Command(
		"docker",
		"compose",
		"-f", filepath.Join(preferencesDirectory, "compose.yml"),
		"-p", environment.ProjectName,
		"exec", "-it", containerName, "opencode",
	)

	command.Env = environmentVariables
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		if exitError, ok := errors.AsType[*exec.ExitError](err); ok {
			return fmt.Errorf("OpenCode exited with status %d", exitError.ExitCode())
		}

		return fmt.Errorf("execute OpenCode: %w", err)
	}

	return nil
}

func printUsage() {
	fmt.Printf(`%s - isolated OpenCode environments using Docker Compose

Usage:
  bop <project-name>

Run this command from the root directory of the project you want OpenCode
to work on. Bop installs its Docker templates under ~/.config/bop on
first use, starts the environment, and opens OpenCode inside the container.

Options:
  -h, --help      Show this help
  -v, --version   Show version

`, appName)
}
