# Homebrew tap (macOS)

Install **sbx-kit** by tapping this repo. No Homebrew-core submission and no separate brew org required — `Formula/sbx-kit.rb` lives in [nkapatos/sbx-kit](https://github.com/nkapatos/sbx-kit).

## Install

```bash
brew tap nkapatos/sbx-kit https://github.com/nkapatos/sbx-kit
brew install sbx-kit

sbx-kit version
sbx-kit agents
```

Before the first version tag, install from the default branch:

```bash
brew install --HEAD nkapatos/sbx-kit/sbx-kit
```

## What you get

| Path | Contents |
| --- | --- |
| `$(brew --prefix)/bin/sbx-kit` | CLI |
| `$(brew --prefix)/share/sbx-kit/` | `config/`, `kits/`, `templates/`, `docs/` |

Normal use does not need `SBX_TREE`. For hacking on this checkout:

```bash
export SBX_TREE=/path/to/sbx-kit
```

Lifecycle commands also use a **host vault** (created on demand, not by brew):

| Path | Contents |
| --- | --- |
| `~/.local/share/sbx-kit/profiles/` | Portable sandbox state archives |
| `~/.local/state/sbx-kit/` | Project ↔ sandbox name bindings |
## Host verify flow

Prerequisites: Docker **`sbx`** CLI signed in; Docker Desktop/Colima **or** Apple `container` + `skopeo`.

```bash
sbx-kit agents
sbx-kit template load --engine docker cursor-mise-docker
sbx template ls

cd ~/my-project
sbx-kit init --agent cursor .    # optional
sbx-kit run cursor .
```

Until images are published to a registry, `template load` is the local import path. Templates resolve from Brew `share/sbx-kit/templates` (or `SBX_TREE`).

## Update / remove

```bash
brew update && brew upgrade sbx-kit
brew uninstall sbx-kit
brew untap nkapatos/sbx-kit   # optional
```

## Cutting a release (maintainers)

1. Tag `v0.1.0` (or similar) on GitHub.
2. `curl -sL https://github.com/nkapatos/sbx-kit/archive/refs/tags/v0.1.0.tar.gz | shasum -a 256`
3. Set `url` + `sha256` in [`Formula/sbx-kit.rb`](../Formula/sbx-kit.rb).
4. Teammates run `brew upgrade sbx-kit`.

## Go (without brew)

Module: `github.com/nkapatos/sbx-kit/cli`

```bash
go install github.com/nkapatos/sbx-kit/cli/cmd/sbx-kit@latest
export SBX_TREE=/path/to/sbx-kit   # needed unless brew share is present
```
