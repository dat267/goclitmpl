package diagnose

import (
	"fmt"

	"github.com/dat267/goclitmpl/pkg/diagnose"
	"github.com/spf13/cobra"
)

// NewDiagnoseRunCmd creates the leaf "diagnose run" command.
func NewDiagnoseRunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Execute a full diagnostics checklist",
		Long:  `Runs all system specifications diagnostics and queries standard network connection endpoints sequentially.`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			timeout, err := cmd.Flags().GetDuration("timeout")
			if err != nil {
				return fmt.Errorf("failed to retrieve timeout: %w", err)
			}

			diagnose.RunSuite(cmd.OutOrStdout(), timeout)
			return nil
		},
	}

	return cmd
}
