#!/usr/bin/env bash
# Print (or optionally apply) host prep for a dedicated VPS.
# Does not install Cursor CLI. Review before --apply.
set -euo pipefail

USER_NAME=cursor
APPLY=0
PRINT=1

usage() {
  cat <<EOF
Usage: $0 [--user NAME] [--print] [--apply]

  --user NAME   Linux account for Cursor CLI + Compose (default: cursor)
  --print       Print recommended commands (default)
  --apply       Run safe subset: create user, add to docker group (requires root)

Install Docker Engine yourself from Docker’s official docs before --apply.
EOF
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --user) USER_NAME="$2"; shift 2 ;;
    --print) PRINT=1; shift ;;
    --apply) APPLY=1; shift ;;
    -h|--help) usage; exit 0 ;;
    *) echo "unknown arg: $1" >&2; usage; exit 1 ;;
  esac
done

cat <<EOF
# sbx-kit deploy — VPS host prep (user=${USER_NAME})

## Manual (always)

1. SSH hardening (key-only, no root password login)
2. Firewall: allow SSH; default deny inbound
3. Install Docker Engine + compose plugin:
   https://docs.docker.com/engine/install/
4. unattended-upgrades / automatic security updates
5. Clone sbx-kit as ${USER_NAME}; see deploy/docs/vps-setup.md

## Account

EOF

if [[ "$APPLY" -eq 1 ]]; then
  if [[ "$(id -u)" -ne 0 ]]; then
    echo "--apply requires root" >&2
    exit 1
  fi
  if ! id -u "$USER_NAME" >/dev/null 2>&1; then
    adduser --disabled-password --gecos 'Cursor CLI / Compose' "$USER_NAME"
  else
    echo "user $USER_NAME already exists"
  fi
  if getent group docker >/dev/null; then
    usermod -aG docker "$USER_NAME"
    echo "added $USER_NAME to docker"
  else
    echo "WARNING: group docker missing — install Docker Engine first" >&2
  fi
  install -d -o "$USER_NAME" -g "$USER_NAME" -m 755 "/home/${USER_NAME}/workspaces"
  echo "created /home/${USER_NAME}/workspaces"
  echo "Done. Next: sudo -iu $USER_NAME  # install Cursor CLI, clone repo"
else
  cat <<EOF
# suggested commands (run as root after Docker is installed):

adduser --disabled-password --gecos 'Cursor CLI / Compose' ${USER_NAME}
usermod -aG docker ${USER_NAME}
# optional: usermod -aG sudo ${USER_NAME}
install -d -o ${USER_NAME} -g ${USER_NAME} -m 755 /home/${USER_NAME}/workspaces
sudo -iu ${USER_NAME}

EOF
fi
