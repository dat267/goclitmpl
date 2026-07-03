package cli

import (
	"fmt"
	"time"

	"github.com/dat267/goclitmpl/pkg/diagnose"
	"github.com/spf13/cobra"
)

func newDiagnoseCheckCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "check <target>",
		Short: "Check network connectivity to a specific target host",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			target := args[0]
			fmt.Printf("Checking target endpoint %s...\n", target)
			if err := diagnose.CheckAddress(target, 5*time.Second); err != nil {
				return fmt.Errorf("reachability check failed: %w", err)
			}
			fmt.Println("Endpoint is reachable!")
			return nil
		},
	}
}
