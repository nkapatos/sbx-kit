# VPS setup guide (Pi deploy stack)

End-to-end path for a dedicated VPS running 24/7 agentic coding with this
repo’s `deploy/` images and Compose. **sbx-kit is not required** on the VPS.

Your intended flow:

1. Create the VPS and apply **host** hardening (SSH, firewall, updates).
2. Create a Linux user that will run **Cursor CLI** (and Docker Compose).
3. With that agent + these docs/scripts, bring up the Pi container stack.

## Roles (keep them separate)

| Who | On the VPS | May use Docker? |
| --- | --- | --- |
| `root` / admin | Host OS, firewall, user creation | yes (bootstrap only) |
| **`cursor`** (or similar) | Cursor CLI, `git`, this repo, Compose | **yes** (`docker` group) |
| container user `agent` (UID 1000) | Inside the Pi image only | **no** — never mount the Docker socket |

Do **not** run Cursor as root. Do **not** put the in-container agent on the
host `docker` group (it never gets a host login).

## 1. Host hardening (before Compose)

Do this as root / cloud-init on a fresh Ubuntu 22.04/24.04 or Debian 12 box:

- [ ] SSH: key-only, disable password auth, non-default user with sudo
- [ ] `ufw` (or nftables): allow `22/tcp` (or your SSH port) only from your IPs if possible; deny public inbound otherwise
- [ ] Block cloud metadata from Docker bridges (see [HARDENING.md](../HARDENING.md)); at minimum do not rely on open egress alone
- [ ] `unattended-upgrades` (or equivalent) for security patches
- [ ] Install Docker Engine + Compose plugin from Docker’s docs (not ancient distro packages)
- [ ] Fail2ban or cloud provider SSH brute-force protection optional but recommended

Helper (prints commands; use `--apply` only after review):

```bash
sudo ./deploy/scripts/vps-host-prep.sh --user cursor --print
# sudo ./deploy/scripts/vps-host-prep.sh --user cursor --apply
```

## 2. Cursor CLI user

```bash
# as root / admin
sudo adduser --disabled-password --gecos 'Cursor CLI' cursor
sudo usermod -aG docker,sudo cursor   # sudo optional; docker required for Compose
sudo -iu cursor
```

Then as `cursor`:

- Install Cursor CLI per Cursor’s docs for Linux.
- Clone this repo (or sync a release tarball) to e.g. `~/sbx-kit`.
- Create a **project** directory that will be the only bind mount, e.g.
  `~/workspaces/my-app` (not `$HOME`).

## 3. Configure secrets and workspace

```bash
cd ~/sbx-kit
./deploy/scripts/env-init.sh
# edit deploy/compose/.env — set ANTHROPIC_API_KEY and WORKSPACE=/home/cursor/workspaces/my-app
chmod 600 deploy/compose/.env
```

`WORKSPACE` must be an **absolute** path owned by `cursor` (UID may differ from
container UID 1000 — if file ownership fights you, either align UIDs or keep
the project mode group-writable carefully; simplest is host UID 1000 = cursor).

## 4. Build images and start (vps tier)

```bash
./deploy/scripts/build.sh
./deploy/scripts/compose.sh up -d
./deploy/scripts/status.sh
```

`compose.sh` always uses `pi.compose.yaml` + `overlays/vps.yaml`.  
Locked tier: `PI_LOCKED=1 ./deploy/scripts/compose.sh up -d`

## 5. Day-2 operations

| Task | Command |
| --- | --- |
| Status | `./deploy/scripts/status.sh` |
| Logs | `./deploy/scripts/compose.sh logs -f pi` |
| Proxy denies | `./deploy/scripts/compose.sh logs -f egress-proxy` |
| Shell in box | `./deploy/scripts/compose.sh exec pi bash -l` |
| Attach TTY | `./deploy/scripts/compose.sh attach pi` |
| Stop | `./deploy/scripts/compose.sh down` |
| Rebuild after policy change | `./deploy/scripts/gen-egress-filter.sh && ./deploy/scripts/build.sh && ./deploy/scripts/compose.sh up -d --build` |

Optional systemd unit (user service) so the stack returns after reboot:

```bash
./deploy/scripts/install-systemd-user.sh   # as cursor; then: systemctl --user enable --now sbx-kit-pi
```

## 6. Sanity checks before walking away

- [ ] `./deploy/scripts/status.sh` shows `pi` and `egress-proxy` healthy/up
- [ ] Proxy log shows **refused** for a junk host (or run the status script’s probe)
- [ ] No ports published on the host for the agent (`docker compose ps`)
- [ ] `.env` is mode `600` and not in git
- [ ] Workspace path is only the project tree

## What this does *not* do

- Provision the cloud VM / DNS / TLS
- Install Cursor CLI for you
- Replace a host firewall
- Give sbx-level microVM isolation

When the checklist above is green, you are good to move on to live agent work
on the VPS.
