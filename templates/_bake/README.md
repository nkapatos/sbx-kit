# Shared bake (`templates/_bake`)

Dockerfile + UX files reused by every thin `*-mise-docker` template via `bake.env`.

```bash
sbx-kit template load --engine docker cursor-mise-docker
```

Thin templates only set `BASE_IMAGE=…`; they do not copy this Dockerfile.

**Includes:** agent CLIs, native build helpers, **mise** binary, **neovim** (in-box / headless ACP), `sqlite3`, `xz-utils`, non-interactive UX (`EDITOR=true`).

**Does not include:** project language pins, GUI IDEs, host editor configs, Playwright/browsers.

**Portable state:** `sbx-kit-state` and its manifest live in [`kits/agent-workspace`](../../kits/agent-workspace/), not this bake — kit config can evolve without rebaking.
