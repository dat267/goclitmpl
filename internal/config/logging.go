package config

import (
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

// LogFormat represents the log output format.
type LogFormat string

const (
	// LogFormatText outputs logs in a human-readable text format.
	LogFormatText LogFormat = "text"
	// LogFormatJSON outputs logs in structured JSON format.
	LogFormatJSON LogFormat = "json"
)

// LogConfig holds configuration for the logger.
type LogConfig struct {
	Level    string `mapstructure:"level"`
	Format   string `mapstructure:"format"`
	FilePath string `mapstructure:"file_path"` // Path to write log files (optional)
}

// DefaultLogConfig returns a default logger configuration.
func DefaultLogConfig() LogConfig {
	return LogConfig{
		Level:    "info",
		Format:   "text",
		FilePath: "",
	}
}

// MultiHandler implements slog.Handler and broadcasts logging calls to multiple handlers.
type MultiHandler struct {
	handlers []slog.Handler
}

// Enabled checks if any of the underlying handlers is enabled for the log level.
func (m *MultiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, h := range m.handlers {
		if h.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

// Handle clones records and broadcasts them to all enabled handlers.
//
//nolint:gocritic // copying slog.Record by value is required to implement the standard library slog.Handler interface
func (m *MultiHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, h := range m.handlers {
		if h.Enabled(ctx, r.Level) {
			if err := h.Handle(ctx, r.Clone()); err != nil {
				return err
			}
		}
	}
	return nil
}

// WithAttrs returns a new MultiHandler with attributes bound to all handlers.
func (m *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		newHandlers[i] = h.WithAttrs(attrs)
	}
	return &MultiHandler{handlers: newHandlers}
}

// WithGroup returns a new MultiHandler with a group name bound to all handlers.
func (m *MultiHandler) WithGroup(name string) slog.Handler {
	newHandlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		newHandlers[i] = h.WithGroup(name)
	}
	return &MultiHandler{handlers: newHandlers}
}

// SetupLogging initializes the global structured logger.
// If cfg.FilePath is specified, it duplicates logs to that file as structured JSON.
func SetupLogging(cfg LogConfig, w io.Writer) {
	if w == nil {
		w = os.Stderr
	}

	var level slog.Level
	if err := level.UnmarshalText([]byte(cfg.Level)); err != nil {
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	var consoleHandler slog.Handler = slog.NewTextHandler(w, opts)
	if LogFormat(strings.ToLower(cfg.Format)) == LogFormatJSON {
		consoleHandler = slog.NewJSONHandler(w, opts)
	}

	// If no log file path is specified, write only to the console stream
	if cfg.FilePath == "" {
		slog.SetDefault(slog.New(consoleHandler))
		return
	}

	// Ensure parent directory for log file exists
	dir := filepath.Dir(cfg.FilePath)
	if dir != "." && dir != "" {
		_ = os.MkdirAll(dir, 0o755)
	}

	// Open the log file in write-only append mode
	file, err := os.OpenFile(cfg.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		// Fallback to console-only log configuration on error
		slog.SetDefault(slog.New(consoleHandler))
		slog.Error("failed to initialize log file, falling back to stderr", slog.String("path", cfg.FilePath), slog.Any("error", err))
		return
	}

	// File logging is always structured JSON for parsing/auditing
	fileHandler := slog.NewJSONHandler(file, opts)

	slog.SetDefault(slog.New(&MultiHandler{
		handlers: []slog.Handler{consoleHandler, fileHandler},
	}))
}
