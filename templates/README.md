# Templates

Sandbox images live here. First-party templates are **our** lean floor — not
extensions of fat `docker/sandbox-templates:*` bases.

```bash
sbx-kit template load --engine docker kit-core
sbx-kit template load --engine docker kit-cursor   # after kit-core
```

Tag convention: `local/sbx-<name>:latest`  
Images are **linux** (arm64/amd64) — they run inside the sbx microVM.

| Path | Tag | sbx agent | Notes |
| --- | --- | --- | --- |
| [kit-core](kit-core/) | `local/sbx-kit-core:latest` | `shell` | Minimal core; mise binary; VPS floor later |
| [kit-cursor](kit-cursor/) | `local/sbx-kit-cursor:latest` | `cursor` | Extends kit-core; Cursor agent CLI |

Load `kit-core` before `kit-cursor`. Recipes: `kit-core`, `cursor` / `kit-cursor`.
Pi will be a `kit-pi` layer later (same pattern).

`ResolveBuild` still honors optional `bake.env` → sibling `_bake` for external
`SBX_TREE` layouts; this repo no longer ships that pattern.
