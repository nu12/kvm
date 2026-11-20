# Agent guidance

## Project overview

KVM is a Go desktop application built with Fyne for browsing Azure
Key Vaults, secrets, certificates, and keys. Azure SDK client construction is
configured in `main.go`.

## Repository layout

- `gui/windows/`: application root, splash screen, and vault-window lifecycle.
- `gui/tabs/`: per-vault secrets, certificates, and keys tabs.
- `gui/components/`: reusable Fyne cells, tables, and widgets.
- `gui/dialogs/`: dialogs and their interaction logic.
- `gui/notifications/`: activity indicators and error handling.
- `services/`: Azure resource operations and pagination.
- `factory/`: Azure client interfaces, factories, and test mocks.

## Implementation conventions

- Keep changes focused and follow existing Go formatting and package patterns.
- Keep Azure operations in `services/`; resolve clients through `factory/`.
  Extend the relevant client interface and mock when adding an Azure operation.
- Keep UI behavior in its relevant window, tab, component, or dialog. Do not put
  dialog/service logic in table-cell widgets.
- UI code should use the service and factory layers rather than creating Azure 
  clients directly
- Run blocking network or file I/O outside the Fyne UI thread. When background
  work completes, wrap UI mutations (including widget refresh/text changes,
  window operations, and dialogs/error notifications) in `fyne.Do`. Keep only
  UI work inside that callback. `notifications.Activity.Start` and `Stop`
  already schedule their own widget updates.
- Surface operation failures using the existing notification/error helpers or a
  Fyne error dialog. Do not silently ignore errors or use `log.Fatal` for
  recoverable UI actions.
- Searchable lists should retain the full collection separately from the visible
  filtered collection, filter from the entry's `OnChanged` callback, and refresh
  the table after updating visible rows. An empty query should restore all rows.
- Avoid adding UI elements or behavior outside the request (for example, do not
  add a count label to the vault secrets tab unless explicitly requested).
- Use the existing factory mocks and Fyne test driver for focused tests. Tests
  that enqueue Fyne UI work may need to wait for the UI callback to complete.

## Validation

- Format modified Go files with `gofmt`.
- Run the narrowest relevant tests while iterating, for example:
  `go test ./gui/tabs`, `go test ./gui/windows`, or `go test ./services`.
- Before completing a change, run `go test ./...` when practical. CI also runs
  Go lint and tests on pull requests.
- The Fyne desktop dependencies may require platform libraries on Linux; CI
  installs Wayland/X11 development packages before linting and testing.

## Safety and scope

- Never place credentials, tokens, or secret values in source, logs, or tests.
- Preserve unrelated worktree changes. Inspect the diff before editing files
  that already contain user changes, and do not revert unrelated modifications.
