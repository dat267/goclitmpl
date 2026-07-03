// Package build holds variables injected at compile time via -ldflags.
// These are kept in a dedicated package so any part of the application
// can import build metadata without importing the CLI layer.
package build

// Variables injected by -ldflags at build time.
// Defaults are used when building without the Makefile or build script.
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)
