package main

import (
	"crypto/sha256"
	_ "embed"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

//go:embed templates/Dockerfile
var DockerfileTmpl string

//go:embed templates/compose.yml
var ComposeTmpl string

const (
	app        = "bop"
	workingDir = "/workspace"
	service    = "harness"
)

var regexProjectName = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]{0,39}$`)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "bop: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		printUsage()
		return errors.New("missing project name")
	}

	switch args[0] {
	case "-h", "--help", "help":
		printUsage()
		return nil
	case "-v", "--version", "version":
		fmt.Println("bop dev")
		return nil
	}

	if len(args) != 1 {
		printUsage()
		return errors.New("expected exactly one project name")
	}

	projectName := args[0]
	if !regexProjectName.MatchString(projectName) {
		return fmt.Errorf(
			"invalid project name %q: use 1-40 letters, numbers, hyphens or underscores; the first character must be alphanumeric",
			projectName,
		)
	}

	configDir, err := configDirectory()
	if err != nil {
		return err
	}

	if err := installTemplates(configDir); err != nil {
		return err
	}

	if err := checkDocker(); err != nil {
		return err
	}

	execDir, err := getExecDir()
	if err != nil {
		return err
	}

	composeProject, containerName := identifyEnvironment(projectName, execDir)
	env := append(os.Environ(),
		"bop_PROJECT_DIR="+execDir,
		"bop_CONFIG_DIR="+configDir,
		fmt.Sprintf("bop_UID=%d", os.Getuid()),
		fmt.Sprintf("bop_GID=%d", os.Getgid()),
		"bop_CONTAINER_NAME="+containerName,
	)

	fmt.Printf("Project:   %s\n", execDir)
	fmt.Printf("Container: %s\n", containerName)
	fmt.Println("Starting Docker environment...")

	if err := compose(configDir, composeProject, env, "up", "-d", "--build"); err != nil {
		return fmt.Errorf("start environment: %w", err)
	}

	fmt.Println("Opening OpenCode. Exit OpenCode to return to your shell.")

	// Pass terminal signals and standard streams through to Docker Compose.
	cmd := exec.Command("docker", "compose",
		"-f", filepath.Join(configDir, "compose.yml"),
		"-p", composeProject,
		"exec", "-it", service, "opencode",
	)
	cmd.Env = env
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if exitErr, ok := errors.AsType[*exec.ExitError](err); ok {
			// Preserve OpenCode's exit status without hiding the Docker error.
			return fmt.Errorf("OpenCode exited with status %d", exitErr.ExitCode())
		}
		return fmt.Errorf("execute OpenCode: %w", err)
	}

	return nil
}

func configDirectory() (string, error) {
	root := os.Getenv("XDG_CONFIG_HOME")
	if root == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("find home directory: %w", err)
		}

		root = filepath.Join(home, ".config")
	}

	return filepath.Join(root, app), nil
}

func identifyEnvironment(name, projectDir string) (string, string) {
	sum := sha256.Sum256([]byte(projectDir))
	suffix := fmt.Sprintf("%x", sum[:3])
	slug := strings.ToLower(name)
	slug = strings.Trim(slug, "-_")
	if slug == "" {
		slug = "project"
	}

	// Compose project names are lowercase and use [a-z0-9_-].
	project := "bop-" + slug + "-" + suffix
	container := project + "-harness"
	return project, container
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

`, app)
}
