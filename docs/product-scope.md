# Product scope (sbx-kit example tree)

This document is the guardrail for what belongs in **this repository’s**
templates, kits, and catalog. sbx-kit the CLI can compose any tree via
`SBX_TREE` or (later) remote registries — that does not expand what we ship
as first-party examples.

## Goals

1. **Own lean templates** that work with sbx (and later Compose/VPS) — not
   purge/strip fat official `docker/sandbox-templates:*` bases.
2. **Manage kits** around those images (thin mixins; agent binaries in images).
3. **Lifecycle** — stable names, portable agent state, recreate/upgrade without
   losing workplace material; later: export the same floor to Compose.
4. Keep the brew/share tree **small and teachable**.

## Layering rule

| Layer | Owns | Examples |
| --- | --- | --- |
| **Template (image)** | Lean floor + optional agent binary | `kit-core`, `kit-cursor` (next: `kit-pi`) |
| **Mixin kit** | Workplace conventions, allowlists, optional credentials | `mise-workspace`, `agent-workspace`, `lsp-mise`, `apt-extras`, `pi` (creds only) |
| **Project** | Language pins, app code | `mise.toml`, repo `AGENTS.md` |
| **Remote / `SBX_TREE`** | Forks, team stacks, taste | Oh My Pi, Playwright FE, host nvim config copy |

**kit-core** ships the mise **binary**. **mise-workspace** activates shims /
allowlists. Project languages stay on mise — never apt toolchains in the floor.
Agent CLIs (Cursor, Pi, …) are **baked layers**, not create-time kit installs.

## In scope (examples we maintain)

- `templates/kit-core`, `templates/kit-cursor`
- Mixins: mise-workspace, agent-workspace; optional lsp-mise, apt-extras; pi creds
- CLI lifecycle + host vault; sbx version floor for tested CLIs
- `deploy/` VPS path converging on the same floor
- Docs that state core vs kits and this scope

## Out of scope / abandoned

- Extending fat official templates then `apt purge` languages
- First-party thins on `cursor-agent-docker` / `shell-docker` / …
- Installing agent binaries via kit `commands.install` (prefer image layers)
- Hermès until core + cursor + pi layers are boring

## Compatibility

- Kits: **schemaVersion `"1"`** on released sbx (see `kits/README.md`)
- CLI: `sbx-kit` requires Docker `sbx` >= floor in `cli/internal/sbxcompat`
