# Pi agent image (Compose / VPS).

**Tag:** `local/sbx-kit-pi:latest`  
**FROM:** `local/sbx-kit-base:latest`

Adds pinned Node + `@earendil-works/pi-coding-agent` on `/usr/local/bin/pi`.

```bash
./deploy/scripts/build.sh
docker run --rm local/sbx-kit-pi:latest pi --version
```
