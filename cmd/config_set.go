package cmd

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var setCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set a configuration value",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		rawValue := args[1]

		var schema map[string]interface{}
		if err := yaml.Unmarshal(defaultConfigFile, &schema); err != nil {
			return fmt.Errorf("failed to load baseline schema: %w", err)
		}

		flatDefaults := make(map[string]interface{})
		flattenSchema("", schema, flatDefaults)

		expectedDefault, knownKey := flatDefaults[key]
		if !knownKey {
			return fmt.Errorf("invalid configuration key %q: key does not exist in default schema template", key)
		}

		var finalValue interface{} = rawValue
		switch expectedDefault.(type) {
		case int:
			val, err := strconv.Atoi(rawValue)
			if err != nil {
				return fmt.Errorf("type validation error: key %q requires an integer value, got %q", key, rawValue)
			}
			finalValue = val
		case bool:
			val, err := strconv.ParseBool(rawValue)
			if err != nil {
				return fmt.Errorf("type validation error: key %q requires a boolean value (true/false), got %q", key, rawValue)
			}
			finalValue = val
		}

		viper.Set(key, finalValue)

		if err := viper.WriteConfig(); err != nil {
			var configFileNotFoundError viper.ConfigFileNotFoundError
			if errors.As(err, &configFileNotFoundError) {
				return errors.New("no active configuration file found to update; please run 'config init' first")
			}
			return fmt.Errorf("failed to save configuration: %w", err)
		}

		fmt.Printf("Updated key %q\n", key)
		return nil
	},
}

func init() {
	configCmd.AddCommand(setCmd)
}
