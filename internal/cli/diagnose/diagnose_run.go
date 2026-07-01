package diagnose

import (
	"fmt"
	"log/slog"
	"net"
	"runtime"
	"time"

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
			out := cmd.OutOrStdout()
			timeout, err := cmd.Flags().GetDuration("timeout")
			if err != nil {
				return fmt.Errorf("failed to retrieve timeout: %w", err)
			}

			slog.Info("starting full diagnostic suite run", slog.Duration("timeout", timeout))

			fmt.Fprintf(out, "Executing Diagnostics Suite...\n\n")

			// Check 1: CPU and hardware profile
			fmt.Fprintf(out, "[1/2] Hardware Specification Profile:\n")
			fmt.Fprintf(out, "  OS/Arch:      %s/%s\n", runtime.GOOS, runtime.GOARCH)
			fmt.Fprintf(out, "  CPU Cores:    %d\n", runtime.NumCPU())
			fmt.Fprintf(out, "  Go Version:   %s\n\n", runtime.Version())

			// Check 2: Core endpoints connection
			targets := []string{"google.com:80", "github.com:80"}
			fmt.Fprintf(out, "[2/2] Remote Endpoints Network Probe:\n")

			for _, target := range targets {
				start := time.Now()
				conn, err := net.DialTimeout("tcp", target, timeout)
				duration := time.Since(start)

				if err != nil {
					slog.Warn("endpoint probe failed", slog.String("target", target), slog.Any("error", err))
					fmt.Fprintf(out, "  %-15s -> FAILED (%v)\n", target, err)
				} else {
					conn.Close()
					slog.Debug("endpoint probe succeeded", slog.String("target", target), slog.Duration("latency", duration))
					fmt.Fprintf(out, "  %-15s -> SUCCESS (latency: %v)\n", target, duration.Round(time.Millisecond))
				}
			}

			fmt.Fprintf(out, "\nDiagnostics run completed successfully!\n")
			return nil
		},
	}

	return cmd
}
