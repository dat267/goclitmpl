// Package cli implements Cobra command routers.
package cli

import (
	"github.com/dat267/goclitmpl/pkg/configinit"
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
			return configinit.Initialize(cmd.OutOrStdout(), force)
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "overwrite existing configuration file if present")

	return cmd
}
