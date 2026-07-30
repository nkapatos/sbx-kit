# Extends the shared shell-mise floor (load shell-mise-docker first):
#   sbx-kit template load --engine docker shell-mise-docker
#   sbx-kit template load --engine docker hermes-mise-docker
#
# Agent binary is baked here; sbx wiring stays in kits/hermes.
# Browser/Playwright is skipped. Shell-family template warning is expected.

ARG PARENT_IMAGE=local/sbx-shell-mise-docker:latest
FROM ${PARENT_IMAGE}

USER root

# Official installer as root → FHS layout (/usr/local/bin/hermes).
RUN set -eux; \
    curl -fsSL https://hermes-agent.nousresearch.com/install.sh \
      | bash -s -- --skip-browser; \
    if [ ! -x /usr/local/bin/hermes ]; then \
      if [ -x /home/agent/.local/bin/hermes ]; then \
        install -m 755 /home/agent/.local/bin/hermes /usr/local/bin/hermes; \
      elif [ -x /root/.local/bin/hermes ]; then \
        install -m 755 /root/.local/bin/hermes /usr/local/bin/hermes; \
      else \
        echo "hermes binary not found after install" >&2; \
        find /home /root /usr/local -name 'hermes' -type f 2>/dev/null | head -50 >&2 || true; \
        exit 1; \
      fi; \
    fi; \
    test -x /usr/local/bin/hermes; \
    /usr/local/bin/hermes --version || true; \
    printf 'hermes\n' > /etc/sbx-kit-agent; \
    chmod 644 /etc/sbx-kit-agent

USER agent
WORKDIR /home/agent
