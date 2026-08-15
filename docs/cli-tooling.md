# Toolkit CLI (`sbx-kit`)

Convenience layer on Docker AI Sandboxes (`sbx`): **recipes**, kit placement,
portable state, and custom images (build/load locally or pull a registry tag,
then import into sbx).

Same words as sbx for **kind** (`sbx run <kind>`), **template** (engine image
store), and **kit**. **Recipe** is the only sbx-kit-specific idea (named kind +
kits + optional custom image). Run `sbx-kit concepts` for a short wiring guide.

There is no `sbx agents` command and no `--agent` flag. Available kinds are in
`sbx run --help`. `sbx template ls` is the engine import store, not our catalog.

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
sbx-kit run                                  # attach sole cwd binding
sbx-kit run <recipe> [--path <dir>] [--sandbox-name <name>] [--yes] …
sbx-kit run --name <sandbox>                 # attach only
sbx-kit rm --recipe <id>|--name … [--keep-state] [--force]
sbx-kit upgrade --recipe <id>|--name …
sbx-kit state export|import --recipe|--name …
sbx-kit status [--path <dir>]
sbx-kit check [--name|--recipe/--path]
sbx-kit init [--recipe <id>] [project-dir]
sbx-kit image ls
sbx-kit image load --engine <docker|container> <name> [tag]
sbx-kit image pull [--engine docker] <registry/tag>
sbx-kit version
```

`--recipe` remains on lifecycle commands when a positional does not apply
(`rm`, `upgrade`, `check`, `init`). Prefer `sbx-kit run cursor` over
`sbx-kit run --recipe cursor`.

## Host vault (XDG)

| Path | Contents |
| --- | --- |
| `~/.local/share/sbx-kit/profiles/<id>/state.tgz` | Portable VM state archives |
| `~/.local/state/sbx-kit/bindings.json` | project + recipe → sandbox name |

## Sandbox identity

| Flag | Intent |
| --- | --- |
| positional recipe / `--recipe <id>` | **Create** from catalog (or resolve binding for rm/upgrade/…) |
| `--path` | Project directory |
| `--sandbox-name` | Name at **create** (default: dirname) |
| `--name` | **Existing** sandbox (attach / rm / check / state) |

```bash
sbx-kit recipes
sbx-kit run shell --yes
sbx-kit check
sbx-kit run --name my-project
```

## Host secrets

On **create**, `run` prints `sbx secret set <service>` for services declared in
the recipe’s kits. Set **any** you use; they are not all required. Extra APIs
belong in that kit (or a personal mixin), not a kit per key. `check` runs
`sbx secret ls` (sandbox-scoped when the box exists).

## Portable state

```bash
sbx-kit rm --recipe cursor --keep-state
sbx-kit run cursor --yes --restore-state
sbx-kit upgrade --recipe cursor
```

## Catalog

Recipes live in [`config/agents.yaml`](../config/agents.yaml) (filename kept for
now). Stock ids match sbx kinds (`shell`, `cursor`) or a sandbox-kit name (`pi`);
custom images use a `kit-` prefix. Catalog default kit is `agent-workspace`.
`sbx-kit recipes` prints `RECIPE | SBX_AGENT | SOURCE | IMAGE | KITS`.

## Custom images

`sbx-kit image ls` lists Dockerfiles under `templates/`. `load` builds one;
`pull` fetches a registry tag. Both import via `sbx template load` so the image
shows up in `sbx template ls`.

## Develop from a checkout

```bash
cd cli
go build -ldflags "-X github.com/nkapatos/sbx-kit/cli/internal/version.Version=dev" \
  -o ../bin/sbx-kit ./cmd/sbx-kit
export SBX_TREE=/path/to/sbx-kit
../bin/sbx-kit concepts
../bin/sbx-kit recipes
../bin/sbx-kit run shell --yes
```
