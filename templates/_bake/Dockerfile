# Shared mise-prepared agent bake.
# Built via `sbx-kit template load` with a thin template's bake.env (BASE_IMAGE=…).
#
# Rules (do not violate):
#   - NEVER apt remove/purge/hold/install docker*, containerd*, runc, iptables*
#   - NEVER apt-get autoremove
#   - NEVER set ENV PATH (replaces base image PATH and breaks agent/DinD)
#   - Mise activate lives in kits/mise-workspace (not here) — binary only in this bake
#
ARG BASE_IMAGE
FROM ${BASE_IMAGE}

USER root

ARG YQ_VERSION=v4.45.4
ARG AST_GREP_VERSION=0.44.1
ARG TARGETARCH

# --- Remove language toolchains only (mise will own versions) ---
RUN set -eux; \
    export DEBIAN_FRONTEND=noninteractive; \
    apt-get update; \
    pkgs=''; \
    for p in \
      golang golang-go golang-src \
      nodejs npm yarn \
      openjdk-17-jdk openjdk-17-jre openjdk-21-jdk openjdk-21-jre default-jdk default-jre \
      python3-pip python3-venv python3-dev \
      rustc cargo \
    ; do \
      if dpkg -l "$p" 2>/dev/null | grep -q '^ii'; then pkgs="$pkgs $p"; fi; \
    done; \
    if [ -n "$pkgs" ]; then apt-get remove -y --purge $pkgs; fi; \
    rm -rf /var/lib/apt/lists/*

# Drop known language install trees only (not agent binaries, not Docker, not system python).
RUN set -eux; \
    rm -rf \
      /usr/local/go \
      /usr/lib/go \
      /usr/local/node \
      /usr/local/nodejs \
      /home/agent/.nvm \
      /home/agent/.cargo \
      /home/agent/.rustup \
      /home/agent/.local/share/mise \
      /root/.nvm \
      /root/.cargo \
      /root/.rustup \
    ; \
    rm -f \
      /usr/local/bin/go /usr/local/bin/gofmt \
      /usr/local/bin/node /usr/local/bin/npm /usr/local/bin/npx \
      /usr/bin/go /usr/bin/gofmt \
      2>/dev/null || true

# Agent CLIs + native build helpers (project language versions stay on mise).
RUN set -eux; \
    export DEBIAN_FRONTEND=noninteractive; \
    apt-get update; \
    apt-get install -y --no-install-recommends \
      build-essential \
      ca-certificates \
      cmake \
      curl \
      fd-find \
      git \
      git-lfs \
      iproute2 \
      locales \
      lsof \
      pkg-config \
      libssl-dev \
      unzip \
    ; \
    locale-gen en_US.UTF-8; \
    update-locale LANG=en_US.UTF-8; \
    ln -sf "$(command -v fdfind)" /usr/local/bin/fd; \
    arch="${TARGETARCH:-$(dpkg --print-architecture)}"; \
    case "$arch" in \
      amd64) yq_arch=amd64; sg_arch=x86_64 ;; \
      arm64) yq_arch=arm64; sg_arch=aarch64 ;; \
      *) echo "unsupported arch=$arch" >&2; exit 1 ;; \
    esac; \
    curl -fsSL "https://github.com/mikefarah/yq/releases/download/${YQ_VERSION}/yq_linux_${yq_arch}" \
      -o /usr/local/bin/yq; \
    chmod 755 /usr/local/bin/yq; \
    curl -fsSL "https://github.com/ast-grep/ast-grep/releases/download/${AST_GREP_VERSION}/app-${sg_arch}-unknown-linux-gnu.zip" \
      -o /tmp/ast-grep.zip; \
    unzip -o /tmp/ast-grep.zip -d /tmp/ast-grep; \
    install -m 755 /tmp/ast-grep/sg /usr/local/bin/sg; \
    install -m 755 /tmp/ast-grep/ast-grep /usr/local/bin/ast-grep; \
    rm -rf /tmp/ast-grep /tmp/ast-grep.zip; \
    git lfs install --system; \
    rm -rf /var/lib/apt/lists/*

# Mise binary + dirs — do NOT set PATH / activate here (kit owns activate).
ENV MISE_DATA_DIR=/mise \
    MISE_CONFIG_DIR=/mise \
    MISE_CACHE_DIR=/mise/cache \
    MISE_INSTALL_PATH=/usr/local/bin/mise \
    LANG=en_US.UTF-8 \
    LC_ALL=en_US.UTF-8 \
    PAGER=cat \
    GIT_PAGER=cat \
    EDITOR=true \
    VISUAL=true

RUN set -eux; \
    curl -fsSL https://mise.run | sh; \
    mkdir -p /mise/cache /mise/installs /mise/shims; \
    chown -R agent:agent /mise; \
    /usr/local/bin/mise --version

COPY gitconfig /etc/gitconfig
COPY agent-env.sh /etc/sbx-agent-env.sh

RUN set -eux; \
    chmod 644 /etc/sbx-agent-env.sh /etc/gitconfig; \
    touch /etc/sandbox-persistent.sh; \
    if ! grep -qF 'sbx-agent-env.sh' /etc/sandbox-persistent.sh; then \
      printf '\n# sbx-agent-ux\n[ -f /etc/sbx-agent-env.sh ] && . /etc/sbx-agent-env.sh\n' \
        >> /etc/sandbox-persistent.sh; \
    fi

USER agent
WORKDIR /home/agent
