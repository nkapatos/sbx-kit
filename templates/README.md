# Templates

Sandbox images live here. Thin directories set `bake.env` (`BASE_IMAGE=…`); the shared body is [`_bake/`](_bake/).
Templates with their own `Dockerfile` extend a parent image (load the parent first).

```bash
sbx-kit template load --engine docker <name>       # Docker Desktop / Colima
sbx-kit template load --engine container <name>    # Apple container + skopeo
```

Tag convention: `local/sbx-<role>-<capability>[-<runtime>]:<tag>`  
Images are **linux** (arm64/amd64), not Darwin — they run inside the sbx microVM.

## Shared bake

| Path | Role |
| --- | --- |
| [_bake](_bake/) | Strip conflicting languages; agent CLIs; mise; **neovim**; sqlite3/xz; non-interactive UX |

## Thin images

| Path | Tag | sbx agent | Notes |
| --- | --- | --- | --- |
| [cursor-mise-docker](cursor-mise-docker/) | `local/sbx-cursor-mise-docker:latest` | `cursor` | bake.env → cursor-agent-docker |
| [opencode-mise-docker](opencode-mise-docker/) | `local/sbx-opencode-mise-docker:latest` | `opencode` | bake.env → opencode-docker |
| [shell-mise-docker](shell-mise-docker/) | `local/sbx-shell-mise-docker:latest` | `shell` | generic BYO / parent for agent layers |
| [pi-mise-docker](pi-mise-docker/) | `local/sbx-pi-mise-docker:latest` | `pi` | extends shell-mise; Node + official Pi |
| [hermes-mise-docker](hermes-mise-docker/) | `local/sbx-hermes-mise-docker:latest` | `hermes` | extends shell-mise; Hermes `--skip-browser` |
| [cursor-mise-ide](cursor-mise-ide/) | `local/sbx-cursor-mise-ide:latest` | `cursor` | extends cursor-mise; IDE install TODO |

Load order for layered images: parent first (`shell-mise-docker` or `cursor-mise-docker`), then the child.

Until Hub publishes images, `template load` is the local import path. CI/registry publish is a follow-up.
