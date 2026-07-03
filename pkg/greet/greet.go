// Package greet provides core business rules for greeting messages.
package greet

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
type Service struct {
	logger *slog.Logger
	env    string
}

// New instantiates a new greeting Service with injected dependencies.
func New(logger *slog.Logger, cfg Config) *Service {
	if logger == nil {
		logger = slog.Default()
	}
	return &Service{
		logger: logger,
		env:    cfg.Env,
	}
}

// Format generates a greeting string according to business rules.
func (s *Service) Format(ctx context.Context, name string, uppercase bool) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
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
