package cmd

import (
	"bufio"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

//go:embed default_config.yaml
var defaultConfigFile []byte

var forceOverwrite bool

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration with interactive prompts",
	RunE: func(cmd *cobra.Command, args []string) error {
		targetPath := viper.ConfigFileUsed()
		if targetPath == "" {
			if cfgFile != "" {
				targetPath = cfgFile
			} else {
				home, err := os.UserHomeDir()
				if err != nil {
					return fmt.Errorf("failed to detect home directory: %w", err)
				}
				targetPath = filepath.Join(home, fmt.Sprintf(".%s.yaml", ConfigName))
			}
		}

		if !forceOverwrite {
			if _, err := os.Stat(targetPath); err == nil {
				fmt.Println("Configuration file already exists. Use --force or -f to overwrite.")
				return nil
			}
		}

		var schema map[string]interface{}
		if err := yaml.Unmarshal(defaultConfigFile, &schema); err != nil {
			return fmt.Errorf("invalid embedded default config: %w", err)
		}

		reader := bufio.NewReader(os.Stdin)
		fmt.Println("--- Interactive CLI Configuration Setup ---")

		flatDefaults := make(map[string]interface{})
		flattenSchema("", schema, flatDefaults)

		for key, defaultValue := range flatDefaults {
			promptLabel := fmt.Sprintf("Set value for %q", key)
			userInput := prompt(reader, promptLabel, fmt.Sprintf("%v", defaultValue))

			var validatedValue interface{} = userInput
			switch defaultValue.(type) {
			case int:
				val, err := strconv.Atoi(userInput)
				if err != nil {
					return fmt.Errorf("validation error: key %q requires an integer, got %q", key, userInput)
				}
				validatedValue = val
			case bool:
				val, err := strconv.ParseBool(userInput)
				if err != nil {
					return fmt.Errorf("validation error: key %q requires a boolean value (true/false), got %q", key, userInput)
				}
				validatedValue = val
			}

			viper.Set(key, validatedValue)
		}

		if err := viper.WriteConfigAs(targetPath); err != nil {
			return fmt.Errorf("failed to write config file: %w", err)
		}

		fmt.Printf("\nSuccessfully initialized configuration in %s\n", targetPath)
		return nil
	},
}

func flattenSchema(prefix string, current map[string]interface{}, result map[string]interface{}) {
	for k, v := range current {
		fullKey := k
		if prefix != "" {
			fullKey = prefix + "." + k
		}

		if nestedMap, ok := v.(map[string]interface{}); ok {
			flattenSchema(fullKey, nestedMap, result)
		} else {
			result[fullKey] = v
		}
	}
}

func prompt(reader *bufio.Reader, label, defaultValue string) string {
	fmt.Printf("%s [%s]: ", label, defaultValue)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		return defaultValue
	}
	return input
}

func init() {
	initCmd.Flags().BoolVarP(&forceOverwrite, "force", "f", false, "Force overwrite existing configuration file")
	configCmd.AddCommand(initCmd)
}
