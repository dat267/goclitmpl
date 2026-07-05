# Go CLI Template (`goclitmpl`)

[![CI](https://github.com/dat267/goclitmpl/actions/workflows/ci.yml/badge.svg)](https://github.com/dat267/goclitmpl/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/dat267/goclitmpl)](https://goreportcard.com/report/github.com/dat267/goclitmpl)

A highly optimized, structured, and production-ready Go CLI template project. Designed to serve as a robust foundation for building high-performance command line tools with modularity, structured logging, configuration overrides, and build size optimizations.

---

## Key Features

- **CLI Routing & Flag Parsing**: Powered by [Cobra](https://github.com/spf13/cobra) — the industry standard for CLI applications.
- **Dynamic Configuration**: Powered by [Viper](https://github.com/spf13/viper) — supports JSON/YAML/TOML config files, environment variables (`GOCLITMPL_*`), and CLI flag overrides.
- **Optimized Compilation**: Ready-to-go compiler flag configurations (`-ldflags="-s -w"`) to strip debug symbols and reduce binary size.
- **Structured Logging**: Built-in Go standard library `slog` with split-stream execution:
  - **Stdout**: Reserved exclusively for command output (perfect for pipelining, e.g., outputting to `jq`).
  - **Stderr**: Receives structured log statements (supports Text and JSON formats).
- **Graceful Shutdowns**: Global context listening for OS termination signals (`SIGINT`, `SIGTERM`) passed down to all subcommands.
- **DX Tooling**:
  - `Makefile` for compiling, testing, coverage analysis, linting, and cross-compiling.
  - `rename.sh` (Unix) and `rename.ps1` (Windows) helpers to instantly re-namespace the Go module path.
  - `.golangci.yml` lint rules preset.
  - GitHub Actions CI workflow setup.

---

## Directory Structure

```text
goclitmpl/
├── .github/workflows/  # CI/CD configurations
│   └── ci.yml
├── bin/                # Compiled binaries (generated on build)
├── cmd/
│   └── goclitmpl/
│       └── main.go     # Program entrypoint (handles signals & context)
├── internal/
│   ├── cli/            # Cobra commands (root, config, diagnose, greet, version)
│   │   ├── config.go   # Config CLI command (and its subcommands)
│   │   ├── diagnose.go # Diagnose CLI command (and its subcommands)
│   │   ├── greet.go    # Greet CLI command
│   │   ├── root.go     # Root CLI configuration & logging initialization
│   │   └── version.go  # Version CLI command
│   └── config/         # Viper configuration & slog logging initialization
│       ├── config.go
│       └── logging.go
├── Makefile            # Project tasks runner (build, test, lint, etc.)
├── pkg/                # Public/reusable library packages (decoupled logic)
│   ├── configinit/     # Business logic for config management
│   │   └── init.go
│   ├── diagnose/       # Business logic for system diagnostics
│   │   ├── check.go
│   │   └── info.go
│   ├── greet/          # Business logic for greeting messages
│   │   └── greet.go
│   └── version/        # Business logic for version formatting
│       └── version.go
├── build.ps1           # Script to build, test, and lint on Windows (PowerShell)
├── rename.ps1          # Script to quickly rename module imports on Windows (PowerShell)
├── rename.sh           # Script to quickly rename module imports on Unix (bash)
├── config.example.yaml # Example YAML configuration file
├── go.mod              # Go module descriptor
└── README.md           # Project documentation
```


---

## Getting Started

### 1. Rename the Module
Run the helper script to update module references from `github.com/dat267/goclitmpl` to your desired repository name:

On Linux/macOS:
```bash
./rename.sh github.com/your-username/your-cli-project
```

On Windows (PowerShell):
```powershell
.\rename.ps1 -NewModule github.com/your-username/your-cli-project
```

### 2. Run the Tests
On Linux/macOS:
```bash
make test
```

On Windows (PowerShell):
```powershell
.\build.ps1 -Task test
```

### 3. Build & Run
Compile the optimized binary into the `bin/` directory:

On Linux/macOS:
```bash
make build
```

On Windows (PowerShell):
```powershell
.\build.ps1 -Task build
```

Then run your command:

```bash
./bin/goclitmpl greet World
```

---

## Configuration & Environments

Configuration values are mapped to structs defined in [internal/config/config.go](file:///home/dat/repos/goclitmpl/internal/config/config.go).
Viper automatically loads configurations using the following precedence (highest to lowest):

1. **CLI Flags**: `--log-format json`
2. **Environment Variables**: Prefixed with `GOCLITMPL_` (e.g. `GOCLITMPL_APP_PORT=9090` overrides `app.port`).
3. **Configuration File**: Looks for a `config.yaml` file in:
   - The current working directory (`.`)
   - The user config directory (`~/.config/goclitmpl/config.yaml`)
   - System etc directory (`/etc/goclitmpl/config.yaml`)
   - Explicitly provided configuration file path via `--config <path>`
4. **Default Values**: Configured in `internal/config/config.go`.

### Example Environment Overrides
```bash
GOCLITMPL_APP_ENV=staging GOCLITMPL_LOG_LEVEL=debug ./bin/goclitmpl greet World
```

---

## Logging Guidelines

To maintain clean pipelines, structured log messages and actual data outputs are strictly separated:

* **Logs (Stderr)**: Use Go's standard library `slog` for warnings, errors, and debug messages:
  ```go
  import "log/slog"
  
  slog.Info("connecting to database", slog.String("host", host))
  ```
* **Data (Stdout)**: Use standard formatters writing directly to `cmd.OutOrStdout()` or the command's writer stream to print the command's output:
  ```go
  fmt.Fprintln(out, "Actual result data")
  ```

---

## Makefile Commands Reference

| Command | Description |
|---|---|
| `make build` | Builds optimized binary for the current OS/Arch (with stripped symbols and injected version vars). |
| `make run` | Runs the binary directly. Additional flags can be passed, e.g., `make run -- greet World`. |
| `make test` | Runs unit tests with race detection. |
| `make test-coverage` | Runs unit tests and opens an interactive HTML coverage report. |
| `make lint` | Validates code formatting and executes golangci-lint analysis. |
| `make fmt` | Automatically formats all Go code. |
| `make cross-compile` | Compiles binaries for Linux, macOS, and Windows. |
| `make clean` | Deletes build outputs and temporary coverage reports. |
