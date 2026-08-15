# Toolkit CLI (`sbx-kit`)

Convenience layer on Docker AI Sandboxes (`sbx`): **recipes**, kit placement,
portable state, and optional local template builds.

Same words as sbx for **agent**, **template**, and **kit**. **Recipe** is the
only sbx-kit-specific idea (named agent + kits shortcut). Run `sbx-kit concepts`
for a short wiring guide.

**Install (macOS):** [homebrew.md](homebrew.md) — `brew tap nkapatos/sbx-kit …`  
**Go module:** `github.com/nkapatos/sbx-kit/cli`

## Required `sbx` CLI

Requires **Docker `sbx` >= 0.34.0**. Kits are authored as schemaVersion `"1"`.

```bash
sbx-kit version
sbx kit validate kits/agent-workspace
```

Escape hatch: `SBX_KIT_SKIP_SBX_CHECK=1`.

## Commands

```text
sbx-kit concepts
sbx-kit recipes
sbx-kit agents                               # sbx agents + custom templates in view
sbx-kit template ls                          # → sbx template ls
sbx-kit run                                  # attach sole cwd binding
sbx-kit run --recipe <id> [--path <dir>] [--sandbox-name <name>] [--yes] …
sbx-kit run --name <sandbox>                 # attach only
sbx-kit rm --recipe <id>|--name … [--keep-state] [--force]
sbx-kit upgrade --recipe <id>|--name … [--force]
sbx-kit state export|import --recipe|--name …
sbx-kit status [--path <dir>]
sbx-kit check [--name|--recipe/--path]
sbx-kit init [--recipe <id>] [project-dir]
sbx-kit template load --engine <docker|container> <name> [tag]
sbx-kit version
```

`--agent` is accepted as an alias for `--recipe` on lifecycle commands.

## Host vault (XDG)

| Path | Contents |
| --- | --- |
| `~/.local/share/sbx-kit/profiles/<id>/state.tgz` | Portable VM state archives |
| `~/.local/state/sbx-kit/bindings.json` | project + recipe → sandbox name |

## Sandbox identity

| Flag | Intent |
| --- | --- |
| `--recipe <id>` | **Create** from catalog (or resolve binding for rm/upgrade/…) |
| `--path` | Project directory |
| `--sandbox-name` | Name at **create** (default: dirname) |
| `--name` | **Existing** sandbox (attach / rm / check / state) |

```bash
sbx-kit recipes
sbx-kit run --recipe shell-hub --yes
sbx-kit check
sbx-kit run --name my-project
```

## Host secrets

On **create**, `run` prints `sbx secret set <service>` for services declared in
the recipe’s kits. `check` runs `sbx secret ls` (sandbox-scoped when the box
exists).

## Portable state

```bash
sbx-kit rm --recipe cursor-hub --keep-state
sbx-kit run --recipe cursor-hub --yes --restore-state
sbx-kit upgrade --recipe cursor-hub
```

## Catalog

Recipes live in [`config/agents.yaml`](../config/agents.yaml) (filename kept for
now). `sbx-kit recipes` prints `RECIPE | SBX_AGENT | SOURCE | IMAGE | KITS`.

## Develop from a checkout

```bash
cd cli
go build -ldflags "-X github.com/nkapatos/sbx-kit/cli/internal/version.Version=dev" \
  -o ../bin/sbx-kit ./cmd/sbx-kit
export SBX_TREE=/path/to/sbx-kit
../bin/sbx-kit concepts
../bin/sbx-kit recipes
../bin/sbx-kit run --recipe shell-hub --yes
```
