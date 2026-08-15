# Product scope (sbx-kit example tree)

Guardrail for this repo’s templates, kits, and catalog. CLI can compose any
tree via `SBX_TREE` — that does not expand what we ship as examples.

## Goals

1. **Recipes + kits for everyone** — novices and power users; works on **official
   Hub templates** (Docker’s supported customization path) and on **local**
   images you build.
2. **Clear template ownership** — Hub/registry vs local build; CLI makes where
   images come from obvious.
3. **Portable agent state** — survive recreate, host moves, and recipe changes
   (`upgrade` / export / import).
4. **Own lean images when you want them** — rare rebuilds; mise for langs; kits
   for preference CLIs; host-side agent refresh (not mid-session).
5. Keep the brew/share tree small and teachable.

## Layering rule (local images)

| Layer | Owns | Examples |
| --- | --- | --- |
| **kit-core** | OS essentials, modern utils, sbx glue, mise **binary** | `templates/kit-core` |
| **Agent image** | Bootstrap install/layout only | `kit-cursor` (next `kit-pi`) |
| **Mixin kit** | Activate, allowlists, state, preference CLIs, creds | `mise-workspace`, `agent-workspace`, `apt-extras`, `pi` |
| **Project** | Language pins | `mise.toml` |
| **Host / CLI** | Recipes, kit placement, state; later agent refresh | |

On **Hub** templates, kits alone adjust the box — no local Dockerfile required.
Languages → mise (when the image has mise). Preference CLIs → kits. Agent
version churn on local images → host refresh, not daily rebake.

## In scope

- Kits + recipes; Hub-first create path; local `template load` when needed
- `templates/kit-core`, `templates/kit-cursor` as optional lean floor examples
- Mixins: mise-workspace, agent-workspace; optional lsp-mise, apt-extras; pi creds
- CLI lifecycle + vault

## Out of scope / parked

- Docker Compose / VPS twin under `deploy/` (archived out of tree; revisit later)
- Official-base + apt purge (abandoned)
- Baking `gh` alone (or any preference CLI) into core “for parity”
- Daily Hub rebuilds for Cursor/Pi releases
- Hot-swapping the agent binary mid-session
- Hermès until core + cursor + pi are boring

## Compatibility

- Kits: schemaVersion `"1"`
- CLI: `sbx` >= floor in `cli/internal/sbxcompat`
