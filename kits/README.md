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
| [deepseek-creds](deepseek-creds/) | mixin | **Trial.** Hub-path proxy creds for `api.deepseek.com` (no agent). |

**Hub path:** attach kits to the stock `sbx` agent template (see recipes
`shell`, `cursor`) — no custom image pin. Extra agents (Pi, etc.) belong as
**sandbox kits** on an official shell template, with secrets via `sbx secret` /
proxyManaged. Create-time `run` prints `sbx secret set …`; `sbx-kit check` shows declared
services and passes through `sbx secret ls`. See `sbx-kit concepts`.

**Custom path:** images under `templates/` (`kit-core` → `kit-cursor`, …).
`template load` for local build/import; pin a registry tag in the recipe once published.
Use `mise-workspace` only on images that ship `/usr/local/bin/mise`.

## Composition

Recipes in [`config/agents.yaml`](../config/agents.yaml). Catalog defaults list
`mise-workspace` + `agent-workspace`; Hub example recipes often override kits
(e.g. `agent-workspace` only).
