package diagnose

import (
	"fmt"
	"log/slog"
	"net"
	"strings"

	"github.com/spf13/cobra"
)

// NewDiagnoseCheckCmd creates the leaf "diagnose check" command.
func NewDiagnoseCheckCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check [address]",
		Short: "Validate network connection to a remote target",
		Long:  `Probes remote server connection availability using TCP handshakes. Default target is google.com:80.`,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Resolve target address
			address := "google.com:80"
			if len(args) > 0 && strings.TrimSpace(args[0]) != "" {
				address = args[0]
			}

			// Retrieve persistent timeout flag
			timeout, err := cmd.Flags().GetDuration("timeout")
			if err != nil {
				return fmt.Errorf("failed to retrieve timeout: %w", err)
			}

			slog.Debug("checking network reachability", slog.String("address", address), slog.Duration("timeout", timeout))

			// Perform TCP connection probe
			conn, err := net.DialTimeout("tcp", address, timeout)
			if err != nil {
				return fmt.Errorf("network connection to %s failed: %w", address, err)
			}
			defer conn.Close()

			fmt.Fprintf(cmd.OutOrStdout(), "Successfully connected to %s!\n", address)
			return nil
		},
	}

	return cmd
}
