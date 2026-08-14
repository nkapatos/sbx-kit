# Agentic sandbox tooling

**Audience:** contributors and users of this toolkit.  
**Goal:** own lean **OS → utils → glue/mise → agent** images + mixin kits. Same
floor for sbx templates and (later) VPS/Compose. CLI owns lifecycle, migrate,
and (later) update reports / Compose export.

---

## Why not official fat templates

Official `docker/sandbox-templates:*` bake language toolchains and preference
CLIs we will not keep current. We omit them: **mise** for languages,
**kits / in-box update** for `gh`/`glab`/taste tools, **host-side agent
refresh** for Cursor/Pi (many releases per week — do not rebake daily).

---

## Architecture

```text
debian:bookworm-slim
        │  (1) CA  (2) essential OS  (3) modern utils
        ▼
   templates/kit-core          # (4) agent user/sudo  (5) persistent-env  (6) mise binary
        │
        ├── kits/mise-workspace          # activate + allowlists
        ├── kits/agent-workspace         # portable state
        ├── kits/apt-extras | lsp-mise | pi(creds)   # optional preference / extras
        │
        ├── templates/kit-cursor         # bootstrap Cursor layout (refresh on host)
        └── templates/kit-pi             # next

# VPS: same kit-core floor via deploy/ / CLI export later
```

| Concern | Lives in | Update cadence |
| --- | --- | --- |
| OS + floor utils + sbx glue + mise binary | **kit-core** (split layers) | Rare / occasional |
| Agent binary layout | **Agent layer** bootstrap pin | Rebake only when layout changes |
| Agent version / new models | **Host refresh before attach** | As needed — not mid-session |
| Languages | **mise** (+ project `mise.toml`) | In-box / agent |
| Preference CLIs (`gh`, `glab`, …) | **Kits** / Compose / in-box | User preference |
| Docker Engine | Future **`-docker`** variant | Not on default core |

**Hard rules**

1. `mise-workspace` is not Cursor-specific; never set `environment.variables.PATH`.
2. Do not install preference CLIs or languages into kit-core “because official has them.”
3. Agent layers bootstrap install paths; **do not** chase daily Hub rebuilds for agent churn.
4. Never put secrets or bash completions in `/etc/sandbox-persistent.sh`.
5. Keep `AGENT_CLI_CREDENTIAL_STORE=memory` on Cursor layers.

---

## First-party templates

| Image tag | Role | Recipe |
| --- | --- | --- |
| `local/sbx-kit-core:latest` | Lean floor | `kit-core` |
| `local/sbx-kit-cursor:latest` | + Cursor bootstrap | `cursor` / `kit-cursor` |

```bash
sbx-kit template load --engine docker kit-core
sbx-kit template load --engine docker kit-cursor
# Cursor tarball: sbx policy allow network downloads.cursor.com
sbx-kit run --agent cursor --yes
```

---

## Follow-ups

| Item | Status |
| --- | --- |
| kit-core cache-split layers | Done |
| kit-cursor bootstrap + host refresh policy | Done (docs) |
| kit-pi layer | Next |
| Kits for preference CLIs (`gh`/`glab`/…) | Open |
| CLI: report/apply mise+apt+agent updates (agent = pre-attach) | Open |
| SSH sock / host env allowlist | Open |
| DinD `-docker` variant | Later |
| Compose / VPS export from same floor | Later |
