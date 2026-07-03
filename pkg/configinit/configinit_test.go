package configinit

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dat267/goclitmpl/internal/config"
)

func TestInitialize(t *testing.T) {
	// Create a temporary directory for isolation
	tempHome, err := os.MkdirTemp("", "configinit-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempHome)

	// Backup and override HOME and XDG_CONFIG_HOME environment variables
	oldHome := os.Getenv("HOME")
	oldXdg := os.Getenv("XDG_CONFIG_HOME")
	t.Setenv("HOME", tempHome)
	t.Setenv("XDG_CONFIG_HOME", tempHome)
	defer func() {
		// Restore in case Setenv doesn't (though Go's t.Setenv does it automatically)
		_ = os.Setenv("HOME", oldHome)
		_ = os.Setenv("XDG_CONFIG_HOME", oldXdg)
	}()

	expectedConfigPath := filepath.Join(tempHome, config.AppName, "config.yaml")

	t.Run("initialize config file successfully", func(t *testing.T) {
		var out bytes.Buffer
		err = Initialize(&out, false)
		if err != nil {
			t.Fatalf("unexpected error running Initialize: %v", err)
		}

		output := out.String()
		if !strings.Contains(output, "Successfully initialized default configuration at") {
			t.Errorf("expected output to contain success message, got: %q", output)
		}

		// Verify file was written to disk
		_, statErr := os.Stat(expectedConfigPath)
		if os.IsNotExist(statErr) {
			t.Errorf("expected config file to be created at %s, but it does not exist", expectedConfigPath)
		}
	})

	t.Run("fail on existing file without force", func(t *testing.T) {
		var out bytes.Buffer
		err = Initialize(&out, false)
		if err == nil {
			t.Error("expected error when trying to write over existing config file, got nil")
		}
		if !strings.Contains(err.Error(), "config file already exists") {
			t.Errorf("expected error message to complain about existing file, got: %v", err)
		}
	})

	t.Run("succeed on existing file with force flag", func(t *testing.T) {
		var out bytes.Buffer
		err = Initialize(&out, true)
		if err != nil {
			t.Fatalf("unexpected error running Initialize with force: %v", err)
		}

		output := out.String()
		if !strings.Contains(output, "Successfully initialized default configuration") {
			t.Errorf("expected output to contain success message, got: %q", output)
		}
	})
}
