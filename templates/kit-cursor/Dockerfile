# syntax=docker/dockerfile:1
# Cursor agent on sbx-kit-core (lean floor — no language toolchains).
#
# Load parent first:
#   sbx-kit template load --engine docker kit-core
#   sbx-kit template load --engine docker kit-cursor
#
# Tag: local/sbx-kit-cursor:latest
# Run: sbx-kit run --agent kit-cursor --yes
#   or: sbx run --template local/sbx-kit-cursor:latest cursor .
#
# Build needs network to downloads.cursor.com. On a locked sandbox host:
#   sbx policy allow network downloads.cursor.com

ARG PARENT_IMAGE=local/sbx-kit-core:latest
FROM ${PARENT_IMAGE}

# Pin when Cursor publishes a new lab build (see curl https://cursor.com/install).
ARG CURSOR_AGENT_VERSION=2026.08.11-e8db854
ARG TARGETARCH

# Official cursor-agent-docker preconfigures these so agent traffic works through
# the sbx host proxy (HTTP/1.1 + SSE). Mirror that on our lean floor.
ENV NODE_USE_ENV_PROXY=1 \
    NODE_EXTRA_CA_CERTS=/etc/ssl/certs/ca-certificates.crt

USER agent
WORKDIR /home/agent

# Same layout as https://cursor.com/install (agent + cursor-agent under ~/.local).
RUN set -eux; \
    arch="${TARGETARCH:-}"; \
    if [ -z "$arch" ]; then arch="$(uname -m)"; fi; \
    case "$arch" in \
      amd64|x86_64) arch=x64 ;; \
      arm64|aarch64) arch=arm64 ;; \
      *) echo "unsupported arch=$arch" >&2; exit 1 ;; \
    esac; \
    ver_dir="/home/agent/.local/share/cursor-agent/versions/${CURSOR_AGENT_VERSION}"; \
    tmp_dir="${ver_dir}.tmp-$$"; \
    mkdir -p "$tmp_dir" /home/agent/.local/bin; \
    curl -fsSL \
      "https://downloads.cursor.com/lab/${CURSOR_AGENT_VERSION}/linux/${arch}/agent-cli-package.tar.gz" \
      | tar --strip-components=1 -xzf - -C "$tmp_dir"; \
    mv "$tmp_dir" "$ver_dir"; \
    ln -sfn "$ver_dir/cursor-agent" /home/agent/.local/bin/cursor-agent; \
    ln -sfn "$ver_dir/cursor-agent" /home/agent/.local/bin/agent; \
    test -x /home/agent/.local/bin/cursor-agent; \
    /home/agent/.local/bin/cursor-agent --help >/dev/null

USER root
# sbx probes agent binaries with a minimal PATH; keep names on /usr/local/bin.
RUN set -eux; \
    ln -sfn /home/agent/.local/bin/cursor-agent /usr/local/bin/cursor-agent; \
    ln -sfn /home/agent/.local/bin/agent /usr/local/bin/agent; \
    printf 'cursor\n' > /etc/sbx-kit-agent; \
    chmod 644 /etc/sbx-kit-agent

USER agent
WORKDIR /home/agent
# sbx agent "cursor" supplies the real entrypoint (cursor-agent --yolo).
CMD ["bash", "-l"]
