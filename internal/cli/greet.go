package cli

import (
	"fmt"
	"log/slog"

	"github.com/dat267/goclitmpl/pkg/greeting"
	"github.com/spf13/cobra"
)

// NewGreetCmd creates the greet command.
func NewGreetCmd() *cobra.Command {
	var uppercase bool

	cmd := &cobra.Command{
		Use:   "greet [name]",
		Short: "Greets a user",
		Long:  `Greets a user with customizable options while referencing loaded configuration.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			cfg := GetConfig(cmd.Context())

			slog.Debug("generating greeting",
				slog.String("name", name),
				slog.Bool("uppercase", uppercase),
			)

			// 1. Map global config to package-specific config struct
			pkgCfg := greeting.Config{
				Env: cfg.App.Env,
			}

			// 2. Instantiate business logic injecting explicit dependencies (slog.Default())
			svc := greeting.NewService(slog.Default(), pkgCfg)

			// 3. Propagate context and check errors
			message, err := svc.Format(cmd.Context(), name, uppercase)
			if err != nil {
				return fmt.Errorf("greeting execution failed: %w", err)
			}

			// 4. Output results to command output stream
			fmt.Fprintln(cmd.OutOrStdout(), message)

			return nil
		},
	}

	cmd.Flags().BoolVarP(&uppercase, "uppercase", "u", false, "convert greeting to UPPERCASE")

	return cmd
}
