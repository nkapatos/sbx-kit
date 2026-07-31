# Pi kit — BROKEN (parked)

> **Status: broken / do not use via `sbx-kit run --agent pi`.**
>
> Create fails with opaque sbx `500 failed to run sandbox container` after
> multiple approaches (dedicated `pi-mise-docker` image; kit setup script on
> `shell-mise`). Further sbx-kit BYO experiments are parked.
>
> **Next plan:** plain `Dockerfile` / Docker Compose for local + VPS deploy,
> translating the working `shell-mise` bake floor — not more kit/template loops.

Experimental artifacts left in-tree for reference only:

| Path | Notes |
| --- | --- |
| `spec.yaml` | sandbox kit (image shell-mise, install → setup script) |
| `files/.../setup-pi.sh` | mise node + npm Pi install (unproven on create) |
| `templates/pi-mise-docker/` | earlier baked-image attempt (also unused) |

Working floor that *does* matter: `templates/shell-mise-docker` + `_bake`
(apt + mise). Reuse that when writing the Docker/Compose path.

Oh My Pi (`omp`) remains out of scope.
