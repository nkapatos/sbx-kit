# syntax=docker/dockerfile:1
# Custom counterpart to docker/sandbox-templates:shell.
# Floor is kit-core. This layer only: tini → login bash, and mark the kind
# as shell. No extra packages, no Node, no agent CLI.
#
#   sbx-kit image load --engine docker kit-core
#   sbx-kit image load --engine docker kit-shell
#   sbx-kit run kit-shell --yes
#
# Tag: local/sbx-kit-shell:latest

ARG PARENT_IMAGE=local/sbx-kit-core:latest
FROM ${PARENT_IMAGE}

USER root
RUN printf 'shell\n' > /etc/sbx-kit-agent && chmod 644 /etc/sbx-kit-agent

USER agent
WORKDIR /home/agent
ENTRYPOINT ["/usr/bin/tini", "--"]
CMD ["bash", "-l"]
