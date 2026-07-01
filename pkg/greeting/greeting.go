// Package greeting provides core business rules for greeting messages.
package greeting

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
)

// Config defines the configuration parameters specific to the greeting service.
type Config struct {
	Env string
}

// Service handles core greeting business rules.
// It is completely decoupled from any CLI parsing frameworks and external global states.
type Service struct {
	logger *slog.Logger
	env    string
}

// NewService instantiates a new greeting Service with injected dependencies (logger) and configuration.
func NewService(logger *slog.Logger, cfg Config) *Service {
	// Fallback to default logger if nil is passed to prevent panics
	if logger == nil {
		logger = slog.Default()
	}

	return &Service{
		logger: logger,
		env:    cfg.Env,
	}
}

// Format generates a greeting string according to business rules.
// It accepts a context for cancellation/timeout propagation and returns errors for validation checks.
func (s *Service) Format(ctx context.Context, name string, uppercase bool) (string, error) {
	// Check if context has been canceled before starting execution
	if err := ctx.Err(); err != nil {
		return "", err
	}

	// Business rule validation check
	if strings.TrimSpace(name) == "" {
		return "", errors.New("name cannot be empty or whitespace")
	}

	s.logger.Debug("generating greeting message", slog.String("name", name), slog.Bool("uppercase", uppercase))

	msg := fmt.Sprintf("Hello, %s! (env: %s)", name, s.env)
	if uppercase {
		msg = strings.ToUpper(msg)
	}

	return msg, nil
}
