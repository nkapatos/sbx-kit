# Templates

First-party images are **our** lean floor (sbx + later VPS) — not fat official
`docker/sandbox-templates:*` bases.

## Why this split

| Layer | Rebuild when | Holds |
| --- | --- | --- |
| OS / essential apt | Debian major/minor | bash, curl, git, locales, tini, sudo |
| Modern utils | Occasional floor tools | fd, rg, jq, git-lfs, sqlite3, … |
| sbx glue + mise binary | Glue or mise binary bumps | user, persistent-env, `/usr/local/bin/mise` |
| Agent layer | Intentional layout bootstrap only | Cursor/Pi install paths — **not** daily releases |
| Kits / in-box / CLI update | Preference + churn | `gh`/`glab`/…, lang tools via mise, agent refresh |

**Agent updates run from the host before attach** (new models / CLI bits). Do not
hot-swap the running agent binary mid-session.

```bash
sbx-kit template load --engine docker kit-core
sbx-kit template load --engine docker kit-cursor   # after kit-core
sbx-kit run --agent cursor --yes
```

| Path | Tag | sbx agent | Notes |
| --- | --- | --- | --- |
| [kit-core](kit-core/) | `local/sbx-kit-core:latest` | `shell` | Floor; cache-split Dockerfile; VPS later |
| [kit-cursor](kit-cursor/) | `local/sbx-kit-cursor:latest` | `cursor` | Bootstrap Cursor layout on core |

Recipes: `kit-core`, `cursor` / `kit-cursor`. Next agent layer: `kit-pi`.

`ResolveBuild` still honors optional `bake.env` → sibling `_bake` for external
`SBX_TREE` layouts; this repo does not ship that pattern.
