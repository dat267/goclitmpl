# Codebase Rules & Guidelines for goclitmpl

This document contains rules and instructions for agentic AI coders working on this CLI template repository. Always adhere to these constraints to preserve the codebase as a high-quality template.

---

## 1. Codebase Integrity & Documentation Sync

* **Documentation Consistency**: Keep `README.md` and `DEVELOPER.md` fully in sync with the codebase structure, configuration file keys, environment variable prefixes, and subcommands.
* **Inline Comments**: All exported packages, structs, and methods must have standard Go documentation comments. Avoid redundant comments, but ensure complex code flows (like diagnostic net probes or custom logging handlers) are well-commented.
* **No Code Duplication**: Keep business logic separated from the command routing interface. Always place CLI interaction in `internal/cli/` and core logic in reusable packages under `pkg/` or `internal/config/`.

---

## 2. Template Flexibility & Renaming Support

* **No Hardcoded Namespaces**: When adding files containing imports or the template namespace `github.com/dat267/goclitmpl`, ensure that:
  * [rename.sh](file:///home/dat/repos/goclitmpl/rename.sh) is updated to replace the old namespace in the new file extension type.
  * [rename.ps1](file:///home/dat/repos/goclitmpl/rename.ps1) is updated to include the new file extension type.
* **Config AppName & EnvPrefix**: Ensure that application names, environment prefixes (e.g., `GOCLITMPL_`), and config home paths are derived dynamically from `config.AppName` and support the automated rename script overrides.

---

## 3. Strict Security & File Permissions

* **Directory Creation**: Any directory created by the application (e.g. for logs or configurations) must be restricted to the owner only (**`0o700`** permissions). Do not use `0o755`.
* **File Creation**: Any configuration or state file written by the CLI must be restricted to the owner only (**`0o600`** permissions). Do not use `0o644`.
* **Credential Protection**: Never hardcode credentials, secrets, or default keys. Ensure `gosec` static analysis remains enabled and runs cleanly.

---

## 4. Linting & Validation Rigor

* **Local Verification**: Before declaring any change complete, run the full validation suite locally:
  ```bash
  make all
  ```
  This runs code formatting (`go fmt`), unit tests with the race detector (`go test -race`), security scanning (`govulncheck`), and standard linting.
* **Linter Warnings**: Zero linter warnings are permitted. Add inline `//nolint:...` overrides only for verified false positives (such as `slog.Record` copy-by-value in custom log handlers).
* **Test Coverage**: Every new subcommand or package must have a matching unit test file (e.g. `_test.go`) covering successful execution paths and error states.
