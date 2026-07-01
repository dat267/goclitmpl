// Package config loads and validates application parameters.
package config

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

//go:embed default_config.yaml
var defaultConfigTemplate []byte

// GetDefaultConfigTemplate returns the embedded default configuration file.
func GetDefaultConfigTemplate() []byte {
	return defaultConfigTemplate
}

const (
	// EnvPrefix is the prefix used for environment variables.
	EnvPrefix = "GOCLITMPL"
	// DefaultConfigFileName is the name of the config file without extension.
	DefaultConfigFileName = "config"
	// AppName is the application name used to structure config folder paths.
	AppName = "goclitmpl"
)

// Config represents the master configuration struct for the application.
type Config struct {
	Log LogConfig `mapstructure:"log"`
	App AppConfig `mapstructure:"app"`
}

// AppConfig represents custom application configuration parameters.
type AppConfig struct {
	Env   string `mapstructure:"env"`
	Debug bool   `mapstructure:"debug"`
	Port  int    `mapstructure:"port"`
}

// NewDefaultConfig returns a configuration initialized with defaults.
func NewDefaultConfig() *Config {
	return &Config{
		Log: DefaultLogConfig(),
		App: AppConfig{
			Env:   "production",
			Debug: false,
			Port:  8080,
		},
	}
}

// Load loads the configuration from a file, environment variables, or CLI flag overrides.
// configFile parameter points to a specific configuration file to load (optional).
func Load(configFile string) (*Config, error) {
	v := viper.New()

	// 1. Set Defaults
	v.SetDefault("log.level", "info")
	v.SetDefault("log.format", "text")
	v.SetDefault("app.env", "production")
	v.SetDefault("app.debug", false)
	v.SetDefault("app.port", 8080)

	// 2. Setup Environment Variables override
	v.SetEnvPrefix(EnvPrefix)
	// Bind env variables with nesting, e.g. GOCLITMPL_LOG_LEVEL maps to log.level
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// 3. Load from Config File
	if configFile != "" {
		// Use config file from the flag
		v.SetConfigFile(configFile)
	} else {
		// Search config in standard places
		v.SetConfigName(DefaultConfigFileName)
		v.SetConfigType("yaml") // Search yaml by default

		// Path 1: Current working directory
		v.AddConfigPath(".")

		// Path 2: Platform-specific user config directory (e.g., ~/.config/goclitmpl/config.yaml)
		if home, err := os.UserConfigDir(); err == nil {
			v.AddConfigPath(filepath.Join(home, AppName))
		}

		// Path 3: System configuration (/etc/goclitmpl/config.yaml)
		//nolint:gocritic // /etc is a standard system directory path
		v.AddConfigPath(filepath.Join("/etc", AppName))
	}

	// Attempt to read the config file
	err := v.ReadInConfig()
	if err != nil {
		// If config file is not found, we can proceed if we didn't explicitly request one.
		// If explicitly requested, we return the error.
		if configFile != "" {
			return nil, fmt.Errorf("error reading config file %s: %w", configFile, err)
		}
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error parsing config file: %w", err)
		}
	}

	// Unmarshal configuration into struct
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode configuration: %w", err)
	}

	// 4. Validate configuration parameters
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return &cfg, nil
}

// Validate performs assertions on the config and returns errors for invalid setups.
func (c *Config) Validate() error {
	if c.App.Port < 1 || c.App.Port > 65535 {
		return fmt.Errorf("app.port must be between 1 and 65535, got %d", c.App.Port)
	}

	logLevel := strings.ToLower(c.Log.Level)
	switch logLevel {
	case "debug", "info", "warn", "error":
	default:
		return fmt.Errorf("log.level must be debug, info, warn, or error, got %q", c.Log.Level)
	}

	logFormat := strings.ToLower(c.Log.Format)
	switch logFormat {
	case "text", "json":
	default:
		return fmt.Errorf("log.format must be text or json, got %q", c.Log.Format)
	}

	return nil
}
