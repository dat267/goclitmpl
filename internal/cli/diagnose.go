package cli

import (
	"fmt"
	"time"

	"github.com/dat267/goclitmpl/pkg/diagnose"
	"github.com/spf13/cobra"
)

// newDiagnoseCmd creates the diagnose parent command.
func newDiagnoseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "diagnose",
		Short: "Troubleshoot and verify system health",
		Example: `  goclitmpl diagnose info
  goclitmpl diagnose check google.com:443
  goclitmpl diagnose check google.com:443 --timeout 10s`,
	}

	// Persistent flags shared by subcommands
	cmd.PersistentFlags().Duration("timeout", 5*time.Second, "timeout for connection checks")

	cmd.AddCommand(
		newDiagnoseInfoCmd(),
		newDiagnoseCheckCmd(),
	)
	return cmd
}

// newDiagnoseInfoCmd creates the diagnose info subcommand.
func newDiagnoseInfoCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "info",
		Short:   "Print system diagnostic information",
		Example: "  goclitmpl diagnose info",
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			diagnose.PrintInfo(cmd.OutOrStdout())
		},
	}
}

// newDiagnoseCheckCmd creates the diagnose check subcommand.
func newDiagnoseCheckCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "check <target>",
		Short: "Check network connectivity to a specific target host",
		Example: `  goclitmpl diagnose check google.com:443
  goclitmpl diagnose check google.com:443 --timeout 10s`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			target := args[0]
			timeout, err := cmd.Flags().GetDuration("timeout")
			if err != nil {
				timeout = 5 * time.Second
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Checking target endpoint %s...\n", target)
			if err := diagnose.CheckAddress(target, timeout); err != nil {
				return fmt.Errorf("reachability check failed: %w", err)
			}
			fmt.Fprintln(cmd.OutOrStdout(), "Endpoint is reachable!")
			return nil
		},
	}
}
