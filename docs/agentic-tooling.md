# Agentic sandbox tooling

**Audience:** contributors and users of this toolkit.  
**Goal:** own **lean core** image + thin agent layers, plus runtime-agnostic mixin
kits. CLI stays for lifecycle, migrate, and (later) Compose/VPS export.

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
        ├── kits/lsp-mise | apt-extras | pi(creds)   # optional
        │
        ├── templates/kit-cursor      # Cursor agent CLI
        └── templates/kit-pi          # next — Pi baked via mise on core

# VPS: deploy/ converges on the same floor (Compose); CLI export later
```

| Concern | Lives in | Notes |
| --- | --- | --- |
| sbx glue + mise binary | **kit-core** | No language apt packages |
| Agent binary | **Thin layer** on kit-core | kit-cursor now; kit-pi next |
| Mise activate / allowlist | **mise-workspace** | Runtime-agnostic |
| Portable agent state | **agent-workspace** | Host vault via CLI |
| Project language versions | **Project `mise.toml`** | Never bake into the image |
| Official fat bases + apt purge | **Abandoned** | Do not revive |

**Hard rules**

1. `mise-workspace` is not Cursor-specific.
2. Do not set `environment.variables.PATH` in kits.
3. Do not install agent CLIs in kit `commands.install` — bake image layers.
4. Never put bash completion scripts in `/etc/sandbox-persistent.sh`.

---

## First-party templates

| Image tag | Role | Recipe |
| --- | --- | --- |
| `local/sbx-kit-core:latest` | Lean floor | `kit-core` |
| `local/sbx-kit-cursor:latest` | + Cursor agent | `cursor` / `kit-cursor` |

```bash
sbx-kit template load --engine docker kit-core
sbx-kit template load --engine docker kit-cursor
# Cursor tarball: allow downloads.cursor.com if policy blocks
sbx-kit run --agent cursor --yes
```

---

## Follow-ups

| Item | Status |
| --- | --- |
| kit-core + kit-cursor | Done |
| Drop official-base thins | Done |
| Harden kit-core | Done (sudo group, locales, lean utils) |
| kit-pi image layer + pi creds mixin | Next |
| Rebuild kit-cursor on new core | Host: reload kit-core then kit-cursor |
| SSH auth socket | Open |
| Cursor download domains on a mixin | Open |
| CLI → Compose export | Later |
