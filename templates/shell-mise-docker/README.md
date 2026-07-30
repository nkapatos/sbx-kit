# Template: shell-mise-docker

**Tag:** `local/sbx-shell-mise-docker:latest`  
**Base:** `docker/sandbox-templates:shell-docker` (via [`bake.env`](bake.env))  
**Bake:** [`templates/_bake`](../_bake/)  
**Use:** generic shell workplace, and **parent image** for Pi / Hermes agent layers

Product agents live on dedicated templates ([`pi-mise-docker`](../pi-mise-docker/), [`hermes-mise-docker`](../hermes-mise-docker/)). This image stays the shared floor — not the Pi/Hermes recipe path.

```bash
sbx-kit template load --engine docker shell-mise-docker
# then, for agents:
sbx-kit template load --engine docker pi-mise-docker
# or: sbx-kit template load --engine docker hermes-mise-docker
```
