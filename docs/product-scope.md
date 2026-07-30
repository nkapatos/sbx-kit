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
| **Template (image)** | Agent binary, default runtime deps, shared bake floor | `cursor-mise-docker`, `pi-mise-docker`, `hermes-mise-docker` |
| **Mixin kit** | Workplace conventions, allowlists, optional capabilities | `mise-workspace`, `agent-workspace`, `lsp-mise`, `apt-extras` |
| **Sandbox kit** | Agent id / entrypoint / creds / network / short context — **not** heavy installs | `kits/pi`, `kits/hermes` |
| **Project** | Language pins, app code | `mise.toml`, repo `AGENTS.md` |
| **Remote / `SBX_TREE`** | Forks, team stacks, taste | Oh My Pi, Playwright FE, host nvim config copy |

Bake when install is slow, flaky, or must be pinned. Keep kits for wiring that
should change without rebaking. Create-time `commands.install` for large agents
is an anti-pattern in this tree.

## In scope (examples we maintain)

- Shared `_bake` (mise binary, neovim, agent CLIs, non-interactive UX)
- Thin templates: cursor, opencode, shell (generic parent), **pi**, **hermes**
- Mixins: mise-workspace, agent-workspace; optional lsp-mise, apt-extras
- Thin pi/hermes sandbox kits (creds/network/context only)
- CLI lifecycle + host vault; sbx version floor for tested CLIs
- Docs that state bake vs kit and this scope

## Out of scope (not first-party examples)

- **Oh My Pi (`omp`) and other Pi/Hermes forks** — remote registry or local tree
- Nested sandboxes inside sbx (e.g. pi-docker-sandbox)
- Host `~/.cursor` / nvim config sync inside **sbx-kit** (optional kits only if ever)
- Playwright / browsers / heavy FE in default images (team kits or remote recipes)
- Baking every popular `pi install` package
- Authored kit-spec v2 until released `sbx kit validate` accepts it
- Full CI → Hub publish and compose export (follow-ups; design should stay compatible)

## Pi / Hermes specifically

- Image: `_bake` on `shell-docker` via parent `shell-mise-docker`, then agent layer
  with the CLI forced onto `/usr/local/bin` (sbx probes PATH for the agent binary)
- Kit: thin sandbox kit (creds/network/context); agent name = kit name (`pi` / `hermes`)
- **Expected warning:** template built for `shell` but agent is `pi`/`hermes` —
  Docker templates do not register new runtimes; the sandbox kit does. Ignore the
  warning when `/usr/local/bin/pi` (or `hermes`) exists in the imported image.
- **Hard failure** `agent binary not found`: rebuild parent + child, confirm with
  `docker run --rm --entrypoint which <tag> pi`, then `sbx template load` again
- Extensions: in-box (`pi install`, Hermes setup) or separate recipes elsewhere
- `shell-mise-docker` remains the **generic** parent / BYO floor, not the product recipe for Pi/Hermes
- Oh My Pi / other forks: remote registry or `SBX_TREE` only

## Compatibility

- Kits: **schemaVersion `"1"`** on released sbx (see `kits/README.md`)
- CLI: `sbx-kit` requires Docker `sbx` >= floor in `cli/internal/sbxcompat`

When in doubt: prefer a smaller example tree and a clear extension path over
shipping every popular agent variant.
