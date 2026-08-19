# sbx-kit skill (draft)

Agent-agnostic bridge between sbx-kit and Docker sbx. Spec and verify commands are
still experimental — see `sbx-kit experimental --help`.

## Layers

| Layer | Commands |
|-------|----------|
| Catalog | `setup`, `catalog add\|ls\|status\|update` |
| Recipes | `recipes`, `recipes image …`, `recipes verify` (stub) |
| Box | `box run`, `box bindings`, `box check`, `box upgrade`, `box rm`, `box state` |
| Project | `project readme` |

## Glossary

- **catalog** — path from `setup` (default `~/sbx-kit-catalog`)
- **directory** — child of catalog with `recipes/`, optional `kits/`, `images/`
- **recipe** — `<dir>/<name>` from `recipes/agents.yaml`
- **kit** — sbx create-time YAML under `kits/` (sbx kit spec v1 today, v2 migration via sbx)
- **image** — Dockerfile or registry tag under `images/`

## Workflow

```bash
sbx-kit setup
sbx-kit catalog add <url>
sbx-kit recipes
sbx-kit box run <dir>/<name> --yes
sbx-kit box bindings
```

## Compose a recipe (interim)

Edit `recipes/agents.yaml` in a catalog directory. Kits stack per entry plus directory defaults.

## Validate (experimental)

```bash
sbx-kit experimental verify recipe [id]
sbx-kit experimental verify kit [dir]
sbx-kit experimental spec
```

Kit schema authority: `sbx kit verify` (sbx v2; v1 kits supported until migrated).

## Host vs box

- **Host**: catalog path, secrets (`sbx secret set`), project bindings, state archives
- **Box**: `sbx run`, workplace under `/home/agent`
