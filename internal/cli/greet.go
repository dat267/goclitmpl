package cli

import (
	"fmt"
	"log/slog"

	"github.com/dat267/goclitmpl/pkg/greet"
	"github.com/spf13/cobra"
)

// NewGreetCmd creates the greet command.
func NewGreetCmd() *cobra.Command {
	var uppercase bool

	cmd := &cobra.Command{
		Use:     "greet <name>",
		Short:   "Greets a user",
		Long:    `Greets a user with customizable options while referencing loaded configuration.`,
		Example: "  goclitmpl greet Alice\n  goclitmpl greet Alice --uppercase",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			cfg := GetConfig(cmd.Context())

			slog.Debug("generating greeting", slog.String("name", name), slog.Bool("uppercase", uppercase))

			pkgCfg := greet.Config{
				Env: cfg.App.Env,
			}

			svc := greet.New(slog.Default(), pkgCfg)
			message, err := svc.Format(cmd.Context(), name, uppercase)
			if err != nil {
				return fmt.Errorf("greeting execution failed: %w", err)
			}

			fmt.Fprintln(cmd.OutOrStdout(), message)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&uppercase, "uppercase", "u", false, "convert greeting to UPPERCASE")
	return cmd
}
