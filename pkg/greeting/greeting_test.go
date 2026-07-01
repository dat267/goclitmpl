package greeting

import (
	"context"
	"io"
	"log/slog"
	"testing"
)

func TestGreetingFormat(t *testing.T) {
	// Create a discard logger for unit tests to keep clean test outputs
	discardLogger := slog.New(slog.NewTextHandler(io.Discard, nil))
	cfg := Config{Env: "development"}
	svc := NewService(discardLogger, cfg)

	t.Run("standard greeting", func(t *testing.T) {
		got, err := svc.Format(context.Background(), "Alice", false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := "Hello, Alice! (env: development)"
		if got != expected {
			t.Errorf("expected %q, got %q", expected, got)
		}
	})

	t.Run("uppercase greeting", func(t *testing.T) {
		got, err := svc.Format(context.Background(), "Bob", true)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := "HELLO, BOB! (ENV: DEVELOPMENT)"
		if got != expected {
			t.Errorf("expected %q, got %q", expected, got)
		}
	})

	t.Run("empty name validation error", func(t *testing.T) {
		_, err := svc.Format(context.Background(), "   ", false)
		if err == nil {
			t.Error("expected error for empty name parameter, got nil")
		}
	})

	t.Run("context cancellation propagation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		_, err := svc.Format(ctx, "Charlie", false)
		if err == nil {
			t.Error("expected context canceled error, got nil")
		}
		if err != context.Canceled {
			t.Errorf("expected context.Canceled, got: %v", err)
		}
	})
}
