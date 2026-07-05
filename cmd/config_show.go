package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Print the active configuration values",
	RunE: func(cmd *cobra.Command, args []string) error {
		settings := viper.AllSettings()

		data, err := yaml.Marshal(settings)
		if err != nil {
			return fmt.Errorf("failed to marshal configuration: %w", err)
		}

		if len(data) == 0 || string(data) == "{}\n" {
			fmt.Println("# No configuration values are currently loaded.")
			return nil
		}

		fmt.Print(string(data))
		return nil
	},
}

func init() {
	configCmd.AddCommand(showCmd)
}
