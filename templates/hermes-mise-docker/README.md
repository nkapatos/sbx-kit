# Template: hermes-mise-docker

**Tag:** `local/sbx-hermes-mise-docker:latest`  
**Parent:** `local/sbx-shell-mise-docker:latest`  
**Adds:** Hermes via `install.sh --skip-browser` → `/usr/local/bin/hermes`  
**Agent:** `hermes`  
**CLI:** `sbx-kit run --agent hermes`

Same shell-family caveat as Pi: the “built for shell” warning is expected; a
missing-binary error means reload the layered image into sbx.

```bash
sbx-kit template load --engine docker shell-mise-docker
sbx-kit template load --engine docker hermes-mise-docker
sbx-kit run --agent hermes --yes
```
