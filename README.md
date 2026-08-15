# sbx-kit

Companion toolkit for [Docker AI Sandboxes](https://docs.docker.com/ai/sandboxes/) (`sbx`).
Recipes, kits, optional custom images, and the `sbx-kit` CLI.

**Sandbox = agent workplace; host = human workplace.**

```bash
brew tap nkapatos/sbx-kit https://github.com/nkapatos/sbx-kit
brew install sbx-kit          # or: brew install --HEAD nkapatos/sbx-kit/sbx-kit
sbx-kit version               # required sbx range
sbx-kit concepts
sbx-kit recipes
sbx-kit run cursor --yes
```

Install details: [docs/homebrew.md](docs/homebrew.md).
Scope: [docs/product-scope.md](docs/product-scope.md).
Official: [Customize](https://docs.docker.com/ai/sandboxes/customize/).

## Concepts

| Term | Meaning |
| --- | --- |
| **Kind** | First argument to `sbx run`. See `sbx run --help`. |
| **Template** | Image already in the sbx engine (`sbx template ls`). |
| **Kit** | Create-time YAML. Not an image. |
| **Recipe** | sbx-kit shortcut: kind + kits + optional custom image. |

Live catalog: `sbx-kit recipes`. Wiring: `sbx-kit concepts`. Commands: `sbx-kit --help`.

## Images

| Path | What |
| --- | --- |
| **Official** | Recipe with no image pin. `sbx` uses the Hub kind. |
| **Custom shell** | `kit-shell` — empty tinishell; add kits. |
| **Custom agent** | `FROM kit-core` (e.g. `kit-cursor`). Load the **child**. kit-core is never imported. |

```bash
sbx-kit image ls
sbx-kit image load --engine docker kit-shell   # also kit-cursor
sbx-kit run kit-shell --yes
```

`image load --help` covers engines. Stock recipes skip this.

## Layout

`cli/` (Go) · `config/` (recipes) · `kits/` · `templates/` · `Formula/sbx-kit.rb`

Homebrew ships the binary plus `share/sbx-kit/{config,kits,templates,docs}`.
Override with `SBX_TREE`. Host vault paths: `sbx-kit status --help`.
