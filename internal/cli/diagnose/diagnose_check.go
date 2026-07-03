package diagnose

import (
	"fmt"
	"strings"

	"github.com/dat267/goclitmpl/pkg/diagnose"
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

			// Perform TCP connection probe
			err = diagnose.CheckAddress(address, timeout)
			if err != nil {
				return fmt.Errorf("network connection to %s failed: %w", address, err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Successfully connected to %s!\n", address)
			return nil
		},
	}

	return cmd
}
