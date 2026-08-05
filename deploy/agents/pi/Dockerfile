# syntax=docker/dockerfile:1
# Pi coding agent on the deploy/base floor.

ARG BASE_IMAGE=local/sbx-kit-base:latest
FROM ${BASE_IMAGE}

ARG NODE_VERSION=22.16.0
ARG PI_PACKAGE=@earendil-works/pi-coding-agent
ARG TARGETARCH

USER root

# Node only in this agent layer (not the shared floor). Official tarball —
# no apt nodejs, so the floor stays language-free.
RUN set -eux; \
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
    npm install -g --ignore-scripts "${PI_PACKAGE}"; \
    test -x /usr/local/bin/pi; \
    pi --version; \
    printf 'pi\n' > /etc/sbx-kit-agent; \
    chmod 644 /etc/sbx-kit-agent

USER agent
WORKDIR /home/agent

ENTRYPOINT ["/usr/bin/tini", "--"]
CMD ["pi"]
