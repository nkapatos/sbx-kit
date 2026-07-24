# Toolkit CLI (`sbx-kit`)

Catalog-driven Go CLI for Docker AI Sandboxes templates, kits, and resource profiles.

**Install (macOS):** [homebrew.md](homebrew.md) — `brew tap nkapatos/sbx-kit …`  
**Go module:** `github.com/nkapatos/sbx-kit/cli`

sbx-kit is a convenience layer for **custom templates/kits + lifecycle** (stable identity, state export/restore, upgrade). Default Hub agents work fine with plain `sbx`; this CLI targets local (and later registry) recipes without re-teaching sbx’s agent/`--name` dual forms.

## Commands

```text
sbx-kit agents
sbx-kit run --agent <name> [--path <dir>] [--yes] [--clone] [--restore-state] [-- sbx-args...]
sbx-kit run --name <sandbox> [--restore-state]
sbx-kit rm --agent <name> [--path <dir>] [--keep-state] [--force]
sbx-kit rm --name <sandbox> [--keep-state] [--force]
sbx-kit upgrade --agent <name> [--path <dir>] [--force]
sbx-kit state export|import --agent <name> [--path <dir>]
sbx-kit state export|import --name <sandbox>
sbx-kit status [--path <dir>]
sbx-kit init [--agent <name>] [project-dir]
sbx-kit template load --engine <docker|container> <name-or-path> [image-tag]
sbx-kit version
```

## Host vault (XDG)

Created lazily by lifecycle commands (not by Homebrew):

| Path | Contents |
| --- | --- |
| `~/.local/share/sbx-kit/profiles/<id>/state.tgz` | Portable VM state archives |
| `~/.local/state/sbx-kit/bindings.json` | project + agent → sandbox name |

Honor `XDG_DATA_HOME` / `XDG_STATE_HOME` when set.

## Sandbox identity

Each project+agent binding gets a stable opaque `sbx --name` (`sbxk-<agent>-<hash>`) so two checkouts with the same directory basename do not collide. Humans see a **label** (project basename) via `sbx-kit status`; you rarely type the id.

| Flag | Meaning |
| --- | --- |
| `--agent` | Catalog **recipe** (template + kits) for create / project-scoped attach |
| `--path` | Project directory (default `.`) |
| `--name` | Attach/rm/state by sandbox id (**no create**) |

`sbx-kit run --agent cursor` attaches if the binding’s sandbox exists; otherwise it **prompts** before create (`--yes` skips the prompt).  
`sbx-kit run --name …` only attaches; missing sandboxes error with a pointer to `status` / `sbx ls`.

## Portable state

Pack/unpack is **not** hardcoded in Go. The host CLI runs:

1. Best-effort wait if the sandbox is still `running` (detach first — agents often keep chat history in SQLite WAL files)
2. `sbx exec <name> -- sbx-kit-state pack|unpack /tmp/sbx-kit-state.tgz` (pack checkpoints `*.db` WALs when `sqlite3` is available; otherwise includes `-wal`/`-shm`)
3. `sbx cp` between that path and the host profile archive

`sbx-kit-state` and `state.manifest` ship in the **`agent-workspace` kit**. Manifest INCLUDE/EXCLUDE lists what survives recreate; caches are excluded.

```bash
sbx-kit rm --agent cursor --keep-state    # export then destroy
sbx-kit run --agent cursor --restore-state
# or:
sbx-kit upgrade --agent cursor            # export → rm → create → restore → attach
```

## Catalog

Agents are declared in [`config/agents.yaml`](../config/agents.yaml). Defaults pull resource profile + kits (`mise-workspace`, `agent-workspace`); each agent sets the sbx agent name, image name, and kit list.

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
