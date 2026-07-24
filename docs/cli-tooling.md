# Toolkit CLI (`sbx-kit`)

Catalog-driven Go CLI for Docker AI Sandboxes templates, kits, and resource profiles.

**Install (macOS):** [homebrew.md](homebrew.md) — `brew tap nkapatos/sbx-kit …`  
**Go module:** `github.com/nkapatos/sbx-kit/cli`

## Commands

```text
sbx-kit agents
sbx-kit run <agent> [project-dir] [--restore-state] [--clone] [-- sbx-args...]
sbx-kit rm <agent|name> [project-dir] [--keep-state] [--force]
sbx-kit upgrade <agent> [project-dir] [--force]
sbx-kit state export|import <agent|name> [project-dir]
sbx-kit status [project-dir]
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

`sbx-kit run` always injects a stable `sbx --name` (`sbxk-<agent>-<hash>`) unless you pass `--name` after `--`. The binding is written before attach so `rm` / `state` / `upgrade` can resolve it. Check with `sbx-kit status`.

## Portable state

Pack/unpack is **not** hardcoded in Go. The host CLI runs:

1. `sbx exec <name> -- sbx-kit-state pack|unpack /tmp/sbx-kit-state.tgz`
2. `sbx cp` between that path and the host profile archive

`sbx-kit-state` and `state.manifest` ship in the **`agent-workspace` kit** (installed at startup onto `/usr/local/bin`). Manifest INCLUDE/EXCLUDE lists what survives recreate; caches are excluded.

```bash
sbx-kit rm cursor . --keep-state    # export then destroy
sbx-kit run cursor . --restore-state
# or:
sbx-kit upgrade cursor .            # export → rm → create → restore → attach
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
../bin/sbx-kit run cursor .
../bin/sbx-kit template load --engine docker cursor-mise-docker
../bin/sbx-kit init --agent cursor /tmp/demo
```

Or: `go install github.com/nkapatos/sbx-kit/cli/cmd/sbx-kit@latest`
