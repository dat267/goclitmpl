package cli

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/dat267/goclitmpl/pkg/version"
)

func TestVersionCommand(t *testing.T) {
	Version = "1.2.3"
	Commit = "abcdef"
	Date = "2026-07-02"

	t.Run("standard output", func(t *testing.T) {
		var out bytes.Buffer
		cmd := NewVersionCmd()
		cmd.SetOut(&out)
		cmd.SetArgs([]string{})

		if err := cmd.Execute(); err != nil {
			t.Fatalf("unexpected error executing version: %v", err)
		}

		output := out.String()
		if !strings.Contains(output, "version 1.2.3") {
			t.Errorf("expected output to contain 'version 1.2.3', got: %q", output)
		}
	})

	t.Run("json output", func(t *testing.T) {
		var out bytes.Buffer
		cmd := NewVersionCmd()
		cmd.SetOut(&out)
		cmd.SetArgs([]string{"--json"})

		if err := cmd.Execute(); err != nil {
			t.Fatalf("unexpected error executing version --json: %v", err)
		}

		var info version.Info
		if err := json.Unmarshal(out.Bytes(), &info); err != nil {
			t.Fatalf("failed to parse json output: %v", err)
		}
		if info.Version != "1.2.3" {
			t.Errorf("expected json version to be '1.2.3', got %q", info.Version)
		}
	})
}

func TestGreetCommand(t *testing.T) {
	t.Run("error missing name", func(t *testing.T) {
		cmd := NewGreetCmd()
		cmd.SetArgs([]string{})
		err := cmd.Execute()
		if err == nil {
			t.Error("expected error for missing name argument, got nil")
		}
	})

	t.Run("standard greeting", func(t *testing.T) {
		var out bytes.Buffer
		cmd := NewGreetCmd()
		cmd.SetOut(&out)
		cmd.SetArgs([]string{"Alice"})

		if err := cmd.Execute(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(out.String(), "Hello, Alice!") {
			t.Errorf("expected 'Hello, Alice!' in output, got: %q", out.String())
		}
	})

	t.Run("uppercase greeting", func(t *testing.T) {
		var out bytes.Buffer
		cmd := NewGreetCmd()
		cmd.SetOut(&out)
		cmd.SetArgs([]string{"Alice", "-u"})

		if err := cmd.Execute(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(out.String(), "HELLO, ALICE!") {
			t.Errorf("expected 'HELLO, ALICE!' in output, got: %q", out.String())
		}
	})
}

func TestConfigInitCommand(t *testing.T) {
	cmd := newConfigInitCmd()
	if cmd == nil {
		t.Fatal("expected newConfigInitCmd to return a non-nil command")
	}
	if cmd.Use != "init" {
		t.Errorf("expected Use='init', got %q", cmd.Use)
	}
}

func TestDiagnoseCommand(t *testing.T) {
	cmd := NewDiagnoseCmd()
	if cmd == nil {
		t.Fatal("expected NewDiagnoseCmd to return a non-nil command")
	}
	if cmd.Use != "diagnose" {
		t.Errorf("expected Use='diagnose', got %q", cmd.Use)
	}

	uses := map[string]bool{}
	for _, sub := range cmd.Commands() {
		uses[sub.Use] = true
	}
	for _, want := range []string{"info", "check [address]", "run"} {
		if !uses[want] {
			t.Errorf("expected subcommand %q to be registered under diagnose", want)
		}
	}
}

