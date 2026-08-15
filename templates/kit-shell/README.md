# Template: kit-shell

**Tag:** `local/sbx-kit-shell:latest`  
**FROM:** `local/sbx-kit-core:latest`  
**sbx agent:** `shell`

Emulates official `docker/sandbox-templates:shell`. Same idea: this image is
the **basis**, and you add stuff with **kits** (Pi, apt-extras, …) — you do not
bake another image for that.

The floor is [kit-core](../kit-core/). This layer adds **almost nothing**: tini
as PID 1 and a login bash. No extra packages, no Node, no agent CLI.

Baked agents (Cursor) `FROM` kit-core, not this image.

```bash
sbx-kit image load --engine docker kit-shell
sbx-kit run kit-shell --yes
sbx-kit run kit-pi --yes
```

`image load kit-shell` docker-builds kit-core first (not imported).
