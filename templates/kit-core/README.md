# Template: kit-core

**Tag:** `local/sbx-kit-core:latest` (Docker only)  
**FROM:** `debian:bookworm-slim`  
**Role:** parent image. Not imported into sbx.

You `FROM` this to bake sandbox images (`kit-shell`, `kit-cursor`, later agents)
and, later, Docker-on-VPS hosts. **Do not** `sbx-kit image load` or
`sbx-kit run` this directory.

## Why this image

Official `docker/sandbox-templates:*` ship language toolchains we will not
maintain. This floor is intentionally thin. **kit-cursor** bakes the Cursor CLI
here. **kit-shell** is the minimum empty image you actually load into sbx so
you can emulate `sbx run shell` and add kits (e.g. Pi).

| In image | Out of image (on purpose) |
| --- | --- |
| sbx glue, `agent`+sudo, mise **binary** | Project languages → **mise** |
| git (+ lfs), curl/wget, ssh client, zip/zstd | Preference CLIs (`gh`, `glab`, …) → **kits** / in-box update |
| C toolchain (`build-essential`) for native addons | Docker Engine → future **`-docker`** / VPS host |
| rg/fd/jq/sqlite3 + small debug utils | Agent CLIs → **agent layers** (layout bootstrap only) |

## Docker layer cache

Rebuild only the layer you change: CA → essential OS → modern utils → user →
env glue → mise. See comments in `Dockerfile`.

## Use

```bash
# Builds kit-core in Docker automatically, then imports the child into sbx:
sbx-kit image load --engine docker kit-shell
sbx-kit image load --engine docker kit-cursor
sbx-kit run kit-shell --yes
sbx-kit run kit-cursor --yes
```

Children: [kit-shell](../kit-shell/) (minimum loadable image), [kit-cursor](../kit-cursor/).
