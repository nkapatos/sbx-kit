# Portable material (inside the sandbox, not the git project)

Put agent-first notes, refs, and scratch here. This tree is what
`sbx-kit box state export` packs by default.

| Idea | Use |
|------|------|
| `refs/` | Specs, dumps |
| `docs/` | Longer notes |
| `scratch/` | Ephemeral |

The project worktree is for product code. This directory survives
`sbx-kit box rm --keep-state` / `box upgrade` via the host profile archive.
Do not put secrets here.
