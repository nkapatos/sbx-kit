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
| Agent layer | Intentional layout bootstrap only | Cursor (and later other agents) — **not** daily releases |
| Kits / in-box / CLI refresh | Preference + churn | `gh`/`glab`/…, lang tools via mise, agent refresh |

**Agent updates run from the host before attach** (new models / CLI bits). Do not
hot-swap mid-session.

```bash
# Hub (stock kind, no custom image pin):
sbx-kit run cursor --yes

# Custom image (after image load or image pull):
sbx-kit image load --engine docker kit-core
sbx-kit image load --engine docker kit-cursor
sbx-kit run kit-cursor --yes
```

| Dir | Image tag | Recipe | Role |
| --- | --- | --- | --- |
| [kit-core](kit-core/) | `local/sbx-kit-core:latest` | `kit-core` | Floor; cache-split Dockerfile |
| [kit-cursor](kit-cursor/) | `local/sbx-kit-cursor:latest` | `kit-cursor` | Bootstrap Cursor layout on core |

Recipes: stock `shell` / `cursor`; custom `kit-core` / `kit-cursor`. Further
baked agents on kit-core (local or registry) come later — Hub extras use
sandbox kits, not this tree’s old workaround mixins.
