# VPS hardening (24/7 agentic coding)

Threat model and Compose tiers for [`deploy/`](.). Not sbx-level isolation
(microVM); this is **container + host hygiene** for a dedicated VPS.

## Threat model (what we assume)

| Asset | Risk if agent is malicious or confused |
| --- | --- |
| Host kernel / Docker socket | Full VPS compromise |
| Sibling containers | Lateral movement |
| Cloud metadata (`169.254.169.254`) | Steal instance credentials |
| API keys in env | Exfil over open egress |
| Workspace mount | Wipe/exfil project; not host `$HOME` if you mount only the repo |
| Package installs (sudo/apt) | Persistence, crypto miners, reverse shells |

**Out of scope for this tree:** defending a multi-tenant kernel, or stopping a
rooted Docker daemon. Do **not** mount `/var/run/docker.sock`. Do **not** run
`--privileged` / `pid: host` / `network: host`.

## Hardening tiers

| Tier | Compose entry | sudo in image | Egress | Rootfs | Restart |
| --- | --- | --- | --- | --- | --- |
| **laptop** | `pi.compose.yaml` alone | yes | open bridge | writable | `no` |
| **vps** (default for 24/7) | `pi.compose.yaml` + `overlays/vps.yaml` | yes | **allowlisted proxy** | writable | `unless-stopped` |
| **locked** | + `overlays/locked.yaml` | build without sudo | allowlisted proxy | read-only + tmpfs | `unless-stopped` |

Promote laptop → vps before leaving a box unattended. Use **locked** when the
agent should not self-provision packages (host/bake owns tools).

## Always-on rules (every tier)

1. User `1000:1000` — never root as the process user.
2. `no-new-privileges:true`, `cap_drop: ALL` (add back only what you need).
3. No published ports unless you intentionally expose a UI.
4. Secrets via Compose `env_file` / VPS secret store — never `ENV` in Dockerfile.
5. Bind-mount **only** the project workspace (and optional agent-state volume).
6. Resource caps: memory, CPUs, pids, log rotation.
7. Block link-local / metadata in policy (enforce via proxy + host firewall).

## Egress model (vps / locked)

```text
  [pi container] --HTTP(S)_PROXY--> [tinyproxy] --allowlist--> Internet
         |                              |
    network: agent                 network: egress
    (internal: true)               (outbound OK)
```

Allowlist source of truth: [`policy/pi.network.yaml`](policy/pi.network.yaml).
Regenerate proxy filter:

```bash
./deploy/scripts/gen-egress-filter.sh
```

Host firewall (nftables/ufw) should still deny the agent bridge from reaching
`169.254.169.254` and other tenants; the proxy is defense in depth, not the
only control.

## Ops checklist (first VPS deploy)

- [ ] Build images on a trusted machine; pull by digest on the VPS when possible
- [ ] `WORKSPACE=` absolute path to **one** project tree
- [ ] Copy `compose/.env.example` → `.env`; chmod 600; real API key only there
- [ ] Start with **vps** overlay; confirm proxy logs show denied junk hosts
- [ ] `docker compose … logs -f` + disk alerts; log max-size enabled
- [ ] Unattended-upgrades / patch host OS; keep Docker Engine current
- [ ] Separate Linux user for compose; no agent in `docker` group on host
- [ ] Backup workspace + `agent-state` volume; treat them as sensitive

## Future CLI

`sbx-kit` can later **report** drift: `policy/*.yaml` vs Compose overlays
(caps, restart, proxy filter membership). Lifecycle stay Compose/systemd.
