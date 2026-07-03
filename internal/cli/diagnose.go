package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/dat267/goclitmpl/pkg/diagnose"
	"github.com/spf13/cobra"
)

// NewDiagnoseCmd creates the "diagnose" parent command and registers all subcommands.
// No RunE is defined — Cobra shows help automatically when called with no subcommand.
func NewDiagnoseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "diagnose",
		Short: "System diagnosis and connectivity tools",
		Long:  `Perform local system checks, verify environment variables, and probe remote network connections.`,
		Example: `  goclitmpl diagnose info
  goclitmpl diagnose check github.com:443
  goclitmpl diagnose check --timeout 2s 1.1.1.1:53
  goclitmpl diagnose run
  goclitmpl diagnose run --timeout 10s`,
	}

	cmd.PersistentFlags().Duration("timeout", 5*time.Second, "max timeout duration for running diagnostics")

	cmd.AddCommand(
		newDiagnoseInfoCmd(),
		newDiagnoseCheckCmd(),
		newDiagnoseRunCmd(),
	)

	return cmd
}

// newDiagnoseInfoCmd calls pkg/diagnose.PrintInfo to display hardware and runtime specs.
func newDiagnoseInfoCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "info",
		Short:   "Print local system hardware and runtime specifications",
		Long:    `Gathers and displays local OS version, CPU count, Go environment parameters, and platform architecture.`,
		Example: `  goclitmpl diagnose info`,
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			diagnose.PrintInfo(cmd.OutOrStdout())
			return nil
		},
	}
}

// newDiagnoseCheckCmd calls pkg/diagnose.CheckAddress to probe TCP connectivity.
func newDiagnoseCheckCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "check [address]",
		Short: "Validate network connection to a remote target",
		Long: `Probes remote server connection availability using a TCP handshake.
Defaults to google.com:80 when no address is provided.`,
		Example: `  goclitmpl diagnose check
  goclitmpl diagnose check github.com:443
  goclitmpl diagnose check --timeout 2s 1.1.1.1:53`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			address := "google.com:80"
			if len(args) > 0 && strings.TrimSpace(args[0]) != "" {
				address = args[0]
			}

			timeout, err := cmd.Flags().GetDuration("timeout")
			if err != nil {
				return fmt.Errorf("failed to retrieve timeout: %w", err)
			}

			if err = diagnose.CheckAddress(address, timeout); err != nil {
				return fmt.Errorf("network connection to %s failed: %w", address, err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Successfully connected to %s!\n", address)
			return nil
		},
	}
}

// newDiagnoseRunCmd calls pkg/diagnose.RunSuite to run the full diagnostics suite.
func newDiagnoseRunCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "run",
		Short:   "Execute a full diagnostics checklist",
		Long:    `Runs all system specification diagnostics and queries standard network connection endpoints sequentially.`,
		Example: `  goclitmpl diagnose run
  goclitmpl diagnose run --timeout 10s`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			timeout, err := cmd.Flags().GetDuration("timeout")
			if err != nil {
				return fmt.Errorf("failed to retrieve timeout: %w", err)
			}

			diagnose.RunSuite(cmd.OutOrStdout(), timeout)
			return nil
		},
	}
}
