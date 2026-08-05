#!/usr/bin/env bash
# Thin wrapper around docker compose for the Pi VPS stack.
# Usage: ./deploy/scripts/compose.sh <compose-args...>
# Env: WORKSPACE (optional if set in compose/.env), PI_LOCKED=1, COMPOSE_PROJECT_NAME
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
COMPOSE_DIR="$ROOT/deploy/compose"
cd "$COMPOSE_DIR"

files=( -f pi.compose.yaml -f overlays/vps.yaml )
if [[ "${PI_LOCKED:-0}" == "1" ]]; then
  files+=( -f overlays/locked.yaml )
fi

# Load WORKSPACE / keys from .env for the wrapper process when present.
if [[ -f .env ]]; then
  set -a
  # shellcheck disable=SC1091
  source .env
  set +a
fi

if [[ -z "${WORKSPACE:-}" ]]; then
  echo "warning: WORKSPACE unset — compose defaults to ./workspace under compose/" >&2
fi

exec docker compose "${files[@]}" "$@"
