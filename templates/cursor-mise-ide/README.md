# Template: cursor-mise-ide

**Tag:** `local/sbx-cursor-mise-ide:latest`  
**Parent:** `local/sbx-cursor-mise-docker:latest` (shared bake + cursor-agent)  
**Agent:** `cursor`  
**Kits:** same as cursor (`mise-workspace`, `agent-workspace`)  
**Status:** scaffolding — IDE install layer still TODO

Extends the custom cursor+mise image for **in-box Cursor IDE** use (human
workflow alongside the agent). Secrets stay on the host via `sbx secret`.

```bash
sbx-kit template load --engine docker cursor-mise-docker
sbx-kit template load --engine docker cursor-mise-ide
# recipe remains stub until IDE install is verified:
# sbx-kit run --agent cursor-ide --yes
```

Heavy FE stacks (Playwright, browsers) stay **kits**, not this image.
