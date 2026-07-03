package cli

import (
	pkgconfig "github.com/dat267/goclitmpl/pkg/config"
	"github.com/spf13/cobra"
)

// NewConfigCmd creates the "config" parent command.
// No RunE is defined — Cobra shows help automatically when called with no subcommand.
func NewConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage CLI configuration options",
		Long:  `Initialize or update configuration parameters for the application.`,
		Example: `  goclitmpl config init
  goclitmpl config init --force`,
	}

	cmd.AddCommand(newConfigInitCmd())

	return cmd
}

// newConfigInitCmd calls pkg/config.Initialize to write the default config file.
func newConfigInitCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:     "init",
		Short:   "Initialize default configuration file",
		Long:    `Generates a standard default YAML configuration file in the user's home configuration directory.`,
		Example: `  goclitmpl config init
  goclitmpl config init --force`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return pkgconfig.Initialize(cmd.OutOrStdout(), force)
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "overwrite existing configuration file if present")

	return cmd
}
