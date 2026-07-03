package cli

import (
	"runtime"

	"github.com/dat267/goclitmpl/pkg/version"
	"github.com/spf13/cobra"
)

// Injected variables at build time via ldflags
var (
	Version   = "dev"
	Commit    = "none"
	Date      = "unknown"
	GoVersion = runtime.Version()
)

// NewVersionCmd creates the base "version" command.
func NewVersionCmd() *cobra.Command {
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print application version information",
		Long:  `Displays compiler version, git commit hash, build timestamps, and runtime environment specifications.`,
		Example: `  goclitmpl version
  goclitmpl version --json`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return version.FormatVersion(cmd.OutOrStdout(), Version, Commit, Date, jsonOutput)
		},
	}

	cmd.Flags().BoolVar(&jsonOutput, "json", false, "display version metadata formatted as structured JSON")

	return cmd
}
