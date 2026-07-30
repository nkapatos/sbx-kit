# Build from the shared cursor+mise floor, then layer Cursor IDE when packaging
# is confirmed. Load the parent image first:
#   sbx-kit template load --engine docker cursor-mise-docker
#
# IDE install details (Linux package / channel) are still iterating with UX —
# this template is the extension point; keep stub:true on the recipe until
# the install RUN is filled and verified.

ARG PARENT_IMAGE=local/sbx-cursor-mise-docker:latest
FROM ${PARENT_IMAGE}

USER root

# TODO: install Cursor IDE for Linux (in-box human workflow).
# Auth: use sbx secrets (same as cursor-agent); do not bake API keys.
# RUN set -eux; ...

USER agent
WORKDIR /home/agent
