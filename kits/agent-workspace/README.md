# Agent workplace kit

Runtime-agnostic **mixin** for sandbox-first agent workplace conventions and
**portable state** used by `sbx-kit` lifecycle commands.

Host `~/.cursor` is **not** mounted. Prefer:

1. Project-stamped README / conventions via `sbx-kit init`
2. This kit for portable dirs + `sbx-kit-state` pack/unpack + short agentContext
3. Host vault via `sbx-kit` (`~/.local/share/sbx-kit/profiles/…`) — never the git tree

`portable/` is created on purpose and may look “empty” aside from a seeded
README — that is the agent docs/refs root, not a broken mount.

Detach the agent session before `rm --keep-state` / `upgrade` so SQLite WALs
can checkpoint cleanly (`sbx-kit-state pack` also tries `wal_checkpoint`).

## Layout (inside the VM)

| Path | Role |
| --- | --- |
| `/home/agent/.local/share/sbx-kit/portable/` | Agent-local docs/refs (exportable; starts with README) |
| `/home/agent/.local/share/sbx-kit/state.manifest` | INCLUDE/EXCLUDE list for pack |
| `/home/agent/.local/share/sbx-kit/bin/sbx-kit-state` | Helper (also installed to `/usr/local/bin`) |

## Why kit (not bake) for config

The **manifest and layout** change as agents evolve; kits apply at create /
`sbx kit add` without rebaking images. The helper script ships with this kit
and is installed onto `PATH` at startup. A future bake floor is optional if we
want the binary present even without the kit.

## Pairing

Catalog **default** is this kit only. Stack `mise-workspace` on custom images
that ship `/usr/local/bin/mise` (`kit-shell`, `kit-cursor`, `kit-pi`). Official
Hub recipes do not get mise — they keep the image's own toolchains.

At startup this kit writes `/etc/sbx-kit/floor.md` from live probes (Hub vs
custom). That file is not part of `--keep-state` exports.
