# syntax=docker/dockerfile:1
# Agent workplace floor for VPS / Compose (not an sbx template).
# Bare bootstrap only — languages and agent CLIs belong in agent layers or
# are installed at runtime by the agent / host.

ARG DEBIAN_VERSION=bookworm-slim
FROM debian:${DEBIAN_VERSION}

# Re-declare after FROM so build --build-arg ENABLE_SUDO=0 works in RUN layers.
ARG ENABLE_SUDO=1

ENV DEBIAN_FRONTEND=noninteractive \
    LANG=C.UTF-8 \
    LC_ALL=C.UTF-8

# Minimal bootstrap for HTTPS fetches + git clone + init.
# Deliberately no Node/Python/Go/Rust — agent or host owns those.
#
# bookworm-slim ships without ca-certificates. Cleartext apt (http) is often
# blocked by sandbox proxies; HTTPS needs CAs. Bootstrap once with peer
# verify disabled, install ca-certificates, then use normal HTTPS apt.
RUN set -eux; \
    if [ -f /etc/apt/sources.list.d/debian.sources ]; then \
      sed -i 's|http://|https://|g' /etc/apt/sources.list.d/debian.sources; \
    fi; \
    if [ -f /etc/apt/sources.list ]; then \
      sed -i 's|http://deb.debian.org|https://deb.debian.org|g; s|http://security.debian.org|https://security.debian.org|g' /etc/apt/sources.list; \
    fi; \
    apt-get -o Acquire::https::Verify-Peer=false \
            -o Acquire::https::Verify-Host=false update; \
    apt-get -o Acquire::https::Verify-Peer=false \
            -o Acquire::https::Verify-Host=false \
            install -y --no-install-recommends ca-certificates; \
    apt-get update; \
    apt-get install -y --no-install-recommends \
      bash \
      curl \
      git \
      tini \
      xz-utils \
    ; \
    if [ "${ENABLE_SUDO}" = "1" ]; then \
      apt-get install -y --no-install-recommends sudo; \
    fi; \
    rm -rf /var/lib/apt/lists/*

# Non-root workplace user (UID/GID 1000 matches common host bind mounts).
RUN set -eux; \
    groupadd --gid 1000 agent; \
    useradd --uid 1000 --gid 1000 --create-home --shell /bin/bash agent; \
    if [ "${ENABLE_SUDO}" = "1" ]; then \
      echo 'agent ALL=(ALL) NOPASSWD:ALL' > /etc/sudoers.d/agent; \
      chmod 0440 /etc/sudoers.d/agent; \
    fi

# Soften common footguns for agent shells without fighting Compose PATH mgmt.
RUN printf '%s\n' \
      '# sbx-kit deploy/base' \
      'umask 027' \
      'export LANG="${LANG:-C.UTF-8}"' \
      > /etc/profile.d/sbx-kit-agent.sh \
    && chmod 0644 /etc/profile.d/sbx-kit-agent.sh

USER agent
WORKDIR /home/agent
ENV HOME=/home/agent

# tini as PID 1 when Compose does not set an init.
ENTRYPOINT ["/usr/bin/tini", "--"]
CMD ["bash", "-l"]
