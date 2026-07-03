package cli

import (
	"context"
	"io"
	"log/slog"
	"os"

	"github.com/dat267/goclitmpl/internal/config"
	"github.com/spf13/cobra"
)

type contextKey string

const (
	configKey contextKey = "config"
)

// NewRootCmd creates the root Cobra command.
// We accept stdout and stderr writers to make the commands easily testable.
func NewRootCmd(out, errOut io.Writer) *cobra.Command {
	var configFile string
	var verbose bool
	var logFormat string

	cmd := &cobra.Command{
		Use:   "goclitmpl",
		Short: "goclitmpl is an optimized Go CLI template project",
		Long: `An optimized, structured, and extensible boilerplate template for 
building production-ready CLI applications in Go.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// 1. Load configuration
			cfg, err := config.Load(configFile)
			if err != nil {
				return err
			}

			// 2. Apply persistent flags override to config
			if verbose {
				cfg.Log.Level = "debug"
			}
			if logFormat != "" {
				cfg.Log.Format = logFormat
			}

			// 3. Initialize structured logging
			config.SetupLogging(cfg.Log, errOut)
			cmd.SilenceUsage = true // Don't print usage instructions on execution errors

			// 4. Store configuration in the command context
			ctx := context.WithValue(cmd.Context(), configKey, cfg)
			cmd.SetContext(ctx)

			slog.Debug("logger initialized and configuration loaded",
				slog.String("level", cfg.Log.Level),
				slog.String("format", cfg.Log.Format),
				slog.Int("port", cfg.App.Port),
				slog.String("env", cfg.App.Env),
			)

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// By default, if no subcommand is executed, show help
			return cmd.Help()
		},
	}

	// Set standard output streams (very useful for unit testing commands)
	cmd.SetOut(out)
	cmd.SetErr(errOut)

	// Define persistent flags
	cmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file path (default search paths: ., ~/.goclitmpl, /etc/goclitmpl)")
	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose (debug) logging output")
	cmd.PersistentFlags().StringVar(&logFormat, "log-format", "", "log format: 'text' or 'json'")

	// Add subcommands explicitly for fine-grained control
	cmd.AddCommand(
		NewVersionCmd(),
		NewGreetCmd(),
		NewConfigCmd(),
		NewDiagnoseCmd(),
	)

	return cmd
}

// GetConfig extracts the application configuration from the context.
func GetConfig(ctx context.Context) *config.Config {
	if cfg, ok := ctx.Value(configKey).(*config.Config); ok {
		return cfg
	}
	// Fallback to defaults if not found in context (e.g. in tests that don't run PreRunE)
	return config.NewDefaultConfig()
}

// Execute runs the root command with default OS stdout, stderr, and arguments.
func Execute(ctx context.Context) int {
	cmd := NewRootCmd(os.Stdout, os.Stderr)
	if err := cmd.ExecuteContext(ctx); err != nil {
		// Output error if silencing usage prevented it, but since Cobra handles it,
		// we just return a non-zero exit code.
		return 1
	}
	return 0
}
