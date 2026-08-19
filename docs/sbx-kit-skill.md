# sbx-kit skill (draft)

Agent-agnostic bridge between sbx-kit and Docker sbx.

## Layers

| Layer | Commands |
|-------|----------|
| Catalog | `setup`, `catalog add\|ls\|status\|update` |
| Recipes | `recipes`, `recipes verify`, `recipes verify kits`, `recipes image …` |
| Box | `box run`, `box bindings`, `box check`, `box upgrade`, `box rm`, `box state` |
| Project | `project readme` |

## Glossary

- **catalog** — path from `setup` (default `~/sbx-kit-catalog`)
- **directory** — child of catalog with `recipes/`, optional `kits/`, `images/`
- **recipe** — `<dir>/<name>` from `recipes/agents.yaml`
- **kit** — sbx create-time YAML under `kits/` (schema owned by sbx)
- **image** — Dockerfile or registry tag under `images/`

## Workflow

```bash
sbx-kit setup
sbx-kit catalog add <url>
sbx-kit recipes
sbx-kit recipes verify
sbx-kit box run <dir>/<name> --yes
sbx-kit box bindings
```

## Verify

```bash
sbx-kit recipes verify              # recipe manifests + sbx kit verify
sbx-kit recipes verify mine/cursor
sbx-kit recipes verify --skip-kits  # manifest only
sbx-kit recipes verify kits mine    # sbx kit verify only
```

Recipe manifests are checked by sbx-kit. Kit specs are checked by sbx. Migrate kits with sbx, not sbx-kit.

`sbx-kit version` shows the required sbx range (`MinVersion`, `MinKitVerify`).

## Host vs box

- **Host**: catalog path, secrets (`sbx secret set`), project bindings, state archives
- **Box**: `sbx run`, workplace under `/home/agent`
