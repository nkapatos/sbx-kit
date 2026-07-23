# Templates

Sandbox images live here. Thin directories set `bake.env` (`BASE_IMAGE=…`); the shared body is [`_bake/`](_bake/).

```bash
sbx-kit template load --engine docker <name>       # Docker Desktop / Colima
sbx-kit template load --engine container <name>    # Apple container + skopeo
```

Tag convention: `local/sbx-<role>-<capability>[-<runtime>]:<tag>`  
Images are **linux** (arm64/amd64), not Darwin — they run inside the sbx microVM.

## Shared bake

| Path | Role |
| --- | --- |
| [_bake](_bake/) | Strip conflicting languages; agent CLIs; mise binary; non-interactive UX |

## Thin images

| Path | Tag | sbx agent | Notes |
| --- | --- | --- | --- |
| [cursor-mise-docker](cursor-mise-docker/) | `local/sbx-cursor-mise-docker:latest` | `cursor` | `bake.env` → cursor-agent-docker; pair with [`mise-workspace`](../kits/mise-workspace/) |
| [opencode-mise-docker](opencode-mise-docker/) | `local/sbx-opencode-mise-docker:latest` | `opencode` | stub thin image |
| [shell-mise-docker](shell-mise-docker/) | `local/sbx-shell-mise-docker:latest` | `shell` | base for hermes/pi kits |
