# Template: kit-core

**Tag:** `local/sbx-kit-core:latest`  
**FROM:** `debian:bookworm-slim`  
**Role:** parent image. You `FROM` this to bake agent images (Cursor, later others).
It is not the Hub-shell counterpart — that is [kit-shell](../kit-shell/), which
you customize with **kits**.

## Why this image

Official `docker/sandbox-templates:*` ship language toolchains we will not
maintain. This floor is intentionally thin. **kit-cursor** bakes the Cursor CLI
here. **kit-shell** is a tiny attach image on the same floor so you can emulate
`sbx run shell` and add kits (e.g. Pi).

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
# then bake a child, e.g.:
sbx-kit image load --engine docker kit-cursor
sbx-kit run kit-cursor --yes
```

Children: [kit-cursor](../kit-cursor/) (baked agent), [kit-shell](../kit-shell/)
(tinishell + kits). `sbx-kit run kit-core` only smokes the floor.
