# Shared bake (`templates/_bake`)

Dockerfile + UX files reused by every thin `*-mise-docker` template via `bake.env`.

```bash
sbx-kit template load --engine docker cursor-mise-docker
```

Thin templates only set `BASE_IMAGE=…`; they do not copy this Dockerfile.
