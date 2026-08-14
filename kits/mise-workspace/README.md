# Mise workspace kit

Runtime-agnostic **mixin** for any mise-prepared sbx template (`kit-core`, `cursor`, …). Configures mise; does not install language toolchains or bind to one agent image.

## What it does

- Sets `MISE_DATA_DIR` / `MISE_CONFIG_DIR` / `MISE_CACHE_DIR` to `/mise`
- Ensures persistent shells activate mise and prepend `/mise/shims` (never via `environment.variables.PATH`)
- Allows network domains for mise backends and common language registries
- Injects agentContext: trust → install → prune after pin changes

## What it does not do

- Does **not** choose a template image or agent entrypoint (`kind: mixin`)
- Does **not** run `mise install` at create time
- Does **not** pin language versions (project `mise.toml` owns that)

## Prerequisites

The template must provide the **mise** binary at `/usr/local/bin/mise` (see `templates/kit-core`).

## Apply

Prefer the catalog CLI (pulls this kit automatically):

```bash
sbx-kit run --agent cursor --yes
sbx-kit run --agent kit-core --yes
```

Long form:

```bash
sbx run cursor --template local/sbx-kit-cursor:latest \
  --kit "$(brew --prefix)/share/sbx-kit/kits/mise-workspace" .
```
