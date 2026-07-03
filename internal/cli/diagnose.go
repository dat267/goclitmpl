package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/dat267/goclitmpl/pkg/diagnose"
	"github.com/spf13/cobra"
)

// NewDiagnoseCmd creates the "diagnose" parent command and registers all subcommands.
func NewDiagnoseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "diagnose",
		Short: "System diagnosis and connectivity tools",
		Long:  `Perform local system checks, verify environment variables, and probe remote network connections.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
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
		Use:   "info",
		Short: "Print local system hardware and runtime specifications",
		Long:  `Gathers and displays local OS version, CPU count, Go environment parameters, and platform architecture.`,
		Args:  cobra.NoArgs,
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
		Long:  `Probes remote server connection availability using TCP handshakes. Default target is google.com:80.`,
		Args:  cobra.MaximumNArgs(1),
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
		Use:   "run",
		Short: "Execute a full diagnostics checklist",
		Long:  `Runs all system specification diagnostics and queries standard network connection endpoints sequentially.`,
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
}
