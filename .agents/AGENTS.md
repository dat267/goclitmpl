# Codebase Rules & Guidelines for goclitmpl

This document contains rules and instructions for agentic AI coders working on this CLI template repository. Always adhere to these constraints to preserve the codebase as a high-quality, minimal template.

---

## 1. Codebase Integrity & Documentation Sync

* **Documentation Consistency**: Keep `README.md` and `CONTRIBUTING.md` fully in sync with the codebase structure, configuration file keys, environment variable prefixes, and subcommands.
* **Inline Comments**: All exported packages, structs, and methods must have standard Go documentation comments. Avoid redundant comments, but ensure complex code flows (like diagnostic net probes or custom logging handlers) are well-commented.
* **No Code Duplication**: Keep business logic separated from the command routing interface. Always place CLI interaction in [internal/cli](file:///home/dat/repos/goclitmpl/internal/cli) and core logic in reusable packages under [pkg](file:///home/dat/repos/goclitmpl/pkg) or [internal/config](file:///home/dat/repos/goclitmpl/internal/config).

---

## 2. CLI Command Structure & Output Design

* **Argument Validation**: Always validate user arguments within Cobra commands using explicit arg count limits (e.g., `Args: cobra.NoArgs` or `Args: cobra.ExactArgs(1)`).
* **Return Errors**: Implement Cobra subcommand runners using `RunE` instead of `Run` to propagate execution errors back to [cli.Execute()](file:///home/dat/repos/goclitmpl/internal/cli/root.go) for unified logging and clean exit code processing.
* **Separation of Streams**:
  * **Stdout**: Write raw command output (such as JSON data or version text) using `cmd.OutOrStdout()` or `fmt.Fprintf(cmd.OutOrStdout(), ...)`.
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

## 5. Linting & Build Verification

* **Local Verification**: Before declaring any change complete, run the full validation suite locally:
  ```bash
  make all
  ```
  Or on Windows (PowerShell):
  ```powershell
  .\build.ps1
  ```
  This runs code formatting (`go fmt`), security scanning (`govulncheck`), and standard linting.
* **No Test Files**: This is a clean, minimal template. Test files (`_test.go`) are intentionally omitted, and no tests are needed. Do **not** add `_test.go` files to the repository.
* **Linter Warnings**: Zero linter warnings are permitted. Add inline `//nolint:...` overrides only for verified false positives (such as `slog.Record` copy-by-value in custom log handlers).

---

## 6. Nested Subcommand Best Practices

* **Group commands must not define `RunE`**: A parent command that only groups subcommands (e.g., `config`, `diagnose`) must **not** define a `Run` or `RunE` function. Cobra automatically shows the help page when such a command is invoked with no subcommand. Defining `RunE: func(...) { return cmd.Help() }` is redundant and should be removed.
* **`SilenceErrors: true` on root**: Set `SilenceErrors: true` on the root command and print the error explicitly in `Execute()`. This prevents Cobra from double-printing errors when both the command and the caller handle them.
* **`SilenceUsage` after validation only**: Set `cmd.SilenceUsage = true` inside `PersistentPreRunE` (not at construction time) so that flag-parsing errors still display usage, but runtime errors do not.
* **`Example` fields on every command**: Every command (parent and leaf) must define an `Example` field with at least one concrete invocation string to aid discoverability via `--help`.
* **Required vs optional arg notation in `Use`**: Use angle brackets for required positional arguments (`<name>`) and square brackets for optional ones (`[address]`). Example: `Use: "greet <name>"`.
* **`PersistentFlags` for shared flags**: Flags consumed by multiple subcommands (e.g., `--timeout` on `diagnose`) must be defined via `cmd.PersistentFlags()` on the parent, not repeated per-subcommand. Access them from leaf commands via `cmd.Flags().GetXxx("flag-name")`.
* **Subcommand file layout**: Each CLI command lives in `internal/cli/<cmd>.go`. All subcommands of `<cmd>` are defined in that same file using unexported constructor functions (e.g., `newDiagnoseInfoCmd()`). Business logic lives in the corresponding `pkg/<cmd>/` package.

