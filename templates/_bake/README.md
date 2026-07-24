# Shared bake (`templates/_bake`)

Dockerfile + UX files reused by every thin `*-mise-docker` template via `bake.env`.

```bash
sbx-kit template load --engine docker cursor-mise-docker
```

Thin templates only set `BASE_IMAGE=…`; they do not copy this Dockerfile.

**Portable state:** `sbx-kit-state` and its manifest live in [`kits/agent-workspace`](../../kits/agent-workspace/), not this bake — kit config can evolve without rebaking. A bake floor for the helper is optional later if we want it present even without the kit.
