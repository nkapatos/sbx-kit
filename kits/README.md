# Kits

Kits are create-time YAML (`spec.yaml`): network allowlists, env, agent notes.
They are **not** images — pass them with `sbx run … --kit /path/to/kit`.

This tree authors **kit-spec v1** (`schemaVersion: "1"`).

**Command shapes (easy to mix up):**

| Block | `command` type |
| --- | --- |
| `commands.install` | **string** (passed to `sh -c`) |
| `commands.startup` | **argv array** (e.g. `[bash, -lc, "…"]`) |

## Mixins

| Directory | Kind | Used with |
| --- | --- | --- |
| [mise-workspace](mise-workspace/) | mixin | Templates with `/usr/local/bin/mise` (kit-core). Allowlists, `MISE_*`, activate. Never sets `environment.variables.PATH`. |
| [agent-workspace](agent-workspace/) | mixin | Portable state + `sbx-kit-state` + agentContext. |
| [lsp-mise](lsp-mise/) | mixin | **Optional.** Box-level `/mise/config.toml` for LSPs/helpers. |
| [apt-extras](apt-extras/) | mixin | **Optional.** Small apt packages. |
| [pi](pi/) | mixin | DeepSeek creds for a future **kit-pi** image (no install). |

Cursor is the **kit-cursor** image + `sbx_agent: cursor`. Agent binaries are
baked in templates; kits stay thin.

## Composition

Recipes in [`config/agents.yaml`](../config/agents.yaml). Defaults: `mise-workspace` + `agent-workspace`.
