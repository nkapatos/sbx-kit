# Allowlisted egress proxy (tinyproxy).

Built from [`policy/pi.network.yaml`](../policy/pi.network.yaml) via
`../scripts/gen-egress-filter.sh`. Used by `compose/overlays/vps.yaml`.

```bash
./deploy/scripts/gen-egress-filter.sh
docker build -t local/sbx-kit-egress-proxy:latest -f deploy/egress/Dockerfile deploy/egress
```
