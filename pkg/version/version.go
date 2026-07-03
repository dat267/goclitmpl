// Package version provides running logic and formatting for version information outputs.
package version

import (
	"encoding/json"
	"fmt"
	"io"
	"runtime"
)

// Info holds version metrics for structured output.
type Info struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	Date      string `json:"date"`
	GoVersion string `json:"goVersion"`
	OS        string `json:"os"`
	Arch      string `json:"arch"`
}

// FormatVersion writes version info to the provided writer in JSON or plain text format.
func FormatVersion(w io.Writer, buildVersion, buildCommit, buildDate string, jsonOutput bool) error {
	info := Info{
		Version:   buildVersion,
		Commit:    buildCommit,
		Date:      buildDate,
		GoVersion: runtime.Version(),
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
	}

	if jsonOutput {
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(info); err != nil {
			return fmt.Errorf("failed to marshal version details: %w", err)
		}
		return nil
	}

	fmt.Fprintf(w, "goclitmpl version %s\n", info.Version)
	fmt.Fprintf(w, "  commit:     %s\n", info.Commit)
	fmt.Fprintf(w, "  built at:   %s\n", info.Date)
	fmt.Fprintf(w, "  go version: %s\n", info.GoVersion)
	fmt.Fprintf(w, "  os/arch:    %s/%s\n", info.OS, info.Arch)

	return nil
}
