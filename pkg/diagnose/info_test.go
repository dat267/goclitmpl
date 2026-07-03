package diagnose

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrintInfo(t *testing.T) {
	var out bytes.Buffer
	PrintInfo(&out)
	output := out.String()

	if !strings.Contains(output, "System Runtime Specifications:") {
		t.Errorf("expected system spec header, got: %q", output)
	}
	if !strings.Contains(output, "CPU Count:") {
		t.Errorf("expected CPU Count metric, got: %q", output)
	}
}
