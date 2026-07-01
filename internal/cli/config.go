package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dat267/goclitmpl/internal/config"
	"github.com/spf13/cobra"
)

// NewConfigCmd creates the base "config" command.
func NewConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage CLI configuration options",
		Long:  `Initialize or update configuration parameters for the application.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// If no subcommand is specified, show help
			return cmd.Help()
		},
	}

	// Add subcommands
	cmd.AddCommand(NewConfigInitCmd())

	return cmd
}

// NewConfigInitCmd creates the "config init" command.
func NewConfigInitCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize default configuration file",
		Long:  `Generates a standard default YAML configuration file in the user's home configuration directory.`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
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
			if err := os.MkdirAll(configDir, 0755); err != nil {
				return fmt.Errorf("failed to create config directory %s: %w", configDir, err)
			}

			// Write default configuration file
			if err := os.WriteFile(configPath, templateBytes, 0644); err != nil {
				return fmt.Errorf("failed to write configuration file %s: %w", configPath, err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Successfully initialized default configuration at: %s\n", configPath)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "overwrite existing configuration file if present")

	return cmd
}
