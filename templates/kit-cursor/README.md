# Template: kit-cursor

**Tag:** `local/sbx-kit-cursor:latest`  
**FROM:** `local/sbx-kit-core:latest`  
**sbx agent:** `cursor`

Bootstrap install of the Cursor agent CLI on the lean **floor** (not via
kit-shell). **Do not** rebuild this image for every Cursor lab release —
refresh the agent from the **host** before attach. Updating the binary under
a live session is unsafe.

```bash
sbx-kit image load --engine docker kit-cursor
# if blocked: sbx policy allow network downloads.cursor.com
sbx-kit run kit-cursor --yes
```

`image load kit-cursor` docker-builds kit-core first (not imported).

Kits: `mise-workspace`, `agent-workspace`. Preference CLIs (`gh`, `glab`, …)
stay in kits / in-box updates — not this layer.
