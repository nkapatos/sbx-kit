# sbx-kit

Companion toolkit for [Docker AI Sandboxes](https://docs.docker.com/ai/sandboxes/) (`sbx`): reusable **templates** (sandbox images), **kits** (create-time YAML), and the **`sbx-kit`** CLI.

**Sandbox = agent workplace; host = human workplace** by default. First-party
images are a **lean core** (debian + sbx glue + mise) with thin agent layers —
not fat official `docker/sandbox-templates:*` bases. This repo ships a few
**example** recipes (core / cursor + mise mixins + portable state). Bring
your own templates/kits/catalog when you want a different stack.

Architecture: [docs/agentic-tooling.md](docs/agentic-tooling.md) · Scope: [docs/product-scope.md](docs/product-scope.md) · Homebrew: [docs/homebrew.md](docs/homebrew.md) · CLI: [docs/cli-tooling.md](docs/cli-tooling.md).

Official background: [Customize](https://docs.docker.com/ai/sandboxes/customize/) · [Templates](https://docs.docker.com/ai/sandboxes/customize/templates/).

## Concepts

| Term | Meaning |
| --- | --- |
| **Template** | Linux image the sandbox boots from (`templates/kit-core`, `templates/kit-cursor`, …). |
| **Kit** | Directory with `spec.yaml` applied at sandbox create (mixin or sandbox agent). Not an image. |
| **Import** | `sbx-kit template load` — build → save → load into sbx’s store (`docker` or Apple `container`). |
| **CLI** | `sbx-kit` — compose template + kits + resources; stamp project READMEs; migrate state |

Image tags: `local/sbx-<name>:latest` (e.g. `local/sbx-kit-cursor:latest`).

## Catalog

See [templates/README.md](templates/README.md) and [kits/README.md](kits/README.md).

| Status | Name | Notes |
| --- | --- | --- |
| Shipped | [`templates/kit-core`](templates/kit-core/) | Lean floor (cache-split); sbx + later VPS |
| Shipped | [`templates/kit-cursor`](templates/kit-cursor/) | Cursor bootstrap on kit-core (refresh agent on host) |
| Shipped | [`deploy/`](deploy/) | Docker/Compose VPS twin — see [`deploy/docs/vps-setup.md`](deploy/docs/vps-setup.md) |
| Shipped | [`cli/`](cli/) + [`Formula/sbx-kit.rb`](Formula/sbx-kit.rb) | Toolkit CLI; macOS via Homebrew |
| Shipped | [`kits/mise-workspace`](kits/mise-workspace/), [`kits/agent-workspace`](kits/agent-workspace/) | Default mixins |
| Shipped | [`kits/pi`](kits/pi/) | Thin DeepSeek mixin for future kit-pi image |
| Shipped | [`kits/lsp-mise`](kits/lsp-mise/), [`kits/apt-extras`](kits/apt-extras/) | Optional capability mixins |
| Follow-up | CI → Hub images, Compose export from CLI, SSH auth socket, more agent layers | |

## Layout

```text
.
├── Formula/sbx-kit.rb           # Homebrew formula (macOS)
├── cli/                         # Go toolkit CLI (sbx-kit)
├── config/                      # agents.yaml + resource profiles
├── deploy/                      # Docker/Compose VPS twin (converges with kit-core)
├── docs/                        # architecture, homebrew, CLI
├── kits/<name>/                 # mixins (and optional sandbox kits)
└── templates/<name>/Dockerfile  # kit-core, kit-cursor, …
```

Homebrew installs the binary plus `share/sbx-kit/{config,kits,templates,docs}`.

## Quick start (macOS)

### 1. Install

```bash
brew tap nkapatos/sbx-kit https://github.com/nkapatos/sbx-kit
brew install sbx-kit   # or: brew install --HEAD nkapatos/sbx-kit/sbx-kit
sbx-kit agents
```

You still need the Docker **`sbx` CLI >= 0.34.0** signed in (kits authored as
schemaVersion `"1"` until released sbx accepts v2). `sbx-kit version` reports
the required range. Details: [docs/homebrew.md](docs/homebrew.md).

### 2. Import templates (until Hub publishes images)

```bash
sbx-kit template load --engine docker kit-core
sbx-kit template load --engine docker kit-cursor
# Apple container: --engine container (needs skopeo)
# Cursor package download: allow downloads.cursor.com if policy blocks it
sbx template ls
```

### 3. Run a recipe in a project

```bash
cd ~/my-project
sbx-kit init --agent cursor .   # optional README stamp
sbx-kit run --agent cursor --yes
# later: sbx-kit run                  # re-attach sole binding
# or:    sbx-kit run --name <id>
# parallel / isolated: sbx-kit run --agent cursor --yes --clone
```

### Git workspace modes

| Mode | How | When to use |
| --- | --- | --- |
| **Direct** (default) | Host working tree mounted R/W at the same absolute path | Day-to-day; live host sync |
| **Clone** (`--clone`) | Private clone in the sandbox; host RO at `/run/sandbox/source` | Multiple sandboxes on the same repo / isolation |

Host **git worktrees** are for parallel host-visible checkouts. For VM-private agent trees use `--clone`, not `git worktree add /tmp/…`.

**Resource defaults** (`config/resources-remote-llm.env`): `--memory 4g --cpus 4`, root `10g`, Docker disk `25g`. See [config/](config/).

---

## Adding a template

1. Add `templates/<name>/Dockerfile` (usually `FROM local/sbx-kit-core:latest`).
2. Add a short `README.md`.
3. Import: `sbx-kit template load --engine <docker|container> <name>` (load `kit-core` first when needed).
4. Pair with mixin kits (`mise-workspace`, …).
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
  --template local/sbx-kit-cursor:latest \
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
| Downloads blocked | `sbx policy log` / kit allowlist; Cursor package: `downloads.cursor.com` |
| `cursor-agent` missing after load | Rebuild `kit-core` then `kit-cursor` |

## License

See repository license file if present; otherwise all rights reserved by the author until stated otherwise.
