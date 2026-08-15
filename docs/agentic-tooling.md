# Agentic sandbox tooling

**Audience:** contributors and users of this toolkit.  
**Goal:** recipes + kits on **official** or **local** templates; portable state;
optional lean **OS → utils → glue/mise → agent** images. CLI owns lifecycle and
migrate. VPS/Compose packaging is parked for later (same image idea, different
runtime).

---

## Why kits on official templates first

Docker’s supported customization model is: pick an **official sandbox
template**, then attach **kits**. sbx-kit’s job is to make that easy via
**recipes** (which agent + which kits + optional local tag), place kit paths
correctly, and preserve `/home/agent` state across recreate and host moves.

Local lean images exist when you want a thin floor without Hub toolchain bloat:
**mise** for languages, **kits** for taste tools, **host-side agent refresh**
for Cursor/Pi (many releases per week — do not rebake daily).

---

## Architecture

```text
                    ┌── Hub / registry template (sbx pulls)
sbx agent + kits ───┤
                    └── local/sbx-kit-*:  templates/kit-core → kit-cursor | kit-pi

kits/mise-workspace | agent-workspace | apt-extras | …   (create-time mixins)
CLI: recipes, placement, state export/import/upgrade
```

| Concern | Lives in | Update cadence |
| --- | --- | --- |
| Official agent image | Hub via `sbx` | Upstream |
| OS + floor utils + sbx glue + mise binary | **kit-core** (local) | Rare / occasional |
| Agent binary layout | **Agent layer** bootstrap pin | Rebake only when layout changes |
| Agent version / new models (local) | **Host refresh before attach** | As needed — not mid-session |
| Languages | **mise** (+ project `mise.toml`) when image has mise | In-box / agent |
| Preference CLIs (`gh`, `glab`, …) | **Kits** | User preference |

**Hard rules**

1. `mise-workspace` is not Cursor-specific; never set `environment.variables.PATH`.
2. Do not install preference CLIs or languages into kit-core “because official has them.”
3. Agent layers bootstrap install paths; **do not** chase daily Hub rebuilds for agent churn.
4. Never put secrets or bash completions in `/etc/sandbox-persistent.sh`.
5. Keep `AGENT_CLI_CREDENTIAL_STORE=memory` on Cursor layers.

---

## First-party local templates (optional)

| Image tag | Role | Recipe |
| --- | --- | --- |
| `local/sbx-kit-core:latest` | Lean floor | `kit-core` |
| `local/sbx-kit-cursor:latest` | + Cursor bootstrap | `cursor` / `kit-cursor` |

Hub example recipe: `cursor-hub` (stock `cursor` agent + `agent-workspace`, no local build).

```bash
sbx-kit run --agent cursor-hub --yes
# local:
sbx-kit template load --engine docker kit-core
sbx-kit template load --engine docker kit-cursor
sbx-kit run --agent cursor --yes
```

---

## Follow-ups

| Item | Status |
| --- | --- |
| kit-core cache-split layers | Done |
| kit-cursor bootstrap + host refresh policy | Done (docs) |
| Hub-first recipes + clearer CLI navigation | In progress |
| One-command local `template load` for agent images | Open |
| kit-pi layer | Next |
| Kits for preference CLIs (`gh`/`glab`/…) | Open |
| Host agent refresh (distinct from `upgrade`) | Open |
| SSH sock / host env allowlist | Open |
| DinD `-docker` variant | Later |
| Compose / VPS from same floor | Parked (was `deploy/`) |
