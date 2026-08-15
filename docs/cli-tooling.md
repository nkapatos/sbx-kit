# Toolkit CLI (`sbx-kit`)

Helper CLI for Docker AI Sandboxes: **recipes**, **kits**, **portable state**,
and optional **local template builds**.

Two template paths (same lifecycle commands):

1. **Official / Hub** — omit local image fields in the recipe; `sbx` uses the
   stock agent template. Attach kits to experiment (Docker’s supported model).
2. **Local / registry tag** — set `image_name` / `template_fallback`, or run
   `sbx-kit template load` for images under `templates/`.

This repository ships a **few example recipes**. Point the toolkit root
(`SBX_TREE` / Homebrew share) at your own `config/` + `kits/` + `templates/`
when you want different stacks.

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
sbx-kit run --agent cursor-hub --yes      # Hub template + kits
sbx-kit run --agent cursor --yes          # local kit-cursor image
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
sbx-kit rm --agent cursor-hub --keep-state
sbx-kit run --agent cursor-hub --yes --restore-state   # only when the box does not exist yet
# or:
sbx-kit upgrade --agent cursor-hub                     # export → rm → create → restore → attach
```

`upgrade` recreates the sandbox from the **current recipe** (template + kits).
In-box agent binary refresh (distinct command) is still open.

## Catalog (recipes)

Recipes are declared in [`config/agents.yaml`](../config/agents.yaml). Each key
is a **user-chosen recipe id**; `sbx_agent` is the underlying sbx agent.

| Recipe fields | Meaning |
| --- | --- |
| No `image_name` / `template_fallback` | **Hub** — `sbx` uses the stock agent template |
| `image_name` + `template_fallback` | Resolve from `sbx template ls`, else fallback tag |
| `kits` | Mixin dirs under the toolkit `kits/` tree |

`sbx-kit agents` prints `RECIPE | SBX_AGENT | SOURCE | IMAGE | KITS | STATUS`.

Resource profiles: [`config/resources-remote-llm.env`](../config/resources-remote-llm.env), [`config/resources-local-llm.env`](../config/resources-local-llm.env).

## Develop from a checkout

```bash
cd cli
go build -ldflags "-X github.com/nkapatos/sbx-kit/cli/internal/version.Version=dev" \
  -o ../bin/sbx-kit ./cmd/sbx-kit
export SBX_TREE=/path/to/sbx-kit
../bin/sbx-kit agents
../bin/sbx-kit run --agent cursor-hub --yes
../bin/sbx-kit template load --engine docker kit-core
../bin/sbx-kit template load --engine docker kit-cursor
../bin/sbx-kit run --agent cursor --yes
../bin/sbx-kit init --agent cursor-hub /tmp/demo
```

Or: `go install github.com/nkapatos/sbx-kit/cli/cmd/sbx-kit@latest`
