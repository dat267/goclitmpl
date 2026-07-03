// Package config provides running logic for CLI configuration management commands.
package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	appconfig "github.com/dat267/goclitmpl/internal/config"
)

// Initialize writes the default configuration template to the user's config directory.
func Initialize(w io.Writer, force bool) error {
	home, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("failed to retrieve config directory: %w", err)
	}
	configDir := filepath.Join(home, appconfig.AppName)
	configPath := filepath.Join(configDir, "config.yaml")

	// Protect against accidental overwrite
	if _, statErr := os.Stat(configPath); statErr == nil && !force {
		return fmt.Errorf("config file already exists at %s (use --force or -f to overwrite)", configPath)
	}

	templateBytes := appconfig.GetDefaultConfigTemplate()

	if err = os.MkdirAll(configDir, 0o700); err != nil {
		return fmt.Errorf("failed to create config directory %s: %w", configDir, err)
	}

	if err = os.WriteFile(configPath, templateBytes, 0o600); err != nil {
		return fmt.Errorf("failed to write configuration file %s: %w", configPath, err)
	}

	fmt.Fprintf(w, "Successfully initialized default configuration at: %s\n", configPath)
	return nil
}
