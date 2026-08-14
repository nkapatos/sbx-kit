# Template: kit-cursor

**Tag:** `local/sbx-kit-cursor:latest`  
**FROM:** `local/sbx-kit-core:latest`  
**sbx agent:** `cursor`

Cursor agent CLI on the lean kit-core floor (no apt language toolchains).

```bash
sbx-kit template load --engine docker kit-core
sbx-kit template load --engine docker kit-cursor
sbx-kit run --agent kit-cursor --yes
```

Build needs `downloads.cursor.com` reachable from the build host. If blocked:

```bash
sbx policy allow network downloads.cursor.com
```

Kits (recipe): `mise-workspace`, `agent-workspace`.

Legacy fat recipe remains as `cursor` → `sbx-cursor-mise-docker`.
