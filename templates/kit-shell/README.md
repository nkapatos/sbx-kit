# Template: kit-shell

**Tag:** `local/sbx-kit-shell:latest`  
**FROM:** `local/sbx-kit-core:latest`  
**sbx agent:** `shell`

Custom counterpart to official `docker/sandbox-templates:shell`. The floor is
[kit-core](../kit-core/). This image adds **almost nothing**: tini as PID 1 and
a login bash. No extra packages, no Node, no agent CLI.

Pi (create-time sandbox kit) sits on this image, same as official Pi sits on
Hub shell.

```bash
sbx-kit image load --engine docker kit-core
sbx-kit image load --engine docker kit-shell
sbx-kit run kit-shell --yes
sbx-kit run kit-pi --yes
```
