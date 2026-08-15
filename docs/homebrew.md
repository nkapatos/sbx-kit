# Homebrew tap (macOS)

Install **sbx-kit** by tapping this repo. No Homebrew-core submission and no separate brew org required — `Formula/sbx-kit.rb` lives in [nkapatos/sbx-kit](https://github.com/nkapatos/sbx-kit).

## Install

```bash
brew tap nkapatos/sbx-kit https://github.com/nkapatos/sbx-kit
brew install sbx-kit

sbx-kit version
sbx-kit recipes
```

Before the first version tag, install from the default branch:

```bash
brew install --HEAD nkapatos/sbx-kit/sbx-kit
```

## What you get

| Path | Contents |
| --- | --- |
| `$(brew --prefix)/bin/sbx-kit` | CLI |
| `$(brew --prefix)/share/sbx-kit/` | Example `config/`, `kits/`, `templates/`, `docs/` |

This tap ships **example** recipes. Override with `SBX_TREE=/path/to/your-tree` to use your own templates/kits/catalog. Normal use of the brew share does not need `SBX_TREE`.

```bash
export SBX_TREE=/path/to/sbx-kit   # local checkout or another recipe tree
```

Lifecycle commands also use a **host vault** (created on demand, not by brew):

| Path | Contents |
| --- | --- |
| `~/.local/share/sbx-kit/profiles/` | Portable sandbox state archives |
| `~/.local/state/sbx-kit/` | Project ↔ recipe sandbox bindings |
## Host verify flow

Prerequisites: Docker **`sbx` CLI >= 0.34.0** signed in (kits = schemaVersion `"1"`); Docker Desktop/Colima **or** Apple `container` + `skopeo`.

```bash
sbx-kit version   # sbx-kit + required range + detected sbx
sbx-kit recipes

# Stock kind + kits (no custom image pin):
cd ~/my-project
sbx-kit init --recipe cursor .    # optional
sbx-kit run cursor --yes

# Custom images (optional):
sbx-kit image ls
sbx-kit image load --engine docker kit-shell
sbx-kit image load --engine docker kit-cursor
sbx template ls
sbx-kit run kit-cursor --yes
```

Override the version gate only if you know what you are doing: `SBX_KIT_SKIP_SBX_CHECK=1`.

`image load` imports **local** Dockerfiles into sbx. `image pull` does the same for a registry tag. Stock recipes skip both — `sbx` uses the Hub image for that kind. Images resolve from Brew `share/sbx-kit/templates` (or `SBX_TREE`). CI → registry publish and Formula version tags are a later release step.

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
