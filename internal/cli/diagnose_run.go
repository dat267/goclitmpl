package cli

import (
	"time"

	"github.com/dat267/goclitmpl/pkg/diagnose"
	"github.com/spf13/cobra"
)

func newDiagnoseRunCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "run",
		Short: "Execute all diagnostic diagnostics",
		Run: func(cmd *cobra.Command, args []string) {
			diagnose.RunSuite(cmd.OutOrStdout(), 5*time.Second)
		},
	}
}
