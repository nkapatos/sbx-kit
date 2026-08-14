# Product scope (sbx-kit example tree)

Guardrail for this repo’s templates, kits, and catalog. CLI can compose any
tree via `SBX_TREE` — that does not expand what we ship as examples.

## Goals

1. **Own lean images** for sbx and (later) VPS — not fat official templates.
2. **Clear update ownership** — rare image rebuilds; mise for langs; kits for
   preference CLIs; **host-side agent refresh** before attach.
3. **Lifecycle** — names, portable state, migrate; later Compose export.
4. Keep the brew/share tree small and teachable.

## Layering rule

| Layer | Owns | Examples |
| --- | --- | --- |
| **kit-core** | OS essentials, modern utils, sbx glue, mise **binary** | `templates/kit-core` |
| **Agent image** | Bootstrap install/layout only | `kit-cursor` (next `kit-pi`) |
| **Mixin kit** | Activate, allowlists, state, preference CLIs, creds | `mise-workspace`, `agent-workspace`, `apt-extras`, `pi` |
| **Project** | Language pins | `mise.toml` |
| **Host / CLI** | Agent binary refresh before run; later update reports | (future `sbx-kit update`) |

Languages → mise. Preference CLIs (`gh`, `glab`, …) → kits / in-box — **not**
special-cased into core. Agent version churn → host refresh, not daily rebake.

## In scope

- `templates/kit-core`, `templates/kit-cursor`
- Mixins: mise-workspace, agent-workspace; optional lsp-mise, apt-extras; pi creds
- CLI lifecycle + vault; docs for why/layers/updates
- `deploy/` converging on the same floor

## Out of scope / abandoned

- Official-base + apt purge
- Baking `gh` alone (or any preference CLI) into core “for parity”
- Daily Hub rebuilds for Cursor/Pi releases
- Hot-swapping the agent binary mid-session
- Hermès until core + cursor + pi are boring

## Compatibility

- Kits: schemaVersion `"1"`
- CLI: `sbx` >= floor in `cli/internal/sbxcompat`
