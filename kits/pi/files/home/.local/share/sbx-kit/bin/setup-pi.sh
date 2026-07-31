#!/usr/bin/env bash
# BROKEN / parked: Pi sbx kit create failed (500). Kept for reference only.
# Next: plain Docker/Compose for VPS — do not rely on this under sbx-kit run.
#
# Create-time Pi install for shell-mise boxes (unproven).
# Invoked by kits/pi commands.install (one-liner). Idempotent via marker.
set -euo pipefail

SHARE="${HOME:-/home/agent}/.local/share/sbx-kit"
MARKER="${SHARE}/.pi-installed"
SCRIPT_NAME="setup-pi"

log() { printf '%s: %s\n' "$SCRIPT_NAME" "$*"; }

if [[ -f "$MARKER" ]] && command -v pi >/dev/null 2>&1; then
  log "already installed at $(command -v pi)"
  exit 0
fi

export MISE_DATA_DIR="${MISE_DATA_DIR:-/mise}"
export MISE_CONFIG_DIR="${MISE_CONFIG_DIR:-/mise}"
export MISE_CACHE_DIR="${MISE_CACHE_DIR:-/mise/cache}"
export PATH="/mise/shims:${PATH:-/usr/local/bin:/usr/bin:/bin}"

if [[ ! -x /usr/local/bin/mise ]]; then
  log "ERROR: /usr/local/bin/mise missing — load shell-mise-docker first"
  exit 1
fi
# Activate for this install shell (mise-workspace startup may not have run yet).
eval "$(/usr/local/bin/mise activate bash)"

NODE_VERSION="${SBX_PI_NODE_VERSION:-22}"
PI_PACKAGE="${SBX_PI_PACKAGE:-@earendil-works/pi-coding-agent}"

log "mise use -g node@${NODE_VERSION}"
mise use -g "node@${NODE_VERSION}"
hash -r 2>/dev/null || true

if ! command -v npm >/dev/null 2>&1; then
  log "ERROR: npm not on PATH after mise use node@${NODE_VERSION}"
  mise which npm || true
  exit 1
fi

log "npm install -g --ignore-scripts ${PI_PACKAGE}"
npm install -g --ignore-scripts "${PI_PACKAGE}"
hash -r 2>/dev/null || true

if ! command -v pi >/dev/null 2>&1; then
  # mise/npm globals sometimes land outside a minimal PATH; resolve and link.
  PI_BIN="$(mise which pi 2>/dev/null || true)"
  if [[ -z "$PI_BIN" || ! -x "$PI_BIN" ]]; then
    log "ERROR: pi binary not found after npm install"
    npm root -g || true
    ls -la "$(npm root -g)/../bin" 2>/dev/null || true
    exit 1
  fi
  log "linking $PI_BIN -> /usr/local/bin/pi (sbx agent PATH probe)"
  sudo ln -sfn "$PI_BIN" /usr/local/bin/pi
else
  PI_BIN="$(command -v pi)"
  if [[ "$PI_BIN" != /usr/local/bin/pi ]]; then
    log "linking $PI_BIN -> /usr/local/bin/pi (sbx agent PATH probe)"
    sudo ln -sfn "$PI_BIN" /usr/local/bin/pi
  fi
fi

command -v pi
pi --version 2>/dev/null || true

mkdir -p "$SHARE"
touch "$MARKER"
log "done"
