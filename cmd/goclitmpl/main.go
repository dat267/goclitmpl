package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/dat267/goclitmpl/internal/cli"
)

func main() {
	// Create context that listens for the interrupt signals from the OS.
	// This ensures that any subcommand can react to cancellation signals gracefully
	// by checking cmd.Context().Done().
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	code := cli.Execute(ctx)
	os.Exit(code)
}
