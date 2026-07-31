# Product scope (sbx-kit example tree)

This document is the guardrail for what belongs in **this repository’s**
templates, kits, and catalog. sbx-kit the CLI can compose any tree via
`SBX_TREE` or (later) remote registries — that does not expand what we ship
as first-party examples.

## Goals

1. **Ease custom templates** on official Docker AI Sandboxes bases (`_bake` +
   thin images).
2. **Manage kits** around those images (mixins + thin sandbox kits).
3. **Lifecycle** — stable names, portable agent state, recreate/upgrade without
   losing workplace material.
4. Keep the brew/share tree **small and teachable** so people can copy the
   pattern and extend it.

## Layering rule

| Layer | Owns | Examples |
| --- | --- | --- |
| **Template (image)** | Shared bake floor; official agent binaries when the starter has them | `cursor-mise-docker`, `opencode-mise-docker`, `shell-mise-docker` |
| **Mixin kit** | Workplace conventions, allowlists, optional capabilities | `mise-workspace`, `agent-workspace`, `lsp-mise`, `apt-extras` |
| **Sandbox kit (BYO)** | Agent id / entrypoint / creds / network + create-time setup script | `kits/pi` (on shell-mise); Hermes same pattern when proven |
| **Project** | Language pins, app code | `mise.toml`, repo `AGENTS.md` |
| **Remote / `SBX_TREE`** | Forks, team stacks, taste | Oh My Pi, Playwright FE, host nvim config copy |

Bake the workplace floor (apt + mise + tools). For agents **not** in official
starters, prefer a kit dropped setup script + one-line `commands.install` over
a second template that reinstalls Node/etc. Dedicated agent images remain an
optional later optimization (pin/Hub), not the default BYO path.

## In scope (examples we maintain)

- Shared `_bake` (mise binary, neovim, agent CLIs, non-interactive UX)
- Thin templates: cursor, opencode, shell (workplace floor)
- Mixins: mise-workspace, agent-workspace; optional lsp-mise, apt-extras
- CLI lifecycle + host vault; sbx version floor for tested CLIs
- Docs that state bake vs kit and this scope

## Parked / broken (do not use)

- **`pi` / `hermes` sbx recipes** — stubbed in `config/agents.yaml`. Kit and
  dedicated-image attempts failed create (`500 failed to run sandbox container`
  / agent binary not found). **Next:** plain Dockerfile + Compose for VPS,
  reusing the shell-mise bake ideas — not more sbx BYO kit loops.

## Out of scope (not first-party examples)

- **Oh My Pi (`omp`) and other Pi/Hermes forks** — remote registry or local tree
- Nested sandboxes inside sbx (e.g. pi-docker-sandbox)
- Host `~/.cursor` / nvim config sync inside **sbx-kit** (optional kits only if ever)
- Playwright / browsers / heavy FE in default images (team kits or remote recipes)
- Baking every popular `pi install` package
- Authored kit-spec v2 until released `sbx kit validate` accepts it
- Full CI → Hub publish and compose export (follow-ups; design should stay compatible)

## Pi / Hermes specifically

**Broken under sbx-kit today** (catalog `stub: true`). Do not run these recipes.
Planned recovery is **Docker / Compose** (local + VPS), not further kit churn.
Keep `shell-mise` + `_bake` as the workplace reference for that translation.
Oh My Pi / other forks: remote registry or `SBX_TREE` only.

## Compatibility

- Kits: **schemaVersion `"1"`** on released sbx (see `kits/README.md`)
- CLI: `sbx-kit` requires Docker `sbx` >= floor in `cli/internal/sbxcompat`

When in doubt: prefer a smaller example tree and a clear extension path over
shipping every popular agent variant.
