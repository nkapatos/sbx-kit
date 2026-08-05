# Deploy (Docker / Compose for VPS)

Plain Docker images and Compose for a **24/7 agent workplace** outside Docker
AI Sandboxes. On a VPS you do **not** need `sbx` / `sbx-kit` — only Docker
Engine, this tree, and the scripts below.

**Start here for a new VPS:** [`docs/vps-setup.md`](docs/vps-setup.md)

| Path | Role |
| --- | --- |
| [`docs/vps-setup.md`](docs/vps-setup.md) | Host + Cursor user + bring-up checklist |
| [`HARDENING.md`](HARDENING.md) | Threat model + tiers (laptop / vps / locked) |
| [`base/`](base/) | Hardened floor image |
| [`agents/pi/`](agents/pi/) | Pi layer |
| [`egress/`](egress/) | Allowlisted tinyproxy |
| [`compose/`](compose/) | Stack + `overlays/vps.yaml` / `locked.yaml` |
| [`policy/`](policy/) | Network + runtime intent |
| [`scripts/`](scripts/) | build, compose wrapper, status, host prep |

## Scripts

| Script | Purpose |
| --- | --- |
| `scripts/vps-host-prep.sh` | Print/apply Cursor Linux user + docker group |
| `scripts/env-init.sh` | Create `compose/.env` (mode 600) |
| `scripts/build.sh` | Build base + pi + egress images |
| `scripts/compose.sh` | `docker compose` with vps (optional locked) files |
| `scripts/status.sh` | ps + egress allow/deny probe |
| `scripts/gen-egress-filter.sh` | Policy YAML → tinyproxy filter |
| `scripts/install-systemd-user.sh` | User systemd unit for reboot persistence |

## Laptop smoke

```bash
./deploy/scripts/build.sh
./deploy/scripts/compose.sh run --rm --entrypoint pi pi --version
```

## VPS (short form)

```bash
./deploy/scripts/env-init.sh          # edit .env: WORKSPACE + ANTHROPIC_API_KEY
./deploy/scripts/build.sh
./deploy/scripts/compose.sh up -d
./deploy/scripts/status.sh
```

Full checklist and Cursor-user model: [`docs/vps-setup.md`](docs/vps-setup.md).
