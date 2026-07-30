# Template: opencode-mise-docker

**Tag:** `local/sbx-opencode-mise-docker:latest`  
**Base:** `docker/sandbox-templates:opencode-docker` (via [`bake.env`](bake.env))  
**Bake:** [`templates/_bake`](../_bake/) (mise, neovim, agent CLIs, non-interactive UX)  
**Agent:** `opencode`  
**Kits:** [`mise-workspace`](../../kits/mise-workspace/) + [`agent-workspace`](../../kits/agent-workspace/)  
**CLI:** `sbx-kit run --agent opencode`

```bash
sbx-kit template load --engine docker opencode-mise-docker
sbx-kit run --agent opencode --yes
```

After create, once per sandbox (if the project pins tools):

```bash
sbx exec -it <sandbox> -- bash -lc 'cd /path/to/project && mise trust mise.toml; mise install; mise ls'
```
