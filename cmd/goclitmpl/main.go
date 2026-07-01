// Package main is the entrypoint of the application.
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/dat267/goclitmpl/internal/cli"
)

func main() {
	os.Exit(run())
}

func run() int {
	// Create context that listens for the interrupt signals from the OS.
	// This ensures that any subcommand can react to cancellation signals gracefully
	// by checking cmd.Context().Done().
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	return cli.Execute(ctx)
}
