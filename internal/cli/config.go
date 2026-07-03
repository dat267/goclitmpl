package cli

import "github.com/spf13/cobra"

func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage configuration profiles",
	}

	cmd.AddCommand(newConfigInitCmd())
	return cmd
}
