# Template: kit-core

**Tag:** `local/sbx-kit-core:latest`  
**FROM:** `debian:bookworm-slim`  
**sbx agent:** `shell`

Lean floor for sbx-kit: `agent` user, `/etc/sandbox-persistent.sh` glue, **mise**
binary. No project language toolchains (Node/Go/Java/Rust) — mise + kits own those.

```bash
sbx-kit template load --engine docker kit-core
sbx-kit run --agent kit-core --yes
```

Agent layers (e.g. [kit-cursor](../kit-cursor/)) `FROM` this image. Same floor is
the intended base for Compose/VPS (`deploy/` convergence).
