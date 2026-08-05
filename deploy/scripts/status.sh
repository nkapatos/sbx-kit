#!/usr/bin/env bash
# Show Pi stack health and a quick egress allow/deny probe (vps tier).
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
COMPOSE="$ROOT/deploy/scripts/compose.sh"

echo "==> compose ps"
"$COMPOSE" ps || true

echo
echo "==> pi: which pi / id"
if "$COMPOSE" ps -q pi >/dev/null 2>&1 && [[ -n "$("$COMPOSE" ps -q pi)" ]]; then
  "$COMPOSE" exec -T pi bash -lc 'id; command -v pi; pi --version 2>/dev/null || true' || true
else
  echo "pi container not running"
fi

echo
echo "==> egress probe (allow api.anthropic.com, deny example.com)"
proxy_cid="$("$COMPOSE" ps -q egress-proxy 2>/dev/null || true)"
if [[ -z "$proxy_cid" ]]; then
  echo "egress-proxy not running — start with: ./deploy/scripts/compose.sh up -d"
  exit 1
fi

proxy_ip="$(docker inspect -f '{{(index .NetworkSettings.Networks "sbx-kit-pi_agent").IPAddress}}' "$proxy_cid")"
net=sbx-kit-pi_agent

allow_code="$(docker run --rm --network "$net" \
  -e https_proxy="http://${proxy_ip}:8888" -e http_proxy="http://${proxy_ip}:8888" \
  curlimages/curl:8.5.0 -sS -o /dev/null -w '%{http_code}' --connect-timeout 15 \
  https://api.anthropic.com/ 2>/dev/null || echo fail)"

deny_ec=0
docker run --rm --network "$net" \
  -e https_proxy="http://${proxy_ip}:8888" -e http_proxy="http://${proxy_ip}:8888" \
  curlimages/curl:8.5.0 -sS -o /dev/null -w '%{http_code}' --connect-timeout 15 \
  https://example.com/ >/tmp/sbx-kit-deny.code 2>/tmp/sbx-kit-deny.err || deny_ec=$?

echo "  allow api.anthropic.com → HTTP ${allow_code} (expect 2xx/4xx from origin, not curl tunnel fail)"
echo "  deny  example.com       → curl_exit=${deny_ec} body=$(cat /tmp/sbx-kit-deny.code 2>/dev/null || true) (expect CONNECT fail / 403)"

if [[ "$deny_ec" -eq 0 ]]; then
  echo "WARNING: example.com was not blocked — check egress filter / overlay" >&2
  exit 2
fi
echo "OK: deny path working"
