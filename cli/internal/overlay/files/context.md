# sbx-kit overlay (this box)

This file is **sbx-kit land**. The CLI writes it on create/attach. It is **not**
a user kit and is **not** kit `agentContext`.

Start here. Open the files below only when you need them.

## Docs (overlay)

| Path | When to read |
|------|----------------|
| `/etc/sbx-kit/floor.md` | What *this* image has (Hub vs custom; mise/gh/npm/go) |
| `/etc/sbx-kit/cli.md` | Host `sbx-kit` verbs for this CLI version |
| `/home/agent/.local/share/sbx-kit/portable/README.md` | What `box state` packs by default |

Helper on PATH: `sbx-kit-state` (`check` / `pack` / `unpack`).
Manifest: `$SBX_KIT_STATE_MANIFEST` (default under `/home/agent/.local/share/sbx-kit/`).

## How to read floor.md

- **kind: official** — Docker Hub template. Use the image's preinstalled tools.
  Do not apt/mise a second toolchain unless floor says `mise: yes`.
- **kind: custom** — lean floor (`/etc/sbx-kit-agent` present). Languages via
  **mise** only when `mise: yes`. Do not assume Hub CLIs (`gh`, global `npm`, Go).

Do not assume mise, Cursor, or a SQLite DB. Checkpoint WAL only if a `.db`
exists under an INCLUDE path.

## User land (kits)

User kits are optional. Their `agentContext` (v1) / `agentInstructions.content`
(v2) is injected by **sbx** (often `kits-agent-context/` next to the memory
file). That is extra taste. It does not replace this overlay.

A user kit **may** point at this file so the box agent finds platform docs
without inlining them. Example one-liner for a kit `agentContext`:

```markdown
Platform (sbx-kit overlay, not this kit): read `/etc/sbx-kit/context.md`.
```

## Secrets and git

Prefer HTTPS + `sbx secret` / proxy injection. Do not run `gh auth login` or
overwrite `proxy-managed` sentinels. Do not put secrets in `portable/`.
