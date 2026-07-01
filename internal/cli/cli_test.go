package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dat267/goclitmpl/internal/config"
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
		if !strings.Contains(output, "commit:     abcdef") {
			t.Errorf("expected output to contain commit 'abcdef', got: %q", output)
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

		var info Info
		if err := json.Unmarshal(out.Bytes(), &info); err != nil {
			t.Fatalf("failed to parse json output: %v", err)
		}

		if info.Version != "1.2.3" {
			t.Errorf("expected json version to be '1.2.3', got %q", info.Version)
		}
		if info.Commit != "abcdef" {
			t.Errorf("expected json commit to be 'abcdef', got %q", info.Commit)
		}
	})
}

func TestGreetCommand(t *testing.T) {
	t.Run("error missing name", func(t *testing.T) {
		var out bytes.Buffer
		cmd := NewGreetCmd()
		cmd.SetOut(&out)
		cmd.SetArgs([]string{})

		err := cmd.Execute()
		if err == nil {
			t.Error("expected error for missing name argument, got nil")
		}
		if !strings.Contains(err.Error(), "accepts 1 arg(s), received 0") {
			t.Errorf("expected error message to contain 'accepts 1 arg(s), received 0', got: %v", err)
		}
	})

	t.Run("error too many arguments", func(t *testing.T) {
		var out bytes.Buffer
		cmd := NewGreetCmd()
		cmd.SetOut(&out)
		cmd.SetArgs([]string{"Alice", "Bob"})

		err := cmd.Execute()
		if err == nil {
			t.Error("expected error for too many arguments, got nil")
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

		output := out.String()
		if !strings.Contains(output, "Hello, Alice!") {
			t.Errorf("expected greeting output to contain 'Hello, Alice!', got: %q", output)
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

		output := out.String()
		if !strings.Contains(output, "HELLO, ALICE!") {
			t.Errorf("expected uppercase output to contain 'HELLO, ALICE!', got: %q", output)
		}
	})
}

func TestConfigInitCommand(t *testing.T) {
	// Create a temporary directory for isolation
	tempHome, err := os.MkdirTemp("", "config-init-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempHome)

	// Backup and override HOME and XDG_CONFIG_HOME environment variables
	oldHome := os.Getenv("HOME")
	oldXdg := os.Getenv("XDG_CONFIG_HOME")
	os.Setenv("HOME", tempHome)
	os.Setenv("XDG_CONFIG_HOME", tempHome)
	defer func() {
		os.Setenv("HOME", oldHome)
		os.Setenv("XDG_CONFIG_HOME", oldXdg)
	}()

	expectedConfigPath := filepath.Join(tempHome, config.AppName, "config.yaml")

	t.Run("initialize config file successfully", func(t *testing.T) {
		var out bytes.Buffer
		cmd := NewConfigInitCmd()
		cmd.SetOut(&out)
		cmd.SetArgs([]string{})

		if err := cmd.Execute(); err != nil {
			t.Fatalf("unexpected error running config init: %v", err)
		}

		// Verify output message
		output := out.String()
		if !strings.Contains(output, "Successfully initialized default configuration at") {
			t.Errorf("expected output to contain success message, got: %q", output)
		}

		// Verify file was written to disk
		if _, err := os.Stat(expectedConfigPath); os.IsNotExist(err) {
			t.Errorf("expected config file to be created at %s, but it does not exist", expectedConfigPath)
		}
	})

	t.Run("fail on existing file without force", func(t *testing.T) {
		cmd := NewConfigInitCmd()
		cmd.SetArgs([]string{})

		err := cmd.Execute()
		if err == nil {
			t.Error("expected error when trying to write over existing config file, got nil")
		}
		if !strings.Contains(err.Error(), "config file already exists") {
			t.Errorf("expected error message to complain about existing file, got: %v", err)
		}
	})

	t.Run("succeed on existing file with force flag", func(t *testing.T) {
		var out bytes.Buffer
		cmd := NewConfigInitCmd()
		cmd.SetOut(&out)
		cmd.SetArgs([]string{"--force"})

		if err := cmd.Execute(); err != nil {
			t.Fatalf("unexpected error running config init with force: %v", err)
		}

		output := out.String()
		if !strings.Contains(output, "Successfully initialized default configuration") {
			t.Errorf("expected output to contain success message, got: %q", output)
		}
	})
}
