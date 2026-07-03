package cli

import (
	"fmt"

	"github.com/dat267/goclitmpl/internal/build"
	"github.com/dat267/goclitmpl/pkg/version"
	"github.com/spf13/cobra"
)

// NewVersionCmd creates the "version" command.
func NewVersionCmd() *cobra.Command {
	var format string

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print application version information",
		Long:  `Displays compiler version, git commit hash, build timestamps, and runtime environment specifications.`,
		Example: `  goclitmpl version
  goclitmpl version --format json`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			switch format {
			case "text", "":
				return version.FormatVersion(cmd.OutOrStdout(), build.Version, build.Commit, build.Date, false)
			case "json":
				return version.FormatVersion(cmd.OutOrStdout(), build.Version, build.Commit, build.Date, true)
			default:
				return fmt.Errorf("unknown format %q: must be 'text' or 'json'", format)
			}
		},
	}

	cmd.Flags().StringVar(&format, "format", "text", "output format: 'text' or 'json'")

	return cmd
}
