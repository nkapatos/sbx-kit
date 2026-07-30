# Extends the shared shell-mise floor (load shell-mise-docker first):
#   sbx-kit template load --engine docker shell-mise-docker
#   sbx-kit template load --engine docker pi-mise-docker
#
# Agent binary is baked here; sbx wiring (creds/network) stays in kits/pi.
# sbx treats this as a shell-family template — the "built for shell" warning
# when launching agent "pi" is expected; the sandbox kit supplies the runtime.

ARG PARENT_IMAGE=local/sbx-shell-mise-docker:latest
FROM ${PARENT_IMAGE}

USER root

ARG NODE_VERSION=22.16.0
ARG PI_PACKAGE=@earendil-works/pi-coding-agent
ARG TARGETARCH

# Node was stripped in the shared bake; Pi is an npm global. Force the CLI onto
# /usr/local/bin so sbx's agent-binary probe finds it (minimal PATH).
RUN set -eux; \
    export DEBIAN_FRONTEND=noninteractive; \
    arch="${TARGETARCH:-$(dpkg --print-architecture)}"; \
    case "$arch" in \
      amd64) node_arch=x64 ;; \
      arm64) node_arch=arm64 ;; \
      *) echo "unsupported arch=$arch" >&2; exit 1 ;; \
    esac; \
    curl -fsSL "https://nodejs.org/dist/v${NODE_VERSION}/node-v${NODE_VERSION}-linux-${node_arch}.tar.xz" \
      -o /tmp/node.tar.xz; \
    tar -xJf /tmp/node.tar.xz -C /usr/local --strip-components=1; \
    rm -f /tmp/node.tar.xz; \
    node --version; \
    npm --version; \
    npm install -g --ignore-scripts --prefix /usr/local "${PI_PACKAGE}"; \
    test -x /usr/local/bin/pi; \
    /usr/local/bin/pi --version; \
    # Keep a marker so we can tell host docker images from incomplete loads.
    printf 'pi\n' > /etc/sbx-kit-agent; \
    chmod 644 /etc/sbx-kit-agent

USER agent
WORKDIR /home/agent
