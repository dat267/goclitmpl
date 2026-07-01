# Developer Guide: Extending `goclitmpl`

This guide explains how to extend this CLI template project in the future.

---

## 1. Adding a New Subcommand (Explicit Registration)

For fine-grained control, safety from side-effects, and optimal testability, this template uses **explicit command registration**.

### Step 1: Create a command file
Create a new file in the `internal/cli/` directory, for example, `internal/cli/import.go`:

```go
package cli

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
)

// NewImportCmd creates and configures the subcommand
func NewImportCmd() *cobra.Command {
	var filePath string
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "import [target]",
		Short: "Imports datasets into the application",
		Long:  `Reads data from a local file and processes it through core services.`,
		Args:  cobra.ExactArgs(1), // Enforces exactly one positional argument
		RunE: func(cmd *cobra.Command, args []string) error {
			target := args[0]
			
			// 1. Get the parsed global configuration
			cfg := GetConfig(cmd.Context())

			slog.Info("starting import task", 
				slog.String("target", target), 
				slog.String("file", filePath),
				slog.Bool("dry_run", dryRun),
			)

			// 2. Call your business logic package (e.g. from internal/service)
			// err := service.ProcessImport(cmd.Context(), target, filePath, dryRun, cfg)
			// if err != nil { return err }

			// 3. Print output to stdout using the command's built-in writer stream
			fmt.Fprintf(cmd.OutOrStdout(), "Successfully imported %s!\n", target)
			return nil
		},
	}

	// Define command flags
	cmd.Flags().StringVarP(&filePath, "file", "f", "", "path to the source file (required)")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "simulate import without saving changes")
	
	// Mark flag as required
	_ = cmd.MarkFlagRequired("file")

	return cmd
}
```

### Step 2: Register the Subcommand in `root.go`
To mount the subcommand, open [internal/cli/root.go](file:///home/dat/repos/goclitmpl/internal/cli/root.go) and add your constructor explicitly inside `NewRootCmd()`:

```go
	// Add subcommands explicitly for fine-grained control
	cmd.AddCommand(
		NewVersionCmd(),
		NewGreetCmd(),
		NewImportCmd(), // <-- Add your subcommand constructor here
	)
```

---

## 2. Organizing Nested Subcommands

For subcommands of subcommands (e.g., `goclitmpl diagnose run`), do not clutter `root.go`. Organize them in a sub-package inside `internal/cli/`. This has already been implemented as an example inside [internal/cli/diagnose/](file:///home/dat/repos/goclitmpl/internal/cli/diagnose/):

```text
internal/cli/diagnose/
├── diagnose.go        # Parent command "diagnose" (defines --timeout persistent flag)
├── diagnose_info.go   # Leaf command "diagnose info" (prints CPU/OS specs)
├── diagnose_check.go  # Leaf command "diagnose check" (probes network target)
└── diagnose_run.go    # Leaf command "diagnose run" (runs all diagnostics)
```

### Best Practices Used:
1. **Command Factory Chain (Lazy Mounting)**:
   * `NewDiagnoseCmd()` constructs the root `diagnose` command and mounts `NewDiagnoseInfoCmd()`, `NewDiagnoseCheckCmd()`, and `NewDiagnoseRunCmd()`.
   * `root.go` only imports `diagnose` and mounts `diagnose.NewDiagnoseCmd()`.
2. **Persistent Flag Inheritance**:
   * Defining the `--timeout` flag as `PersistentFlags().Duration(...)` on the root `diagnose` command makes it automatically visible and parseable in child commands (`check` and `run`).
   * Retrieve inherited persistent flags inside child commands using `cmd.Flags().GetDuration("timeout")`.
3. **Domain Isolation**:
   * Keeps related CLI flags and diagnostics logic scoped cleanly inside their own package folder.

---

## 3. Configuration Initialization (`config init`)

To make it easy for users to get started, the template includes a `config init` subcommand. This subcommand generates a well-documented YAML configuration file in the user's home configuration directory (e.g., `~/.config/goclitmpl/config.yaml`).

### Implementation Strategy:
1. **Commented Template File**: The default settings and documentation comments are stored in [internal/config/default_config.yaml](file:///home/dat/repos/goclitmpl/internal/config/default_config.yaml).
2. **Go Embedding**: The YAML file is compiled directly into the binary using Go's `//go:embed` directive inside [internal/config/config.go](file:///home/dat/repos/goclitmpl/internal/config/config.go).
3. **Overwrite Prevention**: By default, `config init` will fail if a configuration file already exists in the destination folder to prevent wiping out user customizations. Passing the `--force` / `-f` flag will bypass this check.

To add new defaults to the generated configuration file, update [internal/config/default_config.yaml](file:///home/dat/repos/goclitmpl/internal/config/default_config.yaml) and re-compile.

---

## 4. Extending Configuration Settings

Configuration management is mapped to structs using Viper and standard validation.

### Step 1: Update config structs
Open [internal/config/config.go](file:///home/dat/repos/goclitmpl/internal/config/config.go) and add parameters under `AppConfig` or create custom sections:

```go
type AppConfig struct {
	Env    string `mapstructure:"env"`
	Debug  bool   `mapstructure:"debug"`
	Port   int    `mapstructure:"port"`
	DBHost string `mapstructure:"db_host"` // <-- Add your new configuration field here
}
```

### Step 2: Bind defaults and configure Viper
In the same file, set fallbacks and defaults in `NewDefaultConfig()` and `Load()`:

```go
// 1. In NewDefaultConfig()
return &Config{
    Log: DefaultLogConfig(),
    App: AppConfig{
        Env:    "production",
        Debug:  false,
        Port:   8080,
        DBHost: "127.0.0.1", // <-- Add default value here
    },
}

// 2. In Load()
v.SetDefault("app.db_host", "127.0.0.1") // <-- Configures Viper default fallback
```

### Step 3: Add Validation rules
Validate parameters inside the `Validate()` method of `internal/config/config.go`:

```go
if c.App.DBHost == "" {
    return fmt.Errorf("app.db_host cannot be empty")
}
```

### Step 4: Access variables from CLI commands
The loader parses and merges parameters from the following targets. For example, to override `app.db_host`:
1. **Config file (`config.yaml`)**:
   ```yaml
   app:
     db_host: "prod-db.example.com"
   ```
2. **Environment variables**: Use prefix `GOCLITMPL_` and double-underscore for nesting levels:
   ```bash
   export GOCLITMPL_APP_DB_HOST="staging-db.example.com"
   ```
3. **Retrieving configuration**: Retrieve structure through context inside your subcommands:
   ```go
   cfg := GetConfig(cmd.Context())
   dbHost := cfg.App.DBHost
   ```

### Step 5: Passing Configuration to Business Logic (Best Practice)
While we use Go `context.Context` to pass configuration variables between Cobra command hooks (like `PersistentPreRunE` to subcommand `RunE`), **passing configurations implicitly through Context to your core business logic is a Go anti-pattern**.

To keep your core business logic packages decoupled from the CLI implementation details, follow these four design pillars:
1. **Pass config explicitly**: Map configurations to package-specific config structures during initialization.
2. **Inject dependencies**: Inject loggers (`*slog.Logger`) or clients explicitly during constructor calls.
3. **Propagate Context**: Pass `context.Context` as the very first argument of execution functions for timeouts and cancellations.
4. **Propagate Errors**: Return both `(result, error)` to handle business validation and execution failures cleanly.

**Example**:
```go
// In internal/cli/greet.go (CLI Package)
RunE: func(cmd *cobra.Command, args []string) error {
    cfg := GetConfig(cmd.Context()) // 1. Extract from CLI Context
    
    pkgCfg := greeting.Config{
        Env: cfg.App.Env,           // 2. Map config explicitly
    }
    
    // 3. Inject explicit dependencies (slog.Default())
    svc := greeting.NewService(slog.Default(), pkgCfg)
    
    // 4. Pass Context and handle error returns
    message, err := svc.Format(cmd.Context(), args[0], uppercase)
    if err != nil {
        return fmt.Errorf("greeting failed: %w", err)
    }
    
    fmt.Fprintln(cmd.OutOrStdout(), message)
    return nil
}

// In pkg/greeting/greeting.go (Decoupled Business Logic Package)
package greeting

type Config struct {
    Env string
}

type Service struct {
    logger *slog.Logger
    env    string
}

// Constructor takes explicit parameters. It has no dependency on Cobra or Viper.
func NewService(logger *slog.Logger, cfg Config) *Service {
    return &Service{logger: logger, env: cfg.Env}
}

// Format takes a standard Context and returns standard (string, error) output.
func (s *Service) Format(ctx context.Context, name string, uppercase bool) (string, error) {
    if err := ctx.Err(); err != nil {
        return "", err
    }
    if name == "" {
        return "", errors.New("name cannot be empty")
    }
    // Business logic execution...
    return "Hello, " + name, nil
}
```
This preserves package isolation, making your business logic highly testable without needing mock contexts, Cobra commands, or Viper setups.

---

## 5. File Logging Best Practices

### Dual-Destination Logging (`MultiHandler`)
To allow users to debug failures during execution without cluttering the console output, CLIs should write detailed auditing logs to a file while showing clean, minimal outputs on standard error.

**Best Practice**: Implement a broadcast logging handler (`MultiHandler`) using standard Go `slog` to split streams:
1. **Console (Stderr)**: Write human-readable, colorized, or basic text logs (`slog.NewTextHandler()`).
2. **Auditing File**: Write highly detailed structured JSON logs (`slog.NewJSONHandler()`) containing timestamps, component domains, and error traces.

This design pattern is implemented in [internal/config/logging.go](file:///home/dat/repos/goclitmpl/internal/config/logging.go), which safely handles log file path directory creations and falls back gracefully to console-only logging if file permissions fail.

---

## 6. Best Practices for Command Output & Logging

To keep commands automation-friendly and compliant with standard UNIX behaviors:
* **Stdout**: Write actual data outputs (JSON strings, text lines, data pipelines) using `fmt.Fprint(cmd.OutOrStdout(), ...)` (which will inherit standard out redirects correctly).
* **Stderr**: Write log statements, debug outputs, warnings, and system status information using `slog.Info`, `slog.Error`, etc. (which output to standard error).
* **OS Signals & Cancellation**: Check context completion `cmd.Context().Done()` during long-running tasks to support graceful termination on Ctrl+C.

---

## 7. Security and Vulnerability Scanning

To maintain a secure software supply chain and catch security bugs early:
* **SAST Scanning (`gosec`)**: Code security checking is enabled natively as a linter within `.golangci.yml`. It runs automatically during `make lint` or CI actions, validating cryptographic usage, file permissions, and data validation boundaries.
* **Vulnerability Auditing (`govulncheck`)**: Go's official vulnerability tool is integrated inside the CI actions run and via `make vulncheck`. It uses call-path aware analysis to identify if your code imports any dependencies with known active CVE vulnerabilities.
```
