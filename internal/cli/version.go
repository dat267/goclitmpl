package cli

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	// Version is the current version of the application, injected at build time.
	Version = "dev"
	// Commit is the git commit hash injected at build time.
	Commit = "none"
	// Date is the build date injected at build time.
	Date = "unknown"
)

// Info holds details about the application build.
type Info struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	Date      string `json:"date"`
	GoVersion string `json:"go_version"`
	OS        string `json:"os"`
	Arch      string `json:"arch"`
}

// GetInfo returns the build information.
func GetInfo() Info {
	return Info{
		Version:   Version,
		Commit:    Commit,
		Date:      Date,
		GoVersion: runtime.Version(),
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
	}
}

// NewVersionCmd creates the version command.
func NewVersionCmd() *cobra.Command {
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version information",
		Long:  `Print compile-time version and build metadata.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			info := GetInfo()
			out := cmd.OutOrStdout()

			if jsonOutput {
				encoder := json.NewEncoder(out)
				encoder.SetIndent("", "  ")
				return encoder.Encode(info)
			}

			fmt.Fprintf(out, "goclitmpl version %s\n", info.Version)
			fmt.Fprintf(out, "  commit:     %s\n", info.Commit)
			fmt.Fprintf(out, "  built at:   %s\n", info.Date)
			fmt.Fprintf(out, "  go version: %s\n", info.GoVersion)
			fmt.Fprintf(out, "  os/arch:    %s/%s\n", info.OS, info.Arch)

			return nil
		},
	}

	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output version information in JSON format")

	return cmd
}
