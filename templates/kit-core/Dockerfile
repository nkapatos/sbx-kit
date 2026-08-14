# syntax=docker/dockerfile:1
# sbx-kit minimal core — lean floor for Docker AI Sandboxes (shell family).
#
# Not based on docker/sandbox-templates:* (those ship language toolchains).
# Ships: agent user, sbx persistent-env glue, mise binary. No Node/Go/Java/Rust
# project toolchains — mise (and kits) own those.
#
# Tag: local/sbx-kit-core:latest
# Load: sbx-kit template load --engine docker kit-core
# Run:  sbx run --template local/sbx-kit-core:latest shell .
#   or: sbx-kit run --agent kit-core --yes

ARG DEBIAN_VERSION=bookworm-slim
FROM debian:${DEBIAN_VERSION}

ARG ENABLE_SUDO=1
ARG TARGETARCH

ENV DEBIAN_FRONTEND=noninteractive \
    LANG=en_US.UTF-8 \
    LC_ALL=en_US.UTF-8 \
    HOME=/home/agent \
    # Standard tool locations; do not fight sandbox PATH management later.
    PATH=/home/agent/.local/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin \
    NO_PROXY=localhost,127.0.0.1,::1,172.17.0.0/16,gateway.docker.internal \
    no_proxy=localhost,127.0.0.1,::1,172.17.0.0/16,gateway.docker.internal \
    MISE_DATA_DIR=/mise \
    MISE_CONFIG_DIR=/mise \
    MISE_CACHE_DIR=/mise/cache \
    MISE_INSTALL_PATH=/usr/local/bin/mise \
    PAGER=cat \
    GIT_PAGER=cat \
    EDITOR=true \
    VISUAL=true

# bookworm-slim has no CAs; cleartext apt is often blocked. Bootstrap CAs first.
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
      dnsutils \
      fd-find \
      git \
      git-lfs \
      iproute2 \
      jq \
      less \
      locales \
      lsof \
      make \
      procps \
      ripgrep \
      rsync \
      socat \
      sqlite3 \
      tini \
      unzip \
      xz-utils \
    ; \
    if [ "${ENABLE_SUDO}" = "1" ]; then \
      apt-get install -y --no-install-recommends sudo; \
    fi; \
    sed -i '/en_US.UTF-8/s/^# //g' /etc/locale.gen; \
    locale-gen en_US.UTF-8; \
    update-locale LANG=en_US.UTF-8 LC_ALL=en_US.UTF-8; \
    ln -sf "$(command -v fdfind)" /usr/local/bin/fd; \
    git lfs install --system; \
    rm -rf /var/lib/apt/lists/*

# Non-root workplace user (UID/GID 1000). Match official sandbox templates.
RUN set -eux; \
    groupadd --gid 1000 agent; \
    useradd --uid 1000 --gid 1000 --create-home --shell /bin/bash agent; \
    if [ "${ENABLE_SUDO}" = "1" ]; then \
      usermod -aG sudo agent; \
      echo 'agent ALL=(ALL) NOPASSWD:ALL' > /etc/sudoers.d/agent; \
      chmod 0440 /etc/sudoers.d/agent; \
      # Keep proxy-related env across sudo (agent + %sudo).
      printf '%s\n' \
        'Defaults:agent env_keep += "http_proxy https_proxy no_proxy HTTP_PROXY HTTPS_PROXY NO_PROXY SSL_CERT_FILE NODE_EXTRA_CA_CERTS REQUESTS_CA_BUNDLE JAVA_TOOL_OPTIONS"' \
        'Defaults:%sudo env_keep += "http_proxy https_proxy no_proxy HTTP_PROXY HTTPS_PROXY NO_PROXY SSL_CERT_FILE NODE_EXTRA_CA_CERTS REQUESTS_CA_BUNDLE JAVA_TOOL_OPTIONS"' \
        > /etc/sudoers.d/proxyconfig; \
      chmod 0440 /etc/sudoers.d/proxyconfig; \
    fi; \
    mkdir -p /home/agent/.local/bin /home/agent/.docker/sandbox/locks; \
    chown -R agent:agent /home/agent

# sbx sources this for every bash (interactive and BASH_ENV). No completions.
RUN set -eux; \
    touch /etc/sandbox-persistent.sh; \
    chmod 644 /etc/sandbox-persistent.sh; \
    chown agent:agent /etc/sandbox-persistent.sh; \
    printf '%s\n' \
      '# sbx-kit-core' \
      'umask 027' \
      'export LANG="${LANG:-en_US.UTF-8}"' \
      'export LC_ALL="${LC_ALL:-en_US.UTF-8}"' \
      > /etc/profile.d/sbx-kit-core.sh; \
    chmod 644 /etc/profile.d/sbx-kit-core.sh; \
    printf '%s\n' \
      'if [ -f /etc/sandbox-persistent.sh ]; then . /etc/sandbox-persistent.sh; fi' \
      'export BASH_ENV=/etc/sandbox-persistent.sh' \
      | tee /etc/profile.d/sandbox-persistent.sh >/dev/null; \
    chmod 644 /etc/profile.d/sandbox-persistent.sh; \
    # Also prepend into system bashrc so non-login interactive shells pick it up.
    if [ -f /etc/bash.bashrc ] && ! grep -qF 'sandbox-persistent.sh' /etc/bash.bashrc; then \
      printf '\n# sbx-kit-core\nif [ -f /etc/sandbox-persistent.sh ]; then . /etc/sandbox-persistent.sh; fi\nexport BASH_ENV=/etc/sandbox-persistent.sh\n' \
        >> /etc/bash.bashrc; \
    fi; \
    if ! grep -qF 'sandbox-persistent.sh' /home/agent/.bashrc 2>/dev/null; then \
      printf '%s\n' \
        'if [ -f /etc/sandbox-persistent.sh ]; then . /etc/sandbox-persistent.sh; fi' \
        'export BASH_ENV=/etc/sandbox-persistent.sh' \
        | cat - /home/agent/.bashrc > /tmp/agent-bashrc; \
      mv /tmp/agent-bashrc /home/agent/.bashrc; \
      chown agent:agent /home/agent/.bashrc; \
    fi; \
    printf 'core\n' > /etc/sbx-kit-agent; \
    chmod 644 /etc/sbx-kit-agent

ENV BASH_ENV=/etc/sandbox-persistent.sh

# Mise binary only — activate/shims belong to kits/mise-workspace.
RUN set -eux; \
    curl -fsSL https://mise.run | sh; \
    test -x /usr/local/bin/mise; \
    mkdir -p /mise/cache /mise/installs /mise/shims; \
    chown -R agent:agent /mise; \
    /usr/local/bin/mise --version

WORKDIR /home/agent
USER agent

ENTRYPOINT ["/usr/bin/tini", "--"]
CMD ["bash", "-l"]
