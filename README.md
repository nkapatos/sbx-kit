# sbx-kit

Companion toolkit for [Docker AI Sandboxes](https://docs.docker.com/ai/sandboxes/) (`sbx`): reusable **templates** (sandbox images), **kits** (create-time YAML), and the **`sbx-kit`** CLI.

**Sandbox = agent workplace; host = human workplace.** Heavy tooling is baked into shared templates; kits are mixins (or BYO sandbox agents). There is no Cursor-specific mise kit — only thin per-agent images plus agnostic mixins.

Architecture: [docs/agentic-tooling.md](docs/agentic-tooling.md) · Homebrew: [docs/homebrew.md](docs/homebrew.md) · CLI: [docs/cli-tooling.md](docs/cli-tooling.md).

Official background: [Customize](https://docs.docker.com/ai/sandboxes/customize/) · [Templates](https://docs.docker.com/ai/sandboxes/customize/templates/).

## Concepts

| Term | Meaning |
| --- | --- |
| **Template** | Linux image the sandbox boots from. Thin dirs set `bake.env`; shared body is `templates/_bake`. |
| **Kit** | Directory with `spec.yaml` applied at sandbox create (mixin or sandbox agent). Not an image. |
| **Import** | `sbx-kit template load` — build → save → load into sbx’s store (`docker` or Apple `container`). |
| **CLI** | `sbx-kit` — compose template + kits + resources; stamp project READMEs |

Image tags follow: `local/sbx-<role>-<capability>[-<runtime>]:<tag>`  
(e.g. `local/sbx-cursor-mise-docker:latest`).

## Catalog

See [templates/README.md](templates/README.md) and [kits/README.md](kits/README.md).

| Status | Name | Notes |
| --- | --- | --- |
| Shipped | [`templates/_bake`](templates/_bake/) | Shared agent bake (strip languages, CLIs, mise binary, UX) |
| Shipped | [`templates/cursor-mise-docker`](templates/cursor-mise-docker/) | Thin: cursor-agent-docker + bake; `sbx-kit run cursor` |
| Shipped | [`cli/`](cli/) + [`Formula/sbx-kit.rb`](Formula/sbx-kit.rb) | Toolkit CLI; macOS via Homebrew |
| Stub | [`templates/opencode-mise-docker`](templates/opencode-mise-docker/) | Same bake on opencode-docker |
| Stub | [`templates/shell-mise-docker`](templates/shell-mise-docker/) | Same bake on shell-docker (Pi/Hermes image) |
| Stub | [`kits/hermes`](kits/hermes/), [`kits/pi`](kits/pi/) | Sandbox kits on shell-mise (entrypoint TBD) |
| Shipped | [`kits/agent-workspace`](kits/agent-workspace/) | Portable state + `sbx-kit-state`; MCP/skills still TBD |
| Planned | `cursor-mise-ide` | Same stack with Cursor IDE baked in |
| Planned | `nvim-*` | ACP / headless nvim remote-dev |

## Layout

```text
.
├── Formula/sbx-kit.rb           # Homebrew formula (macOS)
├── cli/                         # Go toolkit CLI (sbx-kit)
├── config/                      # agents.yaml + resource profiles
├── docs/                        # architecture, homebrew, CLI
├── kits/<name>/                 # mixins and sandbox agent kits
├── templates/_bake/             # shared Dockerfile + UX files
└── templates/<name>/bake.env    # BASE_IMAGE=… for thin images
```

Homebrew installs the binary plus `share/sbx-kit/{config,kits,templates,docs}`.

## Quick start (macOS)

### 1. Install

```bash
brew tap nkapatos/sbx-kit https://github.com/nkapatos/sbx-kit
brew install sbx-kit   # or: brew install --HEAD nkapatos/sbx-kit/sbx-kit
sbx-kit agents
```

You still need the Docker **`sbx`** CLI signed in. Details: [docs/homebrew.md](docs/homebrew.md).

### 2. Import a template (until Hub publishes images)

```bash
sbx-kit template load --engine docker cursor-mise-docker
# Apple container: --engine container (needs skopeo)
sbx template ls
```

### 3. Run an agent in a project

```bash
cd ~/my-project
sbx-kit init --agent cursor .   # optional README stamp
sbx-kit run cursor .
# parallel / isolated: sbx-kit run cursor . --clone
```

### Git workspace modes

| Mode | How | When to use |
| --- | --- | --- |
| **Direct** (default) | Host working tree mounted R/W at the same absolute path | Day-to-day; live host sync |
| **Clone** (`--clone`) | Private clone in the sandbox; host RO at `/run/sandbox/source` | Multiple sandboxes on the same repo / isolation |

Host **git worktrees** are for parallel host-visible checkouts. For VM-private agent trees use `--clone`, not `git worktree add /tmp/…`.

**Resource defaults** (`config/resources-remote-llm.env`): `--memory 4g --cpus 4`, root `10g`, Docker disk `25g`. See [config/](config/).

---

## Adding a thin template

1. Add `templates/<name>/bake.env` with `BASE_IMAGE=docker/sandbox-templates:<official>-docker`.
2. Add a short `README.md`; do **not** copy `_bake/Dockerfile`.
3. Import: `sbx-kit template load --engine <docker|container> <name>`.
4. Pair with mixin kits (`mise-workspace`, …). Sandbox kits only for BYO agents (Pi/Hermes).
5. Add an entry in [`config/agents.yaml`](config/agents.yaml) and [templates/README.md](templates/README.md).

## Reference

### Import engines

| Host | Command |
| --- | --- |
| Docker Desktop / Colima | `sbx-kit template load --engine docker …` |
| Apple `container` | `sbx-kit template load --engine container …` (requires `skopeo`) |

Not supported: OrbStack, Podman. Pass `--engine` explicitly.

Apple `container image save` → OCI; `sbx template load` expects docker-archive. The container path converts with skopeo (`--override-os linux`).

### Long form (no sbx-kit)

```bash
sbx run cursor \
  --template docker.io/local/sbx-cursor-mise-docker:latest \
  --kit "$(brew --prefix)/share/sbx-kit/kits/mise-workspace" \
  .
```

### Troubleshooting

| Symptom | Fix |
| --- | --- |
| “Load complete” but missing from `sbx template ls` | Wrong tar format / engine — rerun `sbx-kit template load` |
| `skopeo` not found | `brew install skopeo` (Apple path) |
| Host has the image, sbx does not | Expected until import succeeds |
| Can’t find toolkit data | Brew share, or `export SBX_TREE=/path/to/checkout` |
| Wrong toolchain versions | Ensure `mise.toml` exists; agent or `sbx exec … mise install`; then `mise ls` |
| Removed pin still on PATH | `mise install && mise prune -y`; fresh `bash -l -c '…'` if env looks stale |
| Downloads blocked | `sbx policy log` / kit allowlist |

```bash
sbx template ls
sbx ls
sbx exec -it <name> bash
```
