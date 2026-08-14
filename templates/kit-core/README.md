# Template: kit-core

**Tag:** `local/sbx-kit-core:latest`  
**FROM:** `debian:bookworm-slim` (not an official `docker/sandbox-templates:*`)  
**sbx agent:** `shell`

Lean floor for sbx-kit: `agent` user, `/etc/sandbox-persistent.sh` glue, **mise**
binary. No project language toolchains (Node/Go/Java/Rust) — mise + kits own those.

```bash
sbx-kit template load --engine docker kit-core
sbx-kit run --agent kit-core --yes
# or: sbx run --template local/sbx-kit-core:latest shell .
```

**Next:** agent layers (e.g. cursor) `FROM local/sbx-kit-core:latest`. Same floor
is the intended base for Compose/VPS export later (`deploy/` convergence).

Official fat templates (`cursor-mise-docker`, `shell-mise-docker`, …) remain as
legacy recipes; this core replaces them for new work.
