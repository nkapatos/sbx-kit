# Sourced from /etc/sandbox-persistent.sh — keep free of shell completions.
# Mise dirs/activate are owned by kits/mise-workspace; this file is UX only.
# neovim is baked for in-box / headless ACP use; keep EDITOR non-interactive
# so agent and git flows never block on an interactive editor.

export PAGER=cat
export GIT_PAGER=cat
export EDITOR=true
export VISUAL=true
export LANG="${LANG:-en_US.UTF-8}"
export LC_ALL="${LC_ALL:-en_US.UTF-8}"
