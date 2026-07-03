// Package diagnose provides system and connection health tools.
package diagnose

import (
	"fmt"
	"io"
	"log/slog"
	"net"
	"runtime"
	"time"
)

// PrintInfo writes the hardware and runtime specifications to the provided writer.
func PrintInfo(w io.Writer) {
	fmt.Fprintf(w, "System Runtime Specifications:\n")
	fmt.Fprintf(w, "  Operating System: %s\n", runtime.GOOS)
	fmt.Fprintf(w, "  Architecture:     %s\n", runtime.GOARCH)
	fmt.Fprintf(w, "  Go Version:       %s\n", runtime.Version())
	fmt.Fprintf(w, "  CPU Count:        %d\n", runtime.NumCPU())
}

// CheckAddress validates connection reachability to a specific address with a timeout.
func CheckAddress(address string, timeout time.Duration) error {
	slog.Debug("checking network reachability", slog.String("address", address), slog.Duration("timeout", timeout))
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return err
	}
	conn.Close()
	return nil
}

// RunSuite executes a full diagnostics check, writing progress and latency details to the writer.
func RunSuite(w io.Writer, timeout time.Duration) {
	slog.Info("starting full diagnostic suite run", slog.Duration("timeout", timeout))
	fmt.Fprintf(w, "Executing Diagnostics Suite...\n\n")

	// Check 1: Hardware specs
	fmt.Fprintf(w, "[1/2] Hardware Specification Profile:\n")
	fmt.Fprintf(w, "  OS/Arch:      %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Fprintf(w, "  CPU Cores:    %d\n", runtime.NumCPU())
	fmt.Fprintf(w, "  Go Version:   %s\n\n", runtime.Version())

	// Check 2: Core endpoints connection
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
