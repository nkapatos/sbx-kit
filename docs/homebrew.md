# Homebrew tap (macOS)

```bash
brew tap nkapatos/sbx-kit https://github.com/nkapatos/sbx-kit
brew install sbx-kit
# before the first version tag:
brew install --HEAD nkapatos/sbx-kit/sbx-kit

sbx-kit version
sbx-kit concepts
sbx-kit recipes
```

| Path | Contents |
| --- | --- |
| `$(brew --prefix)/bin/sbx-kit` | CLI |
| `$(brew --prefix)/share/sbx-kit/` | Example `config/`, `kits/`, `templates/`, `docs/` |

Override the example tree with `SBX_TREE=/path/to/checkout`. Need Docker `sbx`
signed in; the required range is `sbx-kit version`.

```bash
brew update && brew upgrade sbx-kit
brew uninstall sbx-kit
```

Go without brew: `go install github.com/nkapatos/sbx-kit/cli/cmd/sbx-kit@latest`
and set `SBX_TREE`.

## Cutting a release

Tag, `shasum -a 256` the tarball, set `url` + `sha256` in `Formula/sbx-kit.rb`.
