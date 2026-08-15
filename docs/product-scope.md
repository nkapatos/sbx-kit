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
4. **Own lean images when you want them** — rare rebuilds; mise for langs; kits
   for preference CLIs; host-side agent refresh (not mid-session).
5. Keep the brew/share tree small and teachable.

## Two paths

| Path | How |
| --- | --- |
| **Hub + kits** | Stock `sbx` agent/template; recipes attach kits (incl. sandbox kits that start agents Hub doesn’t ship). Secrets via `sbx secret` / proxyManaged. |
| **Custom floor** | `kit-core` (+ optional baked agent layers like `kit-cursor`). Same pattern later for Pi/Hermès/… — pull local or remote; not Hub workarounds. |

## Layering rule (local / custom images)

| Layer | Owns | Examples |
| --- | --- | --- |
| **kit-core** | OS essentials, modern utils, sbx glue, mise **binary** | `templates/kit-core` |
| **Agent image** | Bootstrap install/layout only | `kit-cursor`; later other agents on core |
| **Mixin kit** | Activate, allowlists, state, preference CLIs | `mise-workspace`, `agent-workspace`, `apt-extras` |
| **Sandbox kit** | Official-path agent entry (startup on shell, etc.) | Community / user kits; not shipped here yet |
| **Project** | Language pins | `mise.toml` |
| **Host / CLI** | Recipes, kit placement, state; later agent refresh | |

## In scope

- Kits + recipes; Hub-first create path; local `template load` when needed
- `templates/kit-core`, `templates/kit-cursor` as optional lean floor examples
- Mixins: mise-workspace, agent-workspace; optional lsp-mise, apt-extras;
  deepseek-creds (Hub secret trial)
- CLI lifecycle + vault; thin `template ls` / `check` over sbx

## Out of scope / parked

- Docker Compose / VPS twin under `deploy/` (archived out of tree; revisit later)
- Official-base + apt purge; thin “creds-only” Pi mixins as Hub workarounds
- Baking `gh` alone (or any preference CLI) into core “for parity”
- Daily image rebuilds for agent CLI churn
- Hot-swapping the agent binary mid-session

## Compatibility

- Kits: schemaVersion `"1"`
- CLI: `sbx` >= floor in `cli/internal/sbxcompat`
