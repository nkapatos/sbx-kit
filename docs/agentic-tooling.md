# Agentic sandbox tooling

**Audience:** contributors and users of this toolkit.  
**Goal:** own **lean core** image + thin agent layers, plus runtime-agnostic mixin
kits (`mise-workspace`, `agent-workspace`, optional capability kits). CLI stays
for lifecycle, migrate, and (later) Compose/VPS export.

This file is the source of truth for *what goes where*. Consumer READMEs are secondary.

---

## Architecture

```text
debian:bookworm-slim
        │
        ▼
   templates/kit-core                 # agent user, sbx persistent-env, mise binary
        │                             # no Node/Go/Java/Rust project toolchains
        ├── kits/mise-workspace       # activate + allowlists (not the binary)
        ├── kits/agent-workspace
        ├── kits/lsp-mise | apt-extras | pi   # optional
        │
        └── templates/kit-cursor      # Cursor agent CLI only
                └── sbx_agent: cursor + same mixins

# VPS: deploy/ converges on the same floor (Compose); CLI export later
```

| Concern | Lives in | Notes |
| --- | --- | --- |
| sbx glue + mise binary | **kit-core** | No language apt packages |
| Agent binary (Cursor, …) | **Thin layer** on kit-core | e.g. kit-cursor |
| Mise activate / registry allowlist | **mise-workspace** | Runtime-agnostic |
| Portable agent state | **agent-workspace** | Host vault via CLI |
| Project language versions | **Project `mise.toml`** | Never bake into the image |
| Official fat `sandbox-templates:*` | **Not used** for first-party | Abandoned purge/`_bake` approach |

**Hard rules**

1. `mise-workspace` is **not** Cursor-specific.
2. Do **not** set `environment.variables.PATH` in kits.
3. Do **not** maintain apt-purge lists against official fat bases.
4. Never put bash completion scripts in `/etc/sandbox-persistent.sh`.
5. Default recipes: kit-core or kit-cursor + mise + portable state.

---

## First-party templates

| Image tag | Role | Used with |
| --- | --- | --- |
| `local/sbx-kit-core:latest` | Lean floor | `shell` / recipe `kit-core`, `pi` |
| `local/sbx-kit-cursor:latest` | + Cursor agent | `cursor` / `kit-cursor` |

```bash
sbx-kit template load --engine docker kit-core
sbx-kit template load --engine docker kit-cursor
# Cursor tarball: allow downloads.cursor.com if policy blocks
sbx-kit run --agent cursor --yes
```

---

## Kits

See [kits/README.md](../kits/README.md). Defaults on recipes: `mise-workspace` +
`agent-workspace`. Template must provide `/usr/local/bin/mise` (kit-core).

**After changing pins in mise.toml:** `mise install && mise prune -y` (prefer a
fresh login shell if env looks stale).

---

## Follow-ups

| Item | Status |
| --- | --- |
| kit-core + kit-cursor | Done (proven under sbx) |
| Drop official-base `_bake` thins | Done |
| Pi mixin on kit-core | Wired; harden install/activate later |
| SSH auth socket in lean boxes | Open |
| Cursor download domains on a cursor mixin | Open |
| More agent layers (opencode, …) | Later |
| CLI → Compose export / deploy convergence | Later |
| Hub publish | Later |
