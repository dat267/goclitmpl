package config

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	appconfig "github.com/dat267/goclitmpl/internal/config"
)

func TestInitialize(t *testing.T) {
	tempHome := t.TempDir()
	t.Setenv("HOME", tempHome)
	t.Setenv("XDG_CONFIG_HOME", tempHome)

	expectedConfigPath := filepath.Join(tempHome, appconfig.AppName, "config.yaml")

	t.Run("initialize config file successfully", func(t *testing.T) {
		var out bytes.Buffer
		err := Initialize(&out, false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(out.String(), "Successfully initialized default configuration at") {
			t.Errorf("expected success message, got: %q", out.String())
		}
		if _, err = os.Stat(expectedConfigPath); os.IsNotExist(err) {
			t.Errorf("expected config file at %s, but not found", expectedConfigPath)
		}
	})

	t.Run("fail on existing file without force", func(t *testing.T) {
		err := Initialize(&bytes.Buffer{}, false)
		if err == nil {
			t.Fatal("expected error when overwriting without --force, got nil")
		}
		if !strings.Contains(err.Error(), "config file already exists") {
			t.Errorf("unexpected error message: %v", err)
		}
	})

	t.Run("succeed on existing file with force", func(t *testing.T) {
		var out bytes.Buffer
		err := Initialize(&out, true)
		if err != nil {
			t.Fatalf("unexpected error with --force: %v", err)
		}
		if !strings.Contains(out.String(), "Successfully initialized default configuration") {
			t.Errorf("expected success message, got: %q", out.String())
		}
	})
}
