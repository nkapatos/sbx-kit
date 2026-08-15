# Pi (sandbox kit)

`kind: sandbox` kit for [Pi](https://pi.dev) (`@earendil-works/pi-coding-agent`).
This is the agent definition, not a mixin.

## Recipe vs kit kind

A **recipe** is always `sbx kind + kits + optional image`. Kit `kind:` only
changes what the first `sbx run` argument is:

| Kit kind | Recipe | What sbx sees |
| --- | --- | --- |
| mixin | `cursor` + `agent-workspace` | `sbx run cursor --kit agent-workspace` |
| sandbox | `pi` + `pi` + mixins | `sbx run pi --kit pi --kit agent-workspace …` |

Do **not** write `sbx_agent: shell` for this kit. The sandbox kit's `name: pi`
is the kind. Official shell is the **image** (`sandbox.image`), not the kind.

Credentials live in kits (`network` + `proxyManaged` / `credentials.sources`).
The recipe only lists kits; the host stores values with `sbx secret set`.

## Install paths

| Recipe | Floor | How Pi gets on PATH |
| --- | --- | --- |
| `pi` | Official `docker/sandbox-templates:shell` (has npm) | `npm install -g` |
| `kit-pi` | `kit-core` (mise binary only, no Node) | `mise use -g node@22 pnpm` then `pnpm add -g`; symlink `/usr/local/bin/pi` |

The kit detects the lean floor (`/usr/local/bin/mise` + `/mise`) and does **not**
fall back to npm there. `mise-workspace` must stay on the `kit-pi` recipe.

```bash
sbx secret set deepseek
sbx-kit run pi --yes
sbx-kit image load --engine docker kit-core
sbx-kit run kit-pi --yes
```

`kit-pi` may warn that kind `pi` ≠ template name; attach is still Pi.
