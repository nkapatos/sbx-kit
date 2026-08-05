#!/usr/bin/env bash
# Build deploy floor, Pi agent, and egress proxy images.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
BASE_TAG="${BASE_TAG:-local/sbx-kit-base:latest}"
PI_TAG="${PI_TAG:-local/sbx-kit-pi:latest}"
EGRESS_TAG="${EGRESS_TAG:-local/sbx-kit-egress-proxy:latest}"
ENABLE_SUDO="${ENABLE_SUDO:-1}"

"$ROOT/deploy/scripts/gen-egress-filter.sh"

echo "==> build $BASE_TAG (ENABLE_SUDO=$ENABLE_SUDO)"
docker build \
  -t "$BASE_TAG" \
  --build-arg "ENABLE_SUDO=$ENABLE_SUDO" \
  -f "$ROOT/deploy/base/Dockerfile" \
  "$ROOT/deploy/base"

echo "==> build $PI_TAG (FROM $BASE_TAG)"
docker build \
  -t "$PI_TAG" \
  --build-arg "BASE_IMAGE=$BASE_TAG" \
  -f "$ROOT/deploy/agents/pi/Dockerfile" \
  "$ROOT/deploy/agents/pi"

echo "==> build $EGRESS_TAG"
docker build -t "$EGRESS_TAG" -f "$ROOT/deploy/egress/Dockerfile" "$ROOT/deploy/egress"

echo "==> smoke"
docker run --rm --entrypoint bash "$BASE_TAG" -lc 'id; command -v curl git tini'
docker run --rm --entrypoint pi "$PI_TAG" --version

echo "Done."
echo "  laptop: docker compose -f deploy/compose/pi.compose.yaml run --rm --entrypoint pi pi --version"
echo "  vps:    WORKSPACE=/abs/project docker compose -f deploy/compose/pi.compose.yaml -f deploy/compose/overlays/vps.yaml up -d"
