#!/usr/bin/env bash
# Copy compose/.env.example → compose/.env if missing; remind chmod 600.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
EXAMPLE="$ROOT/deploy/compose/.env.example"
DEST="$ROOT/deploy/compose/.env"

if [[ -f "$DEST" ]]; then
  echo "exists: $DEST (not overwriting)"
  ls -l "$DEST"
  exit 0
fi

cp "$EXAMPLE" "$DEST"
chmod 600 "$DEST"
echo "created $DEST (mode 600)"
echo "Edit WORKSPACE= and ANTHROPIC_API_KEY= before: ./deploy/scripts/compose.sh up -d"
