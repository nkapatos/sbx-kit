# Template: kit-core

**Tag:** `local/sbx-kit-core:latest`  
**FROM:** `debian:bookworm-slim`  
**sbx agent:** `shell`

Lean floor for sbx-kit: `agent` user (sudo), `/etc/sandbox-persistent.sh` +
`/etc/sbx-agent-env.sh` (UX only), **mise** binary, small agent utilities
(`fd`, `rg`, `jq`, `git-lfs`, `sqlite3`, `socat`, …). No project language
toolchains — mise + kits own those.

**Not in this image:** host SSH keys, DinD, language runtimes, secrets.
Prefer HTTPS + `sbx secret` / proxy; keep Docker Engine on a future `-docker`
variant if needed.

```bash
sbx-kit template load --engine docker kit-core
sbx-kit run --agent kit-core --yes
```

Agent layers (e.g. [kit-cursor](../kit-cursor/)) `FROM` this image. Same floor is
the intended base for Compose/VPS (`deploy/` convergence).
