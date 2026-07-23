# Template: cursor-mise-docker

**Tag:** `local/sbx-cursor-mise-docker:latest`  
**Base:** `docker/sandbox-templates:cursor-agent-docker` (via [`bake.env`](bake.env))  
**Bake:** [`templates/_bake`](../_bake/)  
**Agent:** `cursor`  
**Kits:** [`mise-workspace`](../../kits/mise-workspace/) (mixin); optional [`agent-workspace`](../../kits/agent-workspace/)  
**CLI:** `sbx-kit run cursor`

## Guarantees

- Shared bake: strip conflicting language toolchains; agent CLIs; mise binary; non-interactive UX  
- Does **not** set `ENV PATH`  
- Mise activate/shims come from the **mise-workspace** kit (not this image alone)

## Load and run

```bash
sbx-kit template load --engine docker cursor-mise-docker
# or: --engine container

sbx template ls   # expect *sbx-cursor-mise-docker*

sbx-kit run cursor /path/to/project
```

After create, once per sandbox:

```bash
sbx exec -it <sandbox> -- bash -lc 'cd /path/to/project && mise trust mise.toml; mise install; mise ls'
```
