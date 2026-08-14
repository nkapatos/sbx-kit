# Templates

Sandbox images live here.

1. **Own floor (preferred):** [`kit-core/`](kit-core/) — debian + sbx glue + mise, no language bloat.
2. **Legacy official-based:** thin dirs with `bake.env` → shared [`_bake/`](_bake/) on `docker/sandbox-templates:*` (fat bases; leave as-is).

```bash
sbx-kit template load --engine docker <name>       # Docker Desktop / Colima
sbx-kit template load --engine container <name>    # Apple container + skopeo
```

Tag convention: `local/sbx-<name>:latest`  
Images are **linux** (arm64/amd64), not Darwin — they run inside the sbx microVM.

## Own floor

| Path | Tag | sbx agent | Notes |
| --- | --- | --- | --- |
| [kit-core](kit-core/) | `local/sbx-kit-core:latest` | `shell` | Minimal core; base for agent layers + VPS |
| [kit-cursor](kit-cursor/) | `local/sbx-kit-cursor:latest` | `cursor` | Extends kit-core; Cursor agent CLI only |

Load `kit-core` before `kit-cursor`.

## Legacy (official bases + `_bake`)

| Path | Tag | sbx agent | Notes |
| --- | --- | --- | --- |
| [_bake](_bake/) | — | — | Shared body for bake.env thins (mise on fat official) |
| [cursor-mise-docker](cursor-mise-docker/) | `local/sbx-cursor-mise-docker:latest` | `cursor` | bake.env → cursor-agent-docker |
| [opencode-mise-docker](opencode-mise-docker/) | `local/sbx-opencode-mise-docker:latest` | `opencode` | bake.env → opencode-docker |
| [shell-mise-docker](shell-mise-docker/) | `local/sbx-shell-mise-docker:latest` | `shell` | bake.env → shell-docker |
| [cursor-mise-ide](cursor-mise-ide/) | `local/sbx-cursor-mise-ide:latest` | `cursor` | extends cursor-mise; stub |

Load `kit-core` first; agent layers will extend it. Legacy thins still need their official parent via bake.

Until Hub publishes images, `template load` is the local import path.
