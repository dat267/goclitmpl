package cli

import "github.com/spf13/cobra"

func newDiagnoseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "diagnose",
		Short: "Troubleshoot and verify system health",
	}

	cmd.AddCommand(
		newDiagnoseRunCmd(),
		newMakeDiagnoseInfoCmd(),
		newDiagnoseCheckCmd(),
	)
	return cmd
}
