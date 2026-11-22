package config

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var (
	ErrCreateConfigurationDirectory = errors.New("failed to create configuration directory")
	ErrInstallConfigurationTemplate = errors.New("failed to install configuration template")
)

//go:embed templates/Dockerfile
var dockerfileTemplate string

//go:embed templates/compose.yml
var composeTemplate string

var configurationTemplates = map[string]string{
	"Dockerfile":  dockerfileTemplate,
	"compose.yml": composeTemplate,
}

func InstallTemplates(configurationDirectory string) error {
	if err := os.MkdirAll(configurationDirectory, 0o700); err != nil {
		return fmt.Errorf("%w %q: %w",
			ErrCreateConfigurationDirectory,
			configurationDirectory,
			err,
		)
	}

	for templateName, templateContent := range configurationTemplates {
		configurationTemplatePath := filepath.Join(configurationDirectory, templateName)

		if err := os.WriteFile(configurationTemplatePath, []byte(templateContent), 0o600); err != nil {
			if errors.Is(err, os.ErrExist) {
				continue
			}

			return fmt.Errorf("%w %q: %w",
				ErrInstallConfigurationTemplate,
				configurationTemplatePath,
				err,
			)
		}
	}

	return nil
}
