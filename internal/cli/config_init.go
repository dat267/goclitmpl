package cli

import (
	"fmt"
	"github.com/spf13/cobra"
)

func newConfigInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize a default configuration file",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Initializing config file...")
		},
	}
}
