// Package diagnose provides system and connection health tools.
package diagnose

import (
	"time"

	"github.com/spf13/cobra"
)

// NewDiagnoseCmd creates the base "diagnose" command.
func NewDiagnoseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "diagnose",
		Short: "System diagnosis and connectivity tools",
		Long:  `Perform local system checks, verify environment variables, and probe remote network connections.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Show help if no subcommand is provided
			return cmd.Help()
		},
	}

	// Persistent flags inherited by all child commands (run, info, check)
	cmd.PersistentFlags().Duration("timeout", 5*time.Second, "max timeout duration for running diagnostics")

	// Add subcommands
	cmd.AddCommand(
		NewDiagnoseInfoCmd(),
		NewDiagnoseCheckCmd(),
		NewDiagnoseRunCmd(),
	)

	return cmd
}
