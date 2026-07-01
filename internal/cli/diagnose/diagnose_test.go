package diagnose

import (
	"bytes"
	"net"
	"strings"
	"testing"
)

func TestDiagnoseInfoCommand(t *testing.T) {
	var out bytes.Buffer
	cmd := NewDiagnoseInfoCmd()
	cmd.SetOut(&out)
	cmd.SetArgs([]string{})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "System Runtime Specifications:") {
		t.Errorf("expected output to contain system specifications header, got: %q", output)
	}
	if !strings.Contains(output, "CPU Count:") {
		t.Errorf("expected output to contain CPU Count metrics, got: %q", output)
	}
}

func TestDiagnoseCheckCommand(t *testing.T) {
	// 1. Setup a local TCP listener to mock remote endpoint
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to spin up mock TCP listener: %v", err)
	}
	defer l.Close()

	mockAddress := l.Addr().String()

	t.Run("successful connection to mock listener", func(t *testing.T) {
		var out bytes.Buffer
		cmd := NewDiagnoseCmd()
		cmd.SetOut(&out)
		cmd.SetArgs([]string{"check", mockAddress})

		if err := cmd.Execute(); err != nil {
			t.Fatalf("unexpected check error: %v", err)
		}

		output := out.String()
		expectedMsg := "Successfully connected to " + mockAddress
		if !strings.Contains(output, expectedMsg) {
			t.Errorf("expected output to contain %q, got: %q", expectedMsg, output)
		}
	})

	t.Run("failure on invalid connection target", func(t *testing.T) {
		cmd := NewDiagnoseCmd()
		// Try connecting to a port that is highly likely closed/unreachable
		cmd.SetArgs([]string{"check", "127.0.0.1:9999"})

		err := cmd.Execute()
		if err == nil {
			t.Error("expected connection check error on invalid address, got nil")
		}
	})
}

func TestDiagnoseRunCommand(t *testing.T) {
	var out bytes.Buffer
	cmd := NewDiagnoseCmd()
	cmd.SetOut(&out)
	// Override timeout to reduce test runtime if remote targets fail
	cmd.SetArgs([]string{"--timeout", "100ms", "run"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected execution error: %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "Executing Diagnostics Suite...") {
		t.Errorf("expected output to contain suite header, got: %q", output)
	}
	if !strings.Contains(output, "Remote Endpoints Network Probe:") {
		t.Errorf("expected output to contain network probe section, got: %q", output)
	}
	if !strings.Contains(output, "Diagnostics run completed successfully!") {
		t.Errorf("expected output to contain completion suffix, got: %q", output)
	}
}
