# Kits

Kits are create-time YAML (`spec.yaml`): network allowlists, env, agent notes, or full sandbox agents.
They are **not** images — pass them with `sbx run … --kit /path/to/kit`.

This tree authors **kit-spec v1** (`schemaVersion: "1"`). Released `sbx` (through at
least 0.37) validates/runs v1; clean authored v2 is not in a release yet.

Deprecation WARNs on `network.*` / `proxyManaged` are expected on newer sbx.

**Command shapes (easy to mix up):**

| Block | `command` type |
| --- | --- |
| `commands.install` | **string** (passed to `sh -c`) |
| `commands.startup` | **argv array** (e.g. `[bash, -lc, "…"]`) |

Product agents (**Pi**, …) that are not in official starters use **shell-mise**
plus a sandbox kit with a dropped setup script and a one-line
`commands.install`. See [docs/product-scope.md](../docs/product-scope.md).

## Mixins (runtime-agnostic)

| Directory | Kind | Used with |
| --- | --- | --- |
| [mise-workspace](mise-workspace/) | mixin | Any mise-prepared template. Allowlists, `MISE_*`, activate startup, agentContext. Never sets `environment.variables.PATH`. |
| [agent-workspace](agent-workspace/) | mixin | Portable state + `sbx-kit-state` pack/unpack + agentContext. |
| [lsp-mise](lsp-mise/) | mixin | **Optional.** Box-level `/mise/config.toml` for LSPs/helpers (not project pins). |
| [apt-extras](apt-extras/) | mixin | **Optional.** Small Ubuntu packages via apt (not Homebrew). |

## Sandbox agents (BYO — BROKEN / parked)

| Directory | Kind | Status |
| --- | --- | --- |
| [pi](pi/) | sandbox | **Broken** — stubbed; next = Docker/Compose for VPS |
| [hermes](hermes/) | sandbox | **Broken** — stubbed with Pi |

There is **no** cursor-mise kit — Cursor specificity is the thin template (`sbx-kit run --agent cursor`). Language tooling is always the `mise-workspace` mixin.

## Composition

Recipes stack kits in [`config/agents.yaml`](../config/agents.yaml). Defaults are `mise-workspace` + `agent-workspace`. Add `lsp-mise` / `apt-extras` (or team kits) per recipe when needed.

Host editor prefs and forks (e.g. Oh My Pi) belong in **optional / remote / `SBX_TREE` kits**, not this example catalog.
