package config

import (
	"bytes"
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := NewDefaultConfig()
	if cfg.App.Port != 8080 {
		t.Errorf("expected default port to be 8080, got %d", cfg.App.Port)
	}
	if cfg.Log.Level != "info" {
		t.Errorf("expected default log level to be 'info', got %q", cfg.Log.Level)
	}
}

func TestConfigValidation(t *testing.T) {
	cfg := NewDefaultConfig()

	// Valid config
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected default config to be valid, got error: %v", err)
	}

	// Invalid port
	cfg.App.Port = 99999
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for invalid port 99999, got nil")
	}

	// Invalid log level
	cfg = NewDefaultConfig()
	cfg.Log.Level = "invalid"
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for invalid log level, got nil")
	}

	// Invalid log format
	cfg = NewDefaultConfig()
	cfg.Log.Format = "xml"
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for invalid log format, got nil")
	}
}

func TestLoadWithEnvOverrides(t *testing.T) {
	// Set environment variables
	os.Setenv("GOCLITMPL_APP_PORT", "9090")
	os.Setenv("GOCLITMPL_LOG_LEVEL", "debug")
	defer func() {
		os.Unsetenv("GOCLITMPL_APP_PORT")
		os.Unsetenv("GOCLITMPL_LOG_LEVEL")
	}()

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if cfg.App.Port != 9090 {
		t.Errorf("expected port to be overridden by env variable to 9090, got %d", cfg.App.Port)
	}

	if cfg.Log.Level != "debug" {
		t.Errorf("expected log level to be overridden by env variable to 'debug', got %q", cfg.Log.Level)
	}
}

func TestLoadWithFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "config-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	configContent := `
app:
  port: 7070
  env: staging
log:
  level: warn
  format: json
`
	configFilePath := filepath.Join(tempDir, "test-config.yaml")
	if err := os.WriteFile(configFilePath, []byte(configContent), 0o644); err != nil {
		t.Fatalf("failed to write temp config file: %v", err)
	}

	cfg, err := Load(configFilePath)
	if err != nil {
		t.Fatalf("failed to load config from file: %v", err)
	}

	if cfg.App.Port != 7070 {
		t.Errorf("expected port to be 7070 from file, got %d", cfg.App.Port)
	}
	if cfg.App.Env != "staging" {
		t.Errorf("expected env to be 'staging' from file, got %q", cfg.App.Env)
	}
	if cfg.Log.Level != "warn" {
		t.Errorf("expected log level to be 'warn' from file, got %q", cfg.Log.Level)
	}
	if cfg.Log.Format != "json" {
		t.Errorf("expected log format to be 'json' from file, got %q", cfg.Log.Format)
	}
}

func TestLoggingSetupWithFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "logging-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	logFilePath := filepath.Join(tempDir, "app.log")
	cfg := LogConfig{
		Level:    "debug",
		Format:   "text",
		FilePath: logFilePath,
	}

	var consoleOut bytes.Buffer
	SetupLogging(cfg, &consoleOut)

	// Log messages
	slog.Info("console text and file json message", slog.String("key", "value"))

	// 1. Verify console output (should be human-readable text)
	consoleStr := consoleOut.String()
	if !strings.Contains(consoleStr, `level=INFO msg="console text and file json message" key=value`) {
		t.Errorf("unexpected console log format: %q", consoleStr)
	}

	// 2. Verify file output exists
	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		t.Fatalf("log file was not created: %s", logFilePath)
	}

	// 3. Verify file output content (should be JSON)
	fileBytes, err := os.ReadFile(logFilePath)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}
	fileStr := string(fileBytes)

	if !strings.Contains(fileStr, `"level":"INFO"`) ||
		!strings.Contains(fileStr, `"msg":"console text and file json message"`) ||
		!strings.Contains(fileStr, `"key":"value"`) {
		t.Errorf("unexpected file log format: %q", fileStr)
	}
}

func TestMultiHandlerMethods(t *testing.T) {
	var out1, out2 bytes.Buffer
	opts := &slog.HandlerOptions{Level: slog.LevelInfo}
	h1 := slog.NewTextHandler(&out1, opts)
	h2 := slog.NewJSONHandler(&out2, opts)

	mh := &MultiHandler{
		handlers: []slog.Handler{h1, h2},
	}

	// Test Enabled
	if !mh.Enabled(context.Background(), slog.LevelInfo) {
		t.Error("expected MultiHandler to be enabled for LevelInfo")
	}
	if mh.Enabled(context.Background(), slog.LevelDebug) {
		t.Error("expected MultiHandler to be disabled for LevelDebug")
	}

	// Test WithAttrs & WithGroup formatting compilation
	mhWithAttrs := mh.WithAttrs([]slog.Attr{slog.String("a", "b")})
	if len(mhWithAttrs.(*MultiHandler).handlers) != 2 {
		t.Error("expected WithAttrs handlers size to match")
	}

	mhWithGroup := mh.WithGroup("g")
	if len(mhWithGroup.(*MultiHandler).handlers) != 2 {
		t.Error("expected WithGroup handlers size to match")
	}
}
