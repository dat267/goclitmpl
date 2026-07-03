package diagnose

import (
	"github.com/dat267/goclitmpl/pkg/diagnose"
	"github.com/spf13/cobra"
)

// NewDiagnoseInfoCmd creates the leaf "diagnose info" command.
func NewDiagnoseInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info",
		Short: "Print local system hardware and runtime specifications",
		Long:  `Gathers and displays local OS version, CPU count, Go environment parameters, and platform architecture.`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			diagnose.PrintInfo(cmd.OutOrStdout())
			return nil
		},
	}

	return cmd
}
