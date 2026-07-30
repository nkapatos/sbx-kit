# Toolkit CLI (`sbx-kit`)

Helper CLI for **custom** Docker AI Sandboxes templates/kits + lifecycle
(stable identity, state export/restore, upgrade).

This repository ships a **few example recipes** so the shape is obvious. It is
not a large first-party agent catalog — point the toolkit root (`SBX_TREE` /
Homebrew share today; more trees later) at your own or someone else’s
`config/` + `kits/` + `templates/` when you want different stacks. Default Hub
agents work fine with plain `sbx`.

**Install (macOS):** [homebrew.md](homebrew.md) — `brew tap nkapatos/sbx-kit …`  
**Go module:** `github.com/nkapatos/sbx-kit/cli`

## Required `sbx` CLI

Kits/templates are still experimental upstream. This tree requires **Docker
`sbx` >= 0.34.0** and authors kits as **schemaVersion `"1"`** (released CLIs
through at least 0.37 do not yet accept clean authored v2). Lifecycle commands
refuse an older `sbx`.

```bash
sbx-kit version          # prints sbx-kit + required range + detected sbx
brew upgrade docker/tap/sbx
sbx kit validate kits/pi # expect VALID; deprecation WARNs are OK
```

Escape hatch (not recommended): `SBX_KIT_SKIP_SBX_CHECK=1`.
Bump the floor in `cli/internal/sbxcompat` when we depend on newer sbx.

## Commands

```text
sbx-kit agents
sbx-kit run                                      # attach sole binding for cwd
sbx-kit run --agent <recipe> [--path <dir>] [--sandbox-name <name>] [--yes] [--clone] [--restore-state] [-- sbx-args...]
sbx-kit run --name <sandbox> [--restore-state]   # attach only
sbx-kit rm --agent <recipe> [--path <dir>] [--keep-state] [--force]
sbx-kit rm --name <sandbox> [--keep-state] [--force]
sbx-kit upgrade --agent <recipe> [--path <dir>] [--force]
sbx-kit state export|import --agent <recipe> [--path <dir>]
sbx-kit state export|import --name <sandbox>
sbx-kit status [--path <dir>]
sbx-kit init [--agent <recipe>] [project-dir]
sbx-kit template load --engine <docker|container> <name-or-path> [image-tag]
sbx-kit version
```

## Host vault (XDG)

Created lazily by lifecycle commands (not by Homebrew):

| Path | Contents |
| --- | --- |
| `~/.local/share/sbx-kit/profiles/<id>/state.tgz` | Portable VM state archives |
| `~/.local/state/sbx-kit/bindings.json` | project + recipe → sandbox name |

Honor `XDG_DATA_HOME` / `XDG_STATE_HOME` when set.

## Sandbox identity

Friendly **sbx name** (what `sbx ls` shows) defaults to the project directory
basename — same idea as stock sbx. An opaque **profile id**
(`sbxk-<recipe>-<hash>`) keys the host vault only; users rarely type it.

| Flag | Intent |
| --- | --- |
| `--agent <recipe>` | **Create** from catalog recipe |
| `--path` | Project directory for create / bare-run attach |
| `--sandbox-name` | Create-time friendly sbx name (default: dirname; `--yes` skips prompt) |
| `--name` | **Attach** (or rm/state) by friendly sbx name — no create |

```bash
sbx-kit run --agent cursor --yes          # create; name = dirname
sbx-kit run --agent cursor                # interactive: name? then create?
sbx-kit run                               # re-attach sole cwd binding
sbx-kit run --name my-project             # attach from anywhere
```

If `--agent` is used and that sbx name already exists, the CLI errors (sbx
owns uniqueness). To “rename”: `rm --keep-state`, recreate with a new
`--sandbox-name`, `--restore-state`.

## Portable state

Pack/unpack is **not** hardcoded in Go. The host CLI runs:

1. Best-effort wait if the sandbox is still `running` (detach first — agents often keep chat history in SQLite WAL files)
2. `sbx exec <name> -- sbx-kit-state pack|unpack /tmp/sbx-kit-state.tgz` (pack checkpoints `*.db` WALs when `sqlite3` is available; otherwise includes `-wal`/`-shm`)
3. `sbx cp` between that path and the host profile archive

`sbx-kit-state` and `state.manifest` ship in the **`agent-workspace` kit**. Manifest INCLUDE/EXCLUDE lists what survives recreate; caches are excluded.

```bash
sbx-kit rm --agent cursor --keep-state
sbx-kit run --agent cursor --yes --restore-state   # only when the box does not exist yet
# or:
sbx-kit upgrade --agent cursor                     # export → rm → create → restore → attach
```

## Catalog (recipes)

Recipes are declared in [`config/agents.yaml`](../config/agents.yaml). Each key
is a **user-chosen recipe id**; `sbx_agent` is the underlying sbx agent. Example
entries (`cursor`, `opencode`, …) are starting points — rename or replace them
in your own tree.

`sbx-kit agents` prints `RECIPE | SBX_AGENT | IMAGE | KITS | STATUS`.

Resource profiles: [`config/resources-remote-llm.env`](../config/resources-remote-llm.env), [`config/resources-local-llm.env`](../config/resources-local-llm.env).

## Develop from a checkout

```bash
cd cli
go build -ldflags "-X github.com/nkapatos/sbx-kit/cli/internal/version.Version=dev" \
  -o ../bin/sbx-kit ./cmd/sbx-kit
export SBX_TREE=/path/to/sbx-kit
../bin/sbx-kit agents
../bin/sbx-kit run --agent cursor --yes
../bin/sbx-kit template load --engine docker cursor-mise-docker
../bin/sbx-kit init --agent cursor /tmp/demo
```

Or: `go install github.com/nkapatos/sbx-kit/cli/cmd/sbx-kit@latest`
