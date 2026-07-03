// Package diagnose provides system and runtime info reporting.
package diagnose

import (
	"fmt"
	"io"
	"runtime"
)

// PrintInfo writes the hardware and runtime specifications to the provided writer.
func PrintInfo(w io.Writer) {
	fmt.Fprintf(w, "System Runtime Specifications:\n")
	fmt.Fprintf(w, "  Operating System: %s\n", runtime.GOOS)
	fmt.Fprintf(w, "  Architecture:     %s\n", runtime.GOARCH)
	fmt.Fprintf(w, "  Go Version:       %s\n", runtime.Version())
	fmt.Fprintf(w, "  CPU Count:        %d\n", runtime.NumCPU())
}
