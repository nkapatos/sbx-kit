# Template: shell-mise-docker

**Tag:** `local/sbx-shell-mise-docker:latest`  
**Base:** `docker/sandbox-templates:shell-docker` (via [`bake.env`](bake.env))  
**Bake:** [`templates/_bake`](../_bake/)  
**Use:** base image for Pi / Hermes sandbox kits (and plain shell + mise)

```bash
sbx-kit template load --engine docker shell-mise-docker
sbx run shell --template local/sbx-shell-mise-docker:latest \
  --kit "$(brew --prefix)/share/sbx-kit/kits/mise-workspace" .
```
