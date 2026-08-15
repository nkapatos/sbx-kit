# Product scope (sbx-kit example tree)

Guardrail for this repo’s templates, kits, and catalog. CLI can compose any
tree via `SBX_TREE` — that does not expand what we ship as examples.

## Goals

1. **Recipes + kits for everyone** — novices and power users; works on **official
   Hub templates** (Docker’s supported customization path) and on **local /
   remote** custom images.
2. **Clear template ownership** — Hub vs local/registry build; CLI makes where
   images come from obvious.
3. **Portable agent state** — survive recreate, host moves, and recipe changes
   (`upgrade` / export / import).
4. **Own custom images when you want them** — build/load for development;
   publish to a registry and pin tags in recipes; rare rebuilds; mise for langs;
   kits for preference CLIs; host-side agent refresh (not mid-session).
5. Keep the brew/share tree small and teachable.

## Two paths

| Path | How |
| --- | --- |
| **Shell + kits** | Official Hub `shell` or our `kit-shell`. Add stuff with **kits** (Pi, mixins). Do not bake a new image. |
| **Image on core** | `FROM kit-core` for custom images you own (e.g. `kit-cursor`). **Do not** import kit-core into sbx. |

## Layering rule (local / custom images)

| Layer | Owns | Examples |
| --- | --- | --- |
| **kit-core** | Parent image: OS, utils, sbx glue, mise **binary**. Never `sbx template load`. | `FROM` for `kit-shell`, `kit-cursor`; later Docker-on-VPS |
| **kit-shell** | Hub-shell counterpart: tini + login bash only | Kits on top (`kit-pi`), same as official `shell` |
| **Agent image** | Bootstrap install/layout only | `kit-cursor` **FROM kit-core**, not via kit-shell |
| **Mixin kit** | Activate, allowlists, state, preference CLIs | `mise-workspace`, `agent-workspace`, `apt-extras` |
| **Sandbox kit** | Agent Hub doesn’t ship: image + entrypoint + install + provider secrets | `kits/pi` (recipes `pi` / `kit-pi`) |
| **Project** | Language pins | `mise.toml` |
| **Host / CLI** | Recipes, kit placement, state; later agent refresh | |

## In scope

- Kits + recipes; Hub-first create path; local `image load` / `image pull` when needed
- `templates/kit-core` (FROM parent, not imported), `templates/kit-shell` (minimum load), `templates/kit-cursor`
- Mixins: mise-workspace, agent-workspace (default on every recipe); optional lsp-mise, apt-extras
- Sandbox kit example: `kits/pi` (install + provider secrets in the same kit) on official shell (`pi`) and kit-shell (`kit-pi`)
- CLI lifecycle + vault; `image ls` of our Dockerfiles; `check` over sbx

## Out of scope / parked

- Docker Compose / VPS twin under `deploy/` (archived; same **kit-core** parent when it returns)
- Official-base + apt purge; thin “creds-only” Pi mixins as Hub workarounds
- Baking `gh` alone (or any preference CLI) into core “for parity”
- Daily image rebuilds for agent CLI churn
- Hot-swapping the agent binary mid-session

## Compatibility

- Kits: schemaVersion `"1"`
- CLI: `sbx` >= floor in `cli/internal/sbxcompat`
