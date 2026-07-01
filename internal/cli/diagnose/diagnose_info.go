package diagnose

import (
	"fmt"
	"runtime"

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
			out := cmd.OutOrStdout()

			fmt.Fprintf(out, "System Runtime Specifications:\n")
			fmt.Fprintf(out, "  Operating System: %s\n", runtime.GOOS)
			fmt.Fprintf(out, "  Architecture:     %s\n", runtime.GOARCH)
			fmt.Fprintf(out, "  Go Version:       %s\n", runtime.Version())
			fmt.Fprintf(out, "  CPU Count:        %d\n", runtime.NumCPU())

			return nil
		},
	}

	return cmd
}
