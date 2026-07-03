// Package configinit provides running logic to initialize the default configuration file.
package configinit

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/dat267/goclitmpl/internal/config"
)

// Initialize writes the default configuration template to the user's config directory.
func Initialize(w io.Writer, force bool) error {
	// Find user's config directory (e.g. ~/.config/goclitmpl/)
	home, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("failed to retrieve config directory: %w", err)
	}
	configDir := filepath.Join(home, config.AppName)
	configPath := filepath.Join(configDir, "config.yaml")

	// Protect against accidental overwrite of existing settings
	if _, err := os.Stat(configPath); err == nil && !force {
		return fmt.Errorf("config file already exists at %s (use --force or -f to overwrite)", configPath)
	}

	// Get embedded configuration bytes
	templateBytes := config.GetDefaultConfigTemplate()

	// Ensure config directory exists
	if err := os.MkdirAll(configDir, 0o700); err != nil {
		return fmt.Errorf("failed to create config directory %s: %w", configDir, err)
	}

	// Write default configuration file
	if err := os.WriteFile(configPath, templateBytes, 0o600); err != nil {
		return fmt.Errorf("failed to write configuration file %s: %w", configPath, err)
	}

	fmt.Fprintf(w, "Successfully initialized default configuration at: %s\n", configPath)
	return nil
}
