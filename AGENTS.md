# yq — agent instructions

See [agents.md](./agents.md) for contributor workflow (formatting, testing patterns, adding encoders/decoders).

## Cursor Cloud specific instructions

### Overview

**yq** is a Go CLI for querying and transforming YAML, JSON, XML, INI, and other structured formats. There are no long-running services — development is build-and-test against a local `./yq` binary.

### Prerequisites

- **Go ≥ 1.25** (see `go.mod`)
- **Node.js** (for `npx cspell` spelling checks in the full `make test` pipeline)
- **Bash** (acceptance tests)
- **Docker/Podman** is optional; use `make local <target>` to run natively when containers are unavailable

### PATH

After `scripts/devtools.sh`, add Go tool binaries to PATH:

```bash
export PATH="$HOME/go/bin:$PATH"
```

`golangci-lint` installs to `$HOME/go/bin`; `gosec` installs to `./bin/gosec` in the repo root.

### Common commands (local, no Docker)

| Task | Command |
|------|---------|
| Install dev tools | `bash scripts/devtools.sh` |
| Vendor dependencies | `make local vendor` |
| Build binary | `go build -o yq .` or `make local build` |
| Format | `make local format` |
| Lint | `make local check` |
| Unit tests | `make local test` or `bash scripts/test.sh` |
| Acceptance (E2E) | `bash scripts/acceptance.sh` (requires `./yq` built first) |

`make local build` runs the full CI chain (format → spelling → gosec → lint → unit tests → build → acceptance). For a faster loop, build with `go build -o yq .` and run `bash scripts/acceptance.sh`.

### Gotchas

- **`make` without `local`** tries Docker/Podman (`Dockerfile.dev`). In Cloud Agent VMs without Docker, always prefix with `make local`.
- **Spelling step** uses `npx cspell` and may download cspell on first run (network required).
- **`make local test` / `scripts/check.sh`** require `golangci-lint` on PATH (`devtools.sh`).
- As of the current tree, `pkg/yqlib/ini_test.go` calls `NewINIDecoder()` without required `INIPreferences`, which breaks unit-test compilation in `make local check` / `make local test`. The **binary still builds** (`go build -o yq .`) and **all acceptance tests pass** (`bash scripts/acceptance.sh`).
