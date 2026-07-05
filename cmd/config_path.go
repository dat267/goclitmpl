package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var pathCmd = &cobra.Command{
	Use:   "path",
	Short: "Show the active configuration file path",
	Run: func(cmd *cobra.Command, args []string) {
		if activeFile := viper.ConfigFileUsed(); activeFile != "" {
			fmt.Println(activeFile)
			return
		}

		if cfgFile != "" {
			fmt.Println(cfgFile)
			return
		}

		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error detecting home directory:", err)
			os.Exit(1)
		}
		fmt.Println(filepath.Join(home, ".goclitmpl.yaml"))
	},
}

func init() {
	configCmd.AddCommand(pathCmd)
}
