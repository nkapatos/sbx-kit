# Agent workplace

Mixin for portable sandbox state (`sbx-kit-state`, host vault). Catalog default.

Does not mount host `~/.cursor`. Project rules stay in the repo (`sbx-kit init`);
session material stays under `/home/agent`. Detach the agent before
`rm --keep-state` / `upgrade` so SQLite can checkpoint.

At startup writes `/etc/sbx-kit/floor.md` (Hub vs custom probe; not exported).
Details live in `spec.yaml`.
