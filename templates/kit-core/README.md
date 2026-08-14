# Template: kit-core

**Tag:** `local/sbx-kit-core:latest`  
**FROM:** `debian:bookworm-slim`  
**sbx agent:** `shell`

Lean floor for sbx-kit: `agent` user (sudo), `/etc/sandbox-persistent.sh` glue,
**mise** binary, small agent utilities (`fd`, `rg`, `jq`, `git-lfs`, `sqlite3`,
`socat`, …). No project language toolchains — mise + kits own those.

```bash
sbx-kit template load --engine docker kit-core
sbx-kit run --agent kit-core --yes
```

Agent layers (e.g. [kit-cursor](../kit-cursor/)) `FROM` this image. Same floor is
the intended base for Compose/VPS (`deploy/` convergence).
