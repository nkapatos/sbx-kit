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

Product agents (**Pi**, **Hermes**) bake the binary into a template; their sandbox kits are **thin** (creds/network/context only). See [docs/product-scope.md](../docs/product-scope.md).

## Mixins (runtime-agnostic)

| Directory | Kind | Used with |
| --- | --- | --- |
| [mise-workspace](mise-workspace/) | mixin | Any mise-prepared template. Allowlists, `MISE_*`, activate startup, agentContext. Never sets `environment.variables.PATH`. |
| [agent-workspace](agent-workspace/) | mixin | Portable state + `sbx-kit-state` pack/unpack + agentContext. |
| [lsp-mise](lsp-mise/) | mixin | **Optional.** Box-level `/mise/config.toml` for LSPs/helpers (not project pins). |
| [apt-extras](apt-extras/) | mixin | **Optional.** Small Ubuntu packages via apt (not Homebrew). |

## Sandbox agents (thin — image owns the binary)

| Directory | Kind | Image | Notes |
| --- | --- | --- | --- |
| [hermes](hermes/) | sandbox | `local/sbx-hermes-mise-docker` | Network + context; no install |
| [pi](pi/) | sandbox | `local/sbx-pi-mise-docker` | Creds + network + context; no install |

There is **no** cursor-mise kit — Cursor specificity is the thin template (`sbx-kit run --agent cursor`). Language tooling is always the `mise-workspace` mixin.

## Composition

Recipes stack kits in [`config/agents.yaml`](../config/agents.yaml). Defaults are `mise-workspace` + `agent-workspace`. Add `lsp-mise` / `apt-extras` (or team kits) per recipe when needed.

Host editor prefs and forks (e.g. Oh My Pi) belong in **optional / remote / `SBX_TREE` kits**, not this example catalog.
