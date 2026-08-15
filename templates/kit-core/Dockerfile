# syntax=docker/dockerfile:1
# sbx-kit-core — lean floor for local sbx templates (optional vs Hub agents).
#
# Why (not official docker/sandbox-templates:*):
#   Official images ship language toolchains we cannot maintain. We own a thin
#   floor; languages → mise; preference CLIs (gh/glab/…) → kits / in-box update;
#   agent binaries → thin layers with HOST-side refresh (not daily image rebuilds).
#
# Layer cache (rebuild only what changed):
#   1) CA bootstrap
#   2) Essential OS packages (rare: OS / major bumps; tini binary, not PID 1)
#   3) Modern utils (occasional: add/remove floor tools)
#   4) agent user + sudo / proxy env_keep
#   5) sbx persistent-env + UX
#   6) mise binary (occasional: mise releases)
#
# Not here: Node/Go/Java/Rust/Python, gh/glab, Docker Engine (future -docker variant).
# C toolchain is here so mise-managed langs can compile native addons.
#
# Tag: local/sbx-kit-core:latest (Docker FROM parent — not imported into sbx)
#   sbx-kit image load --engine docker kit-shell
#   sbx-kit image load --engine docker kit-cursor

ARG DEBIAN_VERSION=bookworm-slim
FROM debian:${DEBIAN_VERSION}

ARG ENABLE_SUDO=1
ARG TARGETARCH

ENV DEBIAN_FRONTEND=noninteractive \
    LANG=en_US.UTF-8 \
    LC_ALL=en_US.UTF-8 \
    HOME=/home/agent \
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
    VISUAL=true \
    TZ=UTC

# --- 1) CA bootstrap (bookworm-slim has none; cleartext apt often blocked) ---
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
    rm -rf /var/lib/apt/lists/*

# --- 2) Essential OS (rebuild rarely: debian / major floor changes) ---
RUN set -eux; \
    apt-get update; \
    apt-get install -y --no-install-recommends \
      bash \
      build-essential \
      bzip2 \
      curl \
      git \
      locales \
      openssh-client \
      patch \
      pkg-config \
      tini \
      tzdata \
      unzip \
      wget \
      xz-utils \
      zip \
      zstd \
    ; \
    if [ "${ENABLE_SUDO}" = "1" ]; then \
      apt-get install -y --no-install-recommends sudo; \
    fi; \
    sed -i '/en_US.UTF-8/s/^# //g' /etc/locale.gen; \
    locale-gen en_US.UTF-8; \
    update-locale LANG=en_US.UTF-8 LC_ALL=en_US.UTF-8; \
    rm -rf /var/lib/apt/lists/*

# --- 3) Modern utils (floor tools; preference CLIs stay in kits) ---
RUN set -eux; \
    apt-get update; \
    apt-get install -y --no-install-recommends \
      dnsutils \
      fd-find \
      file \
      git-lfs \
      iproute2 \
      iputils-ping \
      jq \
      less \
      lsof \
      openssl \
      procps \
      psmisc \
      ripgrep \
      rsync \
      socat \
      sqlite3 \
    ; \
    ln -sf "$(command -v fdfind)" /usr/local/bin/fd; \
    git lfs install --system; \
    rm -rf /var/lib/apt/lists/*

# --- 4) Workplace user (uid 1000) + sudo / proxy env_keep ---
RUN set -eux; \
    groupadd --gid 1000 agent; \
    useradd --uid 1000 --gid 1000 --create-home --shell /bin/bash agent; \
    if [ "${ENABLE_SUDO}" = "1" ]; then \
      usermod -aG sudo agent; \
      echo 'agent ALL=(ALL) NOPASSWD:ALL' > /etc/sudoers.d/agent; \
      chmod 0440 /etc/sudoers.d/agent; \
      printf '%s\n' \
        'Defaults:agent env_keep += "http_proxy https_proxy no_proxy HTTP_PROXY HTTPS_PROXY NO_PROXY SSL_CERT_FILE NODE_EXTRA_CA_CERTS REQUESTS_CA_BUNDLE JAVA_TOOL_OPTIONS"' \
        'Defaults:%sudo env_keep += "http_proxy https_proxy no_proxy HTTP_PROXY HTTPS_PROXY NO_PROXY SSL_CERT_FILE NODE_EXTRA_CA_CERTS REQUESTS_CA_BUNDLE JAVA_TOOL_OPTIONS"' \
        > /etc/sudoers.d/proxyconfig; \
      chmod 0440 /etc/sudoers.d/proxyconfig; \
    fi; \
    mkdir -p /home/agent/.local/bin /home/agent/.docker/sandbox/locks; \
    chown -R agent:agent /home/agent

# --- 5) sbx persistent-env + UX (no secrets, no completions) ---
COPY agent-env.sh /etc/sbx-agent-env.sh
RUN set -eux; \
    chmod 644 /etc/sbx-agent-env.sh; \
    touch /etc/sandbox-persistent.sh; \
    chmod 644 /etc/sandbox-persistent.sh; \
    chown agent:agent /etc/sandbox-persistent.sh; \
    if ! grep -qF 'sbx-agent-env.sh' /etc/sandbox-persistent.sh; then \
      printf '\n# sbx-agent-ux\n[ -f /etc/sbx-agent-env.sh ] && . /etc/sbx-agent-env.sh\n' \
        >> /etc/sandbox-persistent.sh; \
    fi; \
    printf '%s\n' \
      '# sbx-kit-core' \
      'umask 027' \
      > /etc/profile.d/sbx-kit-core.sh; \
    chmod 644 /etc/profile.d/sbx-kit-core.sh; \
    printf '%s\n' \
      'if [ -f /etc/sandbox-persistent.sh ]; then . /etc/sandbox-persistent.sh; fi' \
      'export BASH_ENV=/etc/sandbox-persistent.sh' \
      | tee /etc/profile.d/sandbox-persistent.sh >/dev/null; \
    chmod 644 /etc/profile.d/sandbox-persistent.sh; \
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

# --- 6) Mise binary only (activate/shims → kits/mise-workspace) ---
RUN set -eux; \
    curl -fsSL https://mise.run | sh; \
    test -x /usr/local/bin/mise; \
    mkdir -p /mise/cache /mise/installs /mise/shims; \
    chown -R agent:agent /mise; \
    /usr/local/bin/mise --version

WORKDIR /home/agent
USER agent

# Floor only — no tinishell. tini is installed for child images; kit-shell
# (and agent layers) declare PID 1. CMD is a fallback for docker run.
CMD ["bash", "-l"]
