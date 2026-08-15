# sbx-kit

Companion toolkit for [Docker AI Sandboxes](https://docs.docker.com/ai/sandboxes/) (`sbx`): **kits**, **recipes**, optional **custom images** (build locally or pull from a registry, then import into sbx), and the **`sbx-kit`** CLI.

**Sandbox = agent workplace; host = human workplace.** Compose kits on **official Hub kinds** or on **custom images**. Recipes and portable state work for both.

Architecture: [docs/agentic-tooling.md](docs/agentic-tooling.md) · Scope: [docs/product-scope.md](docs/product-scope.md) · Homebrew: [docs/homebrew.md](docs/homebrew.md) · CLI: [docs/cli-tooling.md](docs/cli-tooling.md).

Official background: [Customize](https://docs.docker.com/ai/sandboxes/customize/) · [Templates](https://docs.docker.com/ai/sandboxes/customize/templates/).

## Concepts

| Term | Meaning |
| --- | --- |
| **Kind** | First argument to `sbx run` (`shell`, `cursor`, …). See `sbx run --help`. |
| **Template** | Image already imported into the sbx engine (`sbx template ls`). |
| **Kit** | Create-time customization (`spec.yaml`). Not an image. |
| **Recipe** | Named sbx-kit shortcut: kind + kits + optional custom image. |
| **CLI** | `sbx-kit` — recipes, placement, state; see `sbx-kit concepts`. |

## Two ways to get an image

| Path | When | What you do |
| --- | --- | --- |
| **Official / stock** | Day-to-day on Hub kinds | Recipe with no image pin → `sbx run <kind> --kit …` (no `-t`) |
| **Custom** | Lean floor (`kit-core` → `kit-shell` / `kit-cursor`) or your own image | `sbx-kit image load` while developing, or `sbx-kit image pull` for a registry tag; both import into sbx |

## Catalog

See [templates/README.md](templates/README.md) and [kits/README.md](kits/README.md).

| Status | Name | Notes |
| --- | --- | --- |
| Shipped | [`cli/`](cli/) + [`Formula/sbx-kit.rb`](Formula/sbx-kit.rb) | Toolkit CLI; macOS via Homebrew |
| Shipped | [`kits/agent-workspace`](kits/agent-workspace/), [`kits/mise-workspace`](kits/mise-workspace/) | State (default on every recipe) + mise mixins |
| Shipped | [`kits/lsp-mise`](kits/lsp-mise/), [`kits/apt-extras`](kits/apt-extras/) | Optional mixins |
| Shipped | [`kits/pi`](kits/pi/) | Sandbox kit; recipes `pi` (official shell) and `kit-pi` (kit-shell) |
| Shipped | [`templates/kit-core`](templates/kit-core/), [`templates/kit-shell`](templates/kit-shell/), [`templates/kit-cursor`](templates/kit-cursor/) | Floor, tiny shell, Cursor layer |
| Follow-up | Registry publish for templates; recipe/kit discovery; agent refresh vs upgrade | |

## Layout

```text
.
├── Formula/sbx-kit.rb           # Homebrew formula (macOS)
├── cli/                         # Go toolkit CLI (sbx-kit)
├── config/                      # agents.yaml + resource profiles
├── docs/                        # architecture, homebrew, CLI
├── kits/<name>/                 # mixins
└── templates/<name>/Dockerfile  # custom images (build/load; publish to registry)
```

Homebrew installs the binary plus `share/sbx-kit/{config,kits,templates,docs}`.

## Quick start (macOS)

### 1. Install

```bash
brew tap nkapatos/sbx-kit https://github.com/nkapatos/sbx-kit
brew install sbx-kit   # or: brew install --HEAD nkapatos/sbx-kit/sbx-kit
sbx-kit recipes
sbx-kit concepts
```

You still need the Docker **`sbx` CLI >= 0.34.0** signed in (kits authored as
schemaVersion `"1"` until released sbx accepts v2). `sbx-kit version` reports
the required range. Details: [docs/homebrew.md](docs/homebrew.md).

### 2a. Official kind + kits (first path)

No local image build. `sbx` uses the stock Hub image; kits attach at create:

```bash
cd ~/my-project
sbx-kit concepts
sbx-kit recipes
sbx-kit run shell --yes
sbx-kit check
```

### 2b. Custom images (optional)

```bash
sbx-kit image ls
sbx-kit image load --engine docker kit-core
sbx-kit image load --engine docker kit-shell
sbx-kit image load --engine docker kit-cursor
# or: sbx-kit image pull ghcr.io/example/sbx-kit-cursor:latest
# Apple container: --engine container (needs skopeo)
# Cursor package download: allow downloads.cursor.com if policy blocks it
sbx template ls
sbx-kit run kit-cursor --yes
```

### 3. Day-to-day

```bash
sbx-kit run                  # re-attach sole binding for cwd
sbx-kit run --name <id>      # attach by friendly sbx name
sbx-kit rm --recipe cursor --keep-state
sbx-kit upgrade --recipe cursor
```

### Git workspace modes

| Mode | How | When to use |
| --- | --- | --- |
| **Direct** (default) | Host working tree mounted R/W at the same absolute path | Day-to-day; live host sync |
| **Clone** (`--clone`) | Private clone in the sandbox; host RO at `/run/sandbox/source` | Multiple sandboxes on the same repo / isolation |

Host **git worktrees** are for parallel host-visible checkouts. For VM-private agent trees use `--clone`, not `git worktree add /tmp/…`.

**Resource defaults** (`config/resources-remote-llm.env`): `--memory 4g --cpus 4`, root `10g`, Docker disk `25g`. See [config/](config/).

---

## Adding a kit or recipe

1. Add `kits/<name>/spec.yaml` (and optional `files/`).
2. Reference it from [`config/agents.yaml`](config/agents.yaml) on a Hub or local recipe.
3. For a **local** image: add `templates/<name>/Dockerfile`, then `sbx-kit image load`.
4. Run: `sbx-kit run <recipe> --yes`.

## Reference

### Import engines (`image load`)

| Host | Command |
| --- | --- |
| Docker Desktop / Colima | `sbx-kit image load --engine docker …` |
| Apple `container` | `sbx-kit image load --engine container …` (requires `skopeo`) |

Not supported: OrbStack, Podman. Pass `--engine` explicitly.

Apple `container image save` → OCI; `sbx template load` expects docker-archive. The container path converts with skopeo (`--override-os linux`).

Registry tags: `sbx-kit image pull <registry/tag>` (docker pull, then the same import).

### Long form (no sbx-kit)

```bash
# Official agent + kit paths:
sbx run cursor --kit "$(brew --prefix)/share/sbx-kit/kits/agent-workspace" .

# Local template:
sbx run cursor \
  --template local/sbx-kit-cursor:latest \
  --kit "$(brew --prefix)/share/sbx-kit/kits/mise-workspace" \
  .
```

### Troubleshooting

| Symptom | Fix |
| --- | --- |
| “Load complete” but missing from `sbx template ls` | Wrong tar format / engine — rerun `sbx-kit image load` |
| `skopeo` not found | `brew install skopeo` (Apple path) |
| Host has the image, sbx does not | Expected until import succeeds |
| Can’t find toolkit data | Brew share, or `export SBX_TREE=/path/to/checkout` |
| Wrong toolchain versions | Ensure `mise.toml` exists; agent or `sbx exec … mise install`; then `mise ls` |
| Removed pin still on PATH | `mise install && mise prune -y`; fresh `bash -l -c '…'` if env looks stale |
| Downloads blocked | `sbx policy log` / kit allowlist; Cursor package: `downloads.cursor.com` |
| `cursor-agent` missing after local load | Rebuild `kit-core` then `kit-cursor` |

## License

See repository license file if present; otherwise all rights reserved by the author until stated otherwise.
