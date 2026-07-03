package diagnose

import (
	"fmt"
	"io"
	"log/slog"
	"net"
	"runtime"
	"time"
)

// RunSuite executes a full diagnostics check, writing progress and results to the writer.
func RunSuite(w io.Writer, timeout time.Duration) {
	slog.Info("starting full diagnostic suite run", slog.Duration("timeout", timeout))
	fmt.Fprintf(w, "Executing Diagnostics Suite...\n\n")

	// Check 1: Hardware specs
	fmt.Fprintf(w, "[1/2] Hardware Specification Profile:\n")
	fmt.Fprintf(w, "  OS/Arch:      %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Fprintf(w, "  CPU Cores:    %d\n", runtime.NumCPU())
	fmt.Fprintf(w, "  Go Version:   %s\n\n", runtime.Version())

	// Check 2: Core endpoints connectivity
	targets := []string{"google.com:80", "github.com:80"}
	fmt.Fprintf(w, "[2/2] Remote Endpoints Network Probe:\n")

	for _, target := range targets {
		start := time.Now()
		conn, err := net.DialTimeout("tcp", target, timeout)
		duration := time.Since(start)

		if err != nil {
			slog.Warn("endpoint probe failed", slog.String("target", target), slog.Any("error", err))
			fmt.Fprintf(w, "  %-15s -> FAILED (%v)\n", target, err)
		} else {
			conn.Close()
			slog.Debug("endpoint probe succeeded", slog.String("target", target), slog.Duration("latency", duration))
			fmt.Fprintf(w, "  %-15s -> SUCCESS (latency: %v)\n", target, duration.Round(time.Millisecond))
		}
	}

	fmt.Fprintf(w, "\nDiagnostics run completed successfully!\n")
}
