# Agentic sandbox tooling

**Audience:** contributors and users of this toolkit.  
**Goal:** recipes + kits on **official** or **custom** templates; portable state;
optional lean **OS → utils → glue/mise → agent** images. CLI owns lifecycle and
migrate. VPS/Compose packaging is parked for later (same image idea, different
runtime).

---

## Why kits on official templates first

Docker’s supported customization model is: pick an **official sandbox
template**, then attach **kits**. sbx-kit’s job is to make that easy via
**recipes** (which agent + which kits + optional local/remote tag), place kit
paths correctly, and preserve `/home/agent` state across recreate and host
moves.

Agents Hub doesn’t ship (e.g. Pi) fit as **sandbox kits** on an official shell
template (`startup` + trust args), with credentials through `sbx secret` /
proxyManaged — not by stripping or rebaking Hub images.

Local lean images (`kit-core` → baked agent layers) are for when you want a thin
floor or a custom agent image you pull **local or remote**. Same approach for
future Pi/Hermès layers on core; not a substitute for the Hub+kits path.

---

## Architecture

```text
                    ┌── Hub / registry template (sbx pulls)
sbx agent + kits ───┤     sandbox kits start extra agents on shell, etc.
                    └── custom: kit-core → kit-cursor | (later other agent layers)
                              pull local (template load) or remote registry

kits: mise-workspace | agent-workspace | apt-extras | …   (mixins)
CLI: recipes, placement, state export/import/upgrade; later recipe registry
```

| Concern | Lives in | Update cadence |
| --- | --- | --- |
| Official agent image | Hub via `sbx` | Upstream |
| OS + floor utils + sbx glue + mise binary | **kit-core** (custom) | Rare / occasional |
| Agent binary layout (custom) | **Agent layer** bootstrap pin | Rebake only when layout changes |
| Agent version / new models (custom) | **Host refresh before attach** | As needed — not mid-session |
| Extra agents on Hub shell | **Sandbox kit** + secrets | User / community kits |
| Languages | **mise** (+ project `mise.toml`) when image has mise | In-box / agent |
| Preference CLIs (`gh`, `glab`, …) | **Kits** | User preference |

**Hard rules**

1. `mise-workspace` is not Cursor-specific; never set `environment.variables.PATH`.
2. Do not install preference CLIs or languages into kit-core “because official has them.”
3. Agent layers bootstrap install paths; **do not** chase daily Hub rebuilds for agent churn.
4. Never put secrets or bash completions in `/etc/sandbox-persistent.sh`.
5. Keep `AGENT_CLI_CREDENTIAL_STORE=memory` on Cursor layers.
6. Do not keep Hub-workaround kits that only paper over official templates.

---

## First-party local templates (optional)

| Image tag | Role | Recipe |
| --- | --- | --- |
| `local/sbx-kit-core:latest` | Lean floor | `kit-core` |
| `local/sbx-kit-cursor:latest` | + Cursor bootstrap | `cursor` / `kit-cursor` |

Hub example recipes: `shell-hub` (stock shell + deepseek-creds trial),
`cursor-hub` (stock cursor + agent-workspace).

```bash
sbx-kit secrets --agent shell-hub
sbx-kit run --agent shell-hub --yes
# local custom:
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
| `template ls` + `secrets` guide from kit schema | Done (thin) |
| Recipe/kit discovery (remote tree or registry) | Open |
| One-command local `template load` for agent images | Open |
| Example sandbox kit on official shell (community pattern) | Open |
| Baked agents beyond cursor on kit-core (local/remote pull) | Later |
| Kits for preference CLIs (`gh`/`glab`/…) | Open |
| Host agent refresh (distinct from `upgrade`) | Open |
| SSH sock / host env allowlist | Open |
| DinD `-docker` variant | Later |
| Compose / VPS from same floor | Parked (was `deploy/`) |
