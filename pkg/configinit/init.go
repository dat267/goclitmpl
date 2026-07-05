// Package configinit provides business logic for configuration file initialization.
package configinit

import (
	"fmt"
	"io"
	"os"
)

// Initialize writes the default configuration template to the specified path.
// It conforms to the security requirements by using 0o600 permissions for owner-only access.
func Initialize(w io.Writer, path string, template []byte) error {
	fmt.Fprintf(w, "Initializing config file at %s...\n", path)

	// Check if file already exists to prevent accidental overwrites
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("config file already exists at %s", path)
	}

	// Write configuration template using owner-only permissions (0o600)
	if err := os.WriteFile(path, template, 0o600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	fmt.Fprintln(w, "Configuration file initialized successfully.")
	return nil
}
