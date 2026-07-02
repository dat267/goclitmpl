# Codebase Rules & Guidelines for goclitmpl

This document contains rules and instructions for agentic AI coders working on this CLI template repository. Always adhere to these constraints to preserve the codebase as a high-quality template.

---

## 1. Codebase Integrity & Documentation Sync

* **Documentation Consistency**: Keep `README.md` and `DEVELOPER.md` fully in sync with the codebase structure, configuration file keys, environment variable prefixes, and subcommands.
* **Inline Comments**: All exported packages, structs, and methods must have standard Go documentation comments. Avoid redundant comments, but ensure complex code flows (like diagnostic net probes or custom logging handlers) are well-commented.
* **No Code Duplication**: Keep business logic separated from the command routing interface. Always place CLI interaction in [internal/cli](file:///home/dat/repos/goclitmpl/internal/cli) and core logic in reusable packages under [pkg](file:///home/dat/repos/goclitmpl/pkg) or [internal/config](file:///home/dat/repos/goclitmpl/internal/config).

---

## 2. CLI Command Structure & Output Design

* **Argument Validation**: Always validate user arguments within Cobra commands using explicit arg count limits (e.g., `Args: cobra.NoArgs` or `Args: cobra.ExactArgs(1)`).
* **Return Errors**: Implement Cobra subcommand runners using `RunE` instead of `Run` to propagate execution errors back to [cli.Execute()](file:///home/dat/repos/goclitmpl/internal/cli/root.go) for unified logging and clean exit code processing.
* **Separation of Streams**:
  * **Stdout**: Write raw commands output (such as JSON data or version text) using `cmd.OutOrStdout()` or `fmt.Fprintf(cmd.OutOrStdout(), ...)`.
  * **Stderr**: Reserve exclusively for structured logging (`slog`), diagnostic progress bars, warnings, or debug messages. Never mix output streams.

---

## 3. Template Flexibility & Renaming Support

* **No Hardcoded Namespaces**: When adding files containing imports or the template namespace `github.com/dat267/goclitmpl`, ensure that:
  * [rename.sh](file:///home/dat/repos/goclitmpl/rename.sh) is updated to replace the old namespace in the new file extension type.
  * [rename.ps1](file:///home/dat/repos/goclitmpl/rename.ps1) is updated to include the new file extension type.
* **Config AppName & EnvPrefix**: Ensure that application names, environment prefixes (e.g., `GOCLITMPL_`), and config home paths are derived dynamically from `config.AppName` and support the automated rename script overrides.

---

## 4. Strict Security & File Permissions

* **Directory Creation**: Any directory created by the application (e.g., for logs or configurations) must be restricted to the owner only (**`0o700`** permissions). Do not use `0o755`.
* **File Creation**: Any configuration or state file written by the CLI must be restricted to the owner only (**`0o600`** permissions). Do not use `0o644`.
* **Credential Protection**: Never hardcode credentials, secrets, or default keys. Ensure `gosec` static analysis remains enabled and runs cleanly.

---

## 5. Testing, Linting, & Isolation Rigor

* **Local Verification**: Before declaring any change complete, run the full validation suite locally:
  ```bash
  make all
  ```
  Or on Windows (PowerShell):
  ```powershell
  .\build.ps1
  ```
  This runs code formatting (`go fmt`), unit tests with the race detector (`go test -race`), security scanning (`govulncheck`), and standard linting.
* **Test Isolation**: Unit tests must not pollute the developer's home directory. Always isolate file creation tests using Go's built-in `t.TempDir()` or environment overrides via `t.Setenv()`.
* **Linter Warnings**: Zero linter warnings are permitted. Add inline `//nolint:...` overrides only for verified false positives (such as `slog.Record` copy-by-value in custom log handlers).
* **Test Coverage**: Every new subcommand or package must have a matching unit test file (e.g. `_test.go`) covering successful execution paths and error states.
