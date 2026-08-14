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

## Mixins (runtime-agnostic)

| Directory | Kind | Used with |
| --- | --- | --- |
| [mise-workspace](mise-workspace/) | mixin | Any template with `/usr/local/bin/mise` (kit-core). Allowlists, `MISE_*`, activate startup. Never sets `environment.variables.PATH`. |
| [agent-workspace](agent-workspace/) | mixin | Portable state + `sbx-kit-state` pack/unpack + agentContext. |
| [lsp-mise](lsp-mise/) | mixin | **Optional.** Box-level `/mise/config.toml` for LSPs/helpers (not project pins). |
| [apt-extras](apt-extras/) | mixin | **Optional.** Small apt packages (not Homebrew). |
| [pi](pi/) | mixin | Install Pi via mise on kit-core; DeepSeek creds. |

Cursor is the **kit-cursor** image + `sbx_agent: cursor` — there is no cursor kit.
Language tooling is the `mise-workspace` mixin.

## Composition

Recipes stack kits in [`config/agents.yaml`](../config/agents.yaml). Defaults are `mise-workspace` + `agent-workspace`. Add `lsp-mise` / `apt-extras` / `pi` per recipe when needed.

Host editor prefs and forks (e.g. Oh My Pi) belong in **optional / remote / `SBX_TREE` kits**, not this example catalog.
