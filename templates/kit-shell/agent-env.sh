# Sourced from /etc/sandbox-persistent.sh — keep free of shell completions
# and secrets. UX only (matches official sandbox agent-env layout).

export PAGER="${PAGER:-cat}"
export GIT_PAGER="${GIT_PAGER:-cat}"
export EDITOR="${EDITOR:-true}"
export VISUAL="${VISUAL:-true}"
export LANG="${LANG:-en_US.UTF-8}"
export LC_ALL="${LC_ALL:-en_US.UTF-8}"
