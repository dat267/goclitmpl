package version

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestFormatVersion(t *testing.T) {
	t.Run("standard output", func(t *testing.T) {
		var out bytes.Buffer
		err := FormatVersion(&out, "1.2.3", "abcdef", "2026-07-02", false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := out.String()
		if !strings.Contains(output, "goclitmpl version 1.2.3") {
			t.Errorf("expected output to contain version '1.2.3', got: %q", output)
		}
		if !strings.Contains(output, "commit:     abcdef") {
			t.Errorf("expected output to contain commit 'abcdef', got: %q", output)
		}
	})

	t.Run("json output", func(t *testing.T) {
		var out bytes.Buffer
		err := FormatVersion(&out, "1.2.3", "abcdef", "2026-07-02", true)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
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
