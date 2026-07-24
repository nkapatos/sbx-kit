# Kits

Kits are create-time YAML (`spec.yaml`): network allowlists, env, agent notes, or full sandbox agents.
They are **not** images — pass them with `sbx run … --kit /path/to/kit`.

## Mixins (runtime-agnostic)

| Directory | Kind | Used with |
| --- | --- | --- |
| [mise-workspace](mise-workspace/) | mixin | Any mise-prepared template. Allowlists, `MISE_*`, activate startup, agentContext (trust/install/prune). Never sets `environment.variables.PATH`. |
| [agent-workspace](agent-workspace/) | mixin | Portable state layout + `sbx-kit-state` pack/unpack + agentContext. MCP/skills still TBD. |

## Sandbox agents (BYO on shell-mise)

| Directory | Kind | Image | Notes |
| --- | --- | --- | --- |
| [hermes](hermes/) | sandbox | `local/sbx-shell-mise-docker` | **Stub.** Entrypoint/install/auth TBD. |
| [pi](pi/) | sandbox | `local/sbx-shell-mise-docker` | **Stub.** Entrypoint/install/auth TBD. |

There is **no** cursor-mise kit — Cursor specificity is the thin template (`sbx-kit run cursor`). Language tooling is always the `mise-workspace` mixin.

## Planned

| Name | Intent |
| --- | --- |
| Optional capability mixins | cloud CLIs, Playwright, etc. |

Reuse mixins across templates; use sandbox kits only when defining an agent entrypoint.
