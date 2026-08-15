# Product scope

Guardrail for this example tree. `SBX_TREE` can point the CLI at another tree;
that does not expand what we ship here. Live catalog: `sbx-kit recipes`.

## Goals

1. Recipes + kits on **official Hub kinds** and on **custom images**.
2. Obvious image source (Hub vs `image load` / `image pull`).
3. Portable `/home/agent` state across recreate and host moves.
4. Thin custom images: mise for languages, kits for preference CLIs, host-side
   agent refresh (not mid-session, not daily rebakes).
5. Small brew/share tree. CLI help is the reference, not these files.

## Two paths

| Path | How |
| --- | --- |
| **Shell + kits** | Hub `shell` or `kit-shell`. Add kits. Do not bake. |
| **Image on core** | `FROM kit-core` for images you own. Do not import kit-core into sbx. |

## Layering (custom images)

| Layer | Owns |
| --- | --- |
| **kit-core** | OS, utils, sbx glue, mise binary. Docker `FROM` parent only. Later VPS hosts. |
| **kit-shell** | Tini + login bash. Hub-shell counterpart. |
| **Agent image** | Bootstrap layout only (`kit-cursor` FROM core, not via kit-shell). |
| **Mixin kit** | Activate, allowlists, state, preference CLIs. |
| **Sandbox kit** | Agent Hub doesn’t ship (`kits/pi`). Kind = kit `name:`. |
| **Project** | Language pins (`mise.toml`). |
| **Host / CLI** | Recipes, placement, state. |

**Hard rules:** no languages or preference CLIs in kit-core; no secrets or bash
completions in `/etc/sandbox-persistent.sh`; Cursor layers keep
`AGENT_CLI_CREDENTIAL_STORE=memory`; never set `environment.variables.PATH` in
mise-workspace; don’t paper over Hub images with workaround kits.

## In scope

Kits, recipes, Hub-first create, optional local/registry custom images,
portable state, `check`.

## Out of scope / parked

Compose/VPS twin (same kit-core parent later); official-base + apt purge;
baking `gh` into core; daily agent-image rebuilds; hot-swap agent binaries
mid-session.

## Kit schema

Stay on **`schemaVersion: "1"`**. Upstream v2 exists from sbx **0.36**
(`permissions`, `setup`, `agentInstructions`, list `credentials`). v1 still
loads through the legacy path. This tree’s sbx floor is `sbx-kit version`
(`cli/internal/sbxcompat`). Move to v2 only when that floor is ≥ 0.36 **and**
every `spec.yaml` is rewritten — do not mix v1 fields into a `"2"` file.
