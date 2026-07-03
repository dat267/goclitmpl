package diagnose

import (
	"net"
	"testing"
	"time"
)

func TestCheckAddress(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start mock listener: %v", err)
	}
	defer l.Close()

	t.Run("successful connection to mock listener", func(t *testing.T) {
		if err := CheckAddress(l.Addr().String(), 100*time.Millisecond); err != nil {
			t.Fatalf("expected success, got: %v", err)
		}
	})

	t.Run("failure on unreachable address", func(t *testing.T) {
		if err := CheckAddress("127.0.0.1:9999", 50*time.Millisecond); err == nil {
			t.Error("expected error for unreachable address, got nil")
		}
	})
}
