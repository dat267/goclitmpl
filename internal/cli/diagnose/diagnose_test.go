package diagnose

import (
	"testing"
)

func TestDiagnoseCommands(t *testing.T) {
	t.Run("info command wiring", func(t *testing.T) {
		cmd := NewDiagnoseInfoCmd()
		if cmd == nil {
			t.Fatal("expected NewDiagnoseInfoCmd to return a non-nil command")
		}
		if cmd.Use != "info" {
			t.Errorf("expected command Use to be 'info', got %q", cmd.Use)
		}
	})

	t.Run("check command wiring", func(t *testing.T) {
		cmd := NewDiagnoseCheckCmd()
		if cmd == nil {
			t.Fatal("expected NewDiagnoseCheckCmd to return a non-nil command")
		}
		if cmd.Use != "check [address]" {
			t.Errorf("expected command Use to be 'check [address]', got %q", cmd.Use)
		}
	})

	t.Run("run command wiring", func(t *testing.T) {
		cmd := NewDiagnoseRunCmd()
		if cmd == nil {
			t.Fatal("expected NewDiagnoseRunCmd to return a non-nil command")
		}
		if cmd.Use != "run" {
			t.Errorf("expected command Use to be 'run', got %q", cmd.Use)
		}
	})
}
