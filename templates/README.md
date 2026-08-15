# Templates

Optional **custom** images — lean floor when you do not want fat official
`docker/sandbox-templates:*` bases. Develop with `sbx-kit image load`;
**pull a registry tag** with `sbx-kit image pull` and pin it in a recipe for
normal use. Stock recipes need **no** image pin; `sbx` uses the official kind.

| Layer | Rebuild when | Holds |
| --- | --- | --- |
| OS / essential apt | Debian major/minor | bash, curl, git, locales, tini, sudo |
| Modern utils | Occasional floor adds | jq, ripgrep, fd, sqlite3, … |
| sbx glue + mise binary | Glue or mise binary bumps | user, persistent-env, `/usr/local/bin/mise` |
| **kit-shell** | Almost never | tini PID 1 + login bash only |
| Agent layer | Intentional layout bootstrap only | Cursor (and later other agents) — **not** daily releases |
| Kits / in-box / CLI refresh | Preference + churn | `gh`/`glab`/…, lang tools via mise, agent refresh |

**Agent updates run from the host before attach** (new models / CLI bits). Do not
hot-swap mid-session.

```bash
# Hub (stock kind, no custom image pin):
sbx-kit run cursor --yes
sbx-kit run pi --yes

# Custom images (load floor first):
sbx-kit image load --engine docker kit-core
sbx-kit image load --engine docker kit-shell
sbx-kit image load --engine docker kit-cursor
sbx-kit run kit-shell --yes
sbx-kit run kit-pi --yes
sbx-kit run kit-cursor --yes
```

| Dir | Image tag | Recipe | Role |
| --- | --- | --- | --- |
| [kit-core](kit-core/) | `local/sbx-kit-core:latest` | (parent; optional smoke) | `FROM` this to bake images |
| [kit-shell](kit-shell/) | `local/sbx-kit-shell:latest` | `kit-shell` | Hub-shell counterpart: tinishell; add **kits** |
| [kit-cursor](kit-cursor/) | `local/sbx-kit-cursor:latest` | `kit-cursor` | Cursor CLI **FROM kit-core** |

Recipes: stock `shell` / `cursor` / `pi`; custom `kit-shell` / `kit-cursor` / `kit-pi`.
`kit-pi` uses **kit-shell** (same idea as Hub Pi on official shell).
