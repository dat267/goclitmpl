package cli

import (
	"github.com/dat267/goclitmpl/pkg/diagnose"
	"github.com/spf13/cobra"
)

func newMakeDiagnoseInfoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "info",
		Short: "Print system diagnostic information",
		Run: func(cmd *cobra.Command, args []string) {
			diagnose.PrintInfo(cmd.OutOrStdout())
		},
	}
}
