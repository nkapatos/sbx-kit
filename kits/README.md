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
| [agent-workspace](agent-workspace/) | mixin | Portable state + `sbx-kit-state` + agentContext. **Catalog default** on every recipe. |
| [lsp-mise](lsp-mise/) | mixin | **Optional.** Box-level `/mise/config.toml` for LSPs/helpers. |
| [apt-extras](apt-extras/) | mixin | **Optional.** Small apt packages. |
| [deepseek-creds](deepseek-creds/) | mixin | **Trial.** Hub-path proxy creds for `api.deepseek.com` (no agent). |

## Sandbox kits

| Directory | Kind | Used with |
| --- | --- | --- |
| [pi](pi/) | sandbox | Official shell (`pi`, npm -g) or kit-core (`kit-pi`, mise node@22 + pnpm). `sbx_agent` is `pi`, not `shell`. |

A sandbox kit **is** the sbx kind (`name:` → first arg to `sbx run`). Mixins still stack. Credentials stay in kits; `sbx secret set` is host-side.

**Hub path:** stock kinds (`shell`, `cursor`) plus mixins, or a sandbox kit whose `sandbox.image` is an official template. Create-time `run` prints `sbx secret set …`; `sbx-kit check` shows declared services. See `sbx-kit concepts`.

**Custom path:** images under `templates/` (`kit-core` → `kit-cursor`, …).
`sbx-kit image load` for local build/import; `sbx-kit image pull` for a registry
tag. Pin the tag in the recipe once published. Same sandbox kit as Hub; the
recipe `-t` overrides `sandbox.image`.
Use `mise-workspace` only on images that ship `/usr/local/bin/mise`.

## Composition

Recipes in [`config/agents.yaml`](../config/agents.yaml). Catalog defaults list
`agent-workspace`; recipe `kits:` are extra (mise-workspace, pi, deepseek-creds, …).
