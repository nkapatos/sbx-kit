#!/usr/bin/env bash
# Install a systemd --user unit that starts the Pi VPS compose stack on login/boot
# (enable lingering so it survives logout: sudo loginctl enable-linger "$USER").
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
UNIT_DIR="${XDG_CONFIG_HOME:-$HOME/.config}/systemd/user"
UNIT="$UNIT_DIR/sbx-kit-pi.service"
COMPOSE="$ROOT/deploy/scripts/compose.sh"

mkdir -p "$UNIT_DIR"
cat > "$UNIT" <<EOF
[Unit]
Description=sbx-kit Pi agent (Compose vps tier)
After=network-online.target docker.service
Wants=network-online.target

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=$ROOT/deploy/compose
EnvironmentFile=-$ROOT/deploy/compose/.env
ExecStart=$COMPOSE up -d
ExecStop=$COMPOSE down
TimeoutStartSec=0

[Install]
WantedBy=default.target
EOF

systemctl --user daemon-reload
echo "wrote $UNIT"
echo "Enable lingering (once, as root): sudo loginctl enable-linger $USER"
echo "Then: systemctl --user enable --now sbx-kit-pi.service"
echo "Status: systemctl --user status sbx-kit-pi.service"
