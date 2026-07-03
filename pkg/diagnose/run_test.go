package diagnose

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestRunSuite(t *testing.T) {
	var out bytes.Buffer
	RunSuite(&out, 50*time.Millisecond)
	output := out.String()

	if !strings.Contains(output, "Executing Diagnostics Suite...") {
		t.Errorf("expected suite header, got: %q", output)
	}
	if !strings.Contains(output, "Remote Endpoints Network Probe:") {
		t.Errorf("expected network probe section, got: %q", output)
	}
}
