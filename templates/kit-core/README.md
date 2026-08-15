# Template: kit-core

**Tag:** `local/sbx-kit-core:latest`  
**FROM:** `debian:bookworm-slim`  
**sbx agent:** `shell`

## Why this image

Official `docker/sandbox-templates:*` ship language toolchains we will not
maintain. This floor is intentionally thin:

| In image | Out of image (on purpose) |
| --- | --- |
| sbx glue, `agent`+sudo, mise **binary** | Project languages → **mise** |
| git (+ lfs), curl, small modern utils | Preference CLIs (`gh`, `glab`, …) → **kits** / in-box update |
| `/etc/sbx-agent-env.sh` + persistent env | Docker Engine → future **`-docker`** variant |
| | Agent CLIs → **agent layers** (layout bootstrap only) |

## Docker layer cache

Rebuild only the layer you change: CA → essential OS → modern utils → user →
env glue → mise. See comments in `Dockerfile`.

## Use

```bash
sbx-kit image load --engine docker kit-core
sbx-kit run kit-core --yes
```

Agent layers (e.g. [kit-cursor](../kit-cursor/)) `FROM` this image.
