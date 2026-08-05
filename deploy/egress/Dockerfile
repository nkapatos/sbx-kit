# syntax=docker/dockerfile:1
# Allowlisted forward proxy for agent egress (vps / locked tiers).

ARG DEBIAN_VERSION=bookworm-slim
FROM debian:${DEBIAN_VERSION}

ENV DEBIAN_FRONTEND=noninteractive

RUN set -eux; \
    if [ -f /etc/apt/sources.list.d/debian.sources ]; then \
      sed -i 's|http://|https://|g' /etc/apt/sources.list.d/debian.sources; \
    fi; \
    apt-get -o Acquire::https::Verify-Peer=false \
            -o Acquire::https::Verify-Host=false update; \
    apt-get -o Acquire::https::Verify-Peer=false \
            -o Acquire::https::Verify-Host=false \
            install -y --no-install-recommends ca-certificates; \
    apt-get update; \
    apt-get install -y --no-install-recommends tinyproxy; \
    rm -rf /var/lib/apt/lists/*; \
    mkdir -p /etc/tinyproxy /var/log/tinyproxy; \
    chown -R nobody:nogroup /var/log/tinyproxy

COPY tinyproxy.conf /etc/tinyproxy/tinyproxy.conf
COPY filter /etc/tinyproxy/filter

USER nobody
EXPOSE 8888
ENTRYPOINT ["/usr/bin/tinyproxy", "-d", "-c", "/etc/tinyproxy/tinyproxy.conf"]
