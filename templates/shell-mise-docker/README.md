# Template: shell-mise-docker

**Tag:** `local/sbx-shell-mise-docker:latest`  
**Base:** `docker/sandbox-templates:shell-docker` (via [`bake.env`](bake.env))  
**Bake:** [`templates/_bake`](../_bake/)  
**Use:** generic shell workplace (apt + mise bake floor)

```bash
sbx-kit template load --engine docker shell-mise-docker
# shell agent / mixins only — Pi/Hermes sbx recipes are stubbed (broken).
```

Pi/Hermes will move to plain Docker/Compose for VPS; this bake remains the
reference floor to translate.
