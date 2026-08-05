# Base workplace image (no agent binary).

**Tag:** `local/sbx-kit-base:latest`  
**FROM:** `debian:bookworm-slim`

Includes: bash, ca-certificates, curl, git, sudo (passwordless for `agent`),
tini, xz-utils. User `agent` UID/GID 1000.

HTTPS apt bootstrap: slim has no CAs; cleartext apt is often blocked. The
Dockerfile briefly disables HTTPS verify to install `ca-certificates`, then
uses normal apt.

Build:

```bash
./deploy/scripts/build.sh
# or:
docker build -t local/sbx-kit-base:latest -f deploy/base/Dockerfile deploy/base
```
