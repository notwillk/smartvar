# Project Instructions

> This is a template repository. The `build`, `clean`, `format`, and `test` recipes in the `justfile` are stubs — replace the placeholder `echo` commands with the appropriate commands for the project.

## Commands

Use [`just`](https://just.systems) for all common tasks. Run `just help` for a summary.

- `just build` / `just build --watch` — compile the project
- `just test` / `just test --watch` — run tests
- `just clean` / `just clean --deep` — remove build artifacts (`--deep` also removes generated files in `.schemas/`)
- `just format` — format code in-place
- `just format --check` — check formatting without modifying files
- `just doctor` / `just doctor --fix` — check environment health, auto-fix if possible
- `just static` / `just static --fix` — run static checks including format check; `--fix` formats and auto-fixes

## Checks

Checks are defined in YAML and run by `checksy`.

### Doctor checks (`doctor.checksy.yaml`)

Validate the dev environment (tools installed, env vars set, etc.). Run on container creation and with `just doctor`.

To add a check:
```yaml
rules:
  - name: description of check
    check: 'shell command that exits 0 on success'
    severity: error  # error | warning | debug
    fix: 'optional command to auto-fix'
```

### Static checks (`static.checksy.yaml`)

Validate the codebase (lint, format, etc.). Run in CI and with `just static`.

To add a check:
```yaml
rules:
  - name: description of check
    check: 'shell command that exits 0 on success'
    severity: error
    fix: 'optional auto-fix command'
```

## Directories and Files

- `.schemas/` — JSON Schema files for validating config/data. Generated schemas go here so editors and validators have a single known location to reference them. Cleared by `just clean --deep`.
- `.testignore` — file patterns excluded from triggering test re-runs in watch mode (`just test --watch`). Uses the same format as `.gitignore`. Add paths here for files whose changes should not re-run tests, such as documentation, editor config, and other non-code files.

## Devcontainer

Open in VS Code and choose **Reopen in Container**. The container runs `just doctor --fix` on creation, `just doctor` on each start, and includes `just`, `checksy`, and AI coding tools.
