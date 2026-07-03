package diagnose

import (
	"bytes"
	"net"
	"strings"
	"testing"
	"time"
)

func TestPrintInfo(t *testing.T) {
	var out bytes.Buffer
	PrintInfo(&out)

	output := out.String()
	if !strings.Contains(output, "System Runtime Specifications:") {
		t.Errorf("expected output to contain system specifications header, got: %q", output)
	}
	if !strings.Contains(output, "CPU Count:") {
		t.Errorf("expected output to contain CPU Count metrics, got: %q", output)
	}
}

func TestCheckAddress(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to spin up mock TCP listener: %v", err)
	}
	defer l.Close()

	mockAddress := l.Addr().String()

	t.Run("successful connection to mock listener", func(t *testing.T) {
		err := CheckAddress(mockAddress, 100*time.Millisecond)
		if err != nil {
			t.Fatalf("unexpected check error: %v", err)
		}
	})

	t.Run("failure on invalid connection target", func(t *testing.T) {
		err := CheckAddress("127.0.0.1:9999", 50*time.Millisecond)
		if err == nil {
			t.Error("expected connection check error on invalid address, got nil")
		}
	})
}

func TestRunSuite(t *testing.T) {
	var out bytes.Buffer
	RunSuite(&out, 50*time.Millisecond)

	output := out.String()
	if !strings.Contains(output, "Executing Diagnostics Suite...") {
		t.Errorf("expected output to contain suite header, got: %q", output)
	}
	if !strings.Contains(output, "Remote Endpoints Network Probe:") {
		t.Errorf("expected output to contain network probe section, got: %q", output)
	}
}
