package cli

import (
	"github.com/dat267/goclitmpl/internal/config"
	"github.com/dat267/goclitmpl/pkg/configinit"
	"github.com/spf13/cobra"
)

// newConfigCmd creates the config parent command.
func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage configuration profiles",
		Example: `  goclitmpl config init
  goclitmpl config init my-config.yaml`,
	}

	cmd.AddCommand(newConfigInitCmd())
	return cmd
}

// newConfigInitCmd creates the config init subcommand.
func newConfigInitCmd() *cobra.Command {
	var outputPath string

	cmd := &cobra.Command{
		Use:   "init [path]",
		Short: "Initialize a default configuration file",
		Example: `  goclitmpl config init
  goclitmpl config init config.yaml`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := outputPath
			if len(args) > 0 {
				path = args[0]
			}
			if path == "" {
				path = "config.yaml"
			}

			// Get the default configuration template embedded in internal/config
			cfgBytes := config.GetDefaultConfigTemplate()

			// Call business logic in pkg/configinit
			return configinit.Initialize(cmd.OutOrStdout(), path, cfgBytes)
		},
	}

	cmd.Flags().StringVarP(&outputPath, "output", "o", "", "output file path (default: config.yaml)")
	return cmd
}
