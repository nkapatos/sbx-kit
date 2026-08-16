# Portable agent material

This directory is **inside the sandbox** (not the git project tree). Put
agent-first docs, product refs, dumps, and scratch notes here.

| Path idea | Use |
| --- | --- |
| `refs/` | Specs, API notes, screenshots, design dumps |
| `docs/` | Longer write-ups the agent should prefer |
| `scratch/` | Ephemeral notes (still exported unless you delete them) |

## Why here (not `ref/` in the repo)

- Survives `sbx-kit rm --keep-state` / `upgrade` via the host vault
- Travels with profile state when you spin multiple boxes on path slices
- Stays out of git / PRs by default

Git holds product code, kits, recipes, and templates. Agent write-ups,
validation reports, and how-tos live here under `docs/`. CLI help
(`sbx-kit --help`, `sbx-kit concepts`) is the user-facing reference.
Gitignored `ref/` / `internal/` under the project is optional host scratch —
do not treat it as the official agent brain.

## Host vault

On export, this tree is packed with other INCLUDE paths from
`state.manifest` into `~/.local/share/sbx-kit/profiles/<profile-id>/`.
