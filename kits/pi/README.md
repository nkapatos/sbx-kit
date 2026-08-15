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

## Recipes in this tree

```bash
sbx secret set deepseek          # mixin deepseek-creds
sbx-kit run pi --yes             # official shell image, no -t
sbx-kit image load --engine docker kit-core
sbx-kit run kit-pi --yes         # same kit, recipe pins kit-core
```

`kit-pi` may warn that kind `pi` ≠ template name; attach is still Pi.
`mise-workspace` is on `kit-pi` so install can `mise install node@22` when
the lean image has no npm.
