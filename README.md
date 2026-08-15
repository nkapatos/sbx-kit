# sbx-kit

Companion toolkit for [Docker AI Sandboxes](https://docs.docker.com/ai/sandboxes/) (`sbx`): **kits** (create-time YAML), **recipes**, optional **local templates**, and the **`sbx-kit`** CLI.

**Sandbox = agent workplace; host = human workplace.** Use this CLI to compose kits on top of **official Hub templates** (the supported Docker customization path) or on **local lean images** you build yourself. Recipes and portable state work for both.

Architecture: [docs/agentic-tooling.md](docs/agentic-tooling.md) · Scope: [docs/product-scope.md](docs/product-scope.md) · Homebrew: [docs/homebrew.md](docs/homebrew.md) · CLI: [docs/cli-tooling.md](docs/cli-tooling.md).

Official background: [Customize](https://docs.docker.com/ai/sandboxes/customize/) · [Templates](https://docs.docker.com/ai/sandboxes/customize/templates/).

## Concepts

| Term | Meaning |
| --- | --- |
| **Template** | Linux image the sandbox boots from — Hub/official via `sbx`, or a local build under `templates/`. |
| **Kit** | Directory with `spec.yaml` applied at sandbox create (mixin). Not an image. Official way to adjust any template. |
| **Recipe** | Named combo in `config/agents.yaml`: which `sbx` agent, which template source, which kits. |
| **CLI** | `sbx-kit` — recipes + kit placement + state migrate; `template load` only when you build locally |

## Two ways to get a template

| Path | When | What you do |
| --- | --- | --- |
| **Official / registry** | Day-to-day experiments, Docker’s supported model | Recipe with no local image → `sbx` uses the Hub agent template; kits layer on top |
| **Local build** | Lean floor (`kit-core` / `kit-cursor`) or your own Dockerfile | `sbx-kit template load --engine docker <name>` then run that recipe |

## Catalog

See [templates/README.md](templates/README.md) and [kits/README.md](kits/README.md).

| Status | Name | Notes |
| --- | --- | --- |
| Shipped | [`cli/`](cli/) + [`Formula/sbx-kit.rb`](Formula/sbx-kit.rb) | Toolkit CLI; macOS via Homebrew |
| Shipped | [`kits/agent-workspace`](kits/agent-workspace/), [`kits/mise-workspace`](kits/mise-workspace/) | State + mise mixins (mise needs a mise-ready image) |
| Shipped | [`kits/lsp-mise`](kits/lsp-mise/), [`kits/apt-extras`](kits/apt-extras/), [`kits/pi`](kits/pi/) | Optional mixins |
| Shipped | [`templates/kit-core`](templates/kit-core/), [`templates/kit-cursor`](templates/kit-cursor/) | Optional lean local floor (not required for Hub recipes) |
| Follow-up | Clearer Hub recipe UX, one-shot local `template load` for agent images, agent refresh vs upgrade | |

## Layout

```text
.
├── Formula/sbx-kit.rb           # Homebrew formula (macOS)
├── cli/                         # Go toolkit CLI (sbx-kit)
├── config/                      # agents.yaml + resource profiles
├── docs/                        # architecture, homebrew, CLI
├── kits/<name>/                 # mixins
└── templates/<name>/Dockerfile  # optional local images (kit-core, kit-cursor, …)
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

### 2a. Official template + kits (first path)

No local image build. `sbx` pulls/uses the stock agent template; kits attach at create:

```bash
cd ~/my-project
sbx-kit run --agent cursor-hub --yes    # stock cursor + agent-workspace
# or mix more kits via your own recipe in config/agents.yaml
```

### 2b. Local lean templates (optional)

```bash
sbx-kit template load --engine docker kit-core
sbx-kit template load --engine docker kit-cursor
# Apple container: --engine container (needs skopeo)
# Cursor package download: allow downloads.cursor.com if policy blocks it
sbx template ls
sbx-kit run --agent cursor --yes
```

### 3. Day-to-day

```bash
sbx-kit run                  # re-attach sole binding for cwd
sbx-kit run --name <id>      # attach by friendly sbx name
sbx-kit rm --agent cursor-hub --keep-state
sbx-kit upgrade --agent cursor-hub   # recreate from recipe + restore state
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
3. For a **local** image: add `templates/<name>/Dockerfile`, then `sbx-kit template load`.
4. Run: `sbx-kit run --agent <recipe> --yes`.

## Reference

### Import engines (local builds only)

| Host | Command |
| --- | --- |
| Docker Desktop / Colima | `sbx-kit template load --engine docker …` |
| Apple `container` | `sbx-kit template load --engine container …` (requires `skopeo`) |

Not supported: OrbStack, Podman. Pass `--engine` explicitly.

Apple `container image save` → OCI; `sbx template load` expects docker-archive. The container path converts with skopeo (`--override-os linux`).

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
| “Load complete” but missing from `sbx template ls` | Wrong tar format / engine — rerun `sbx-kit template load` |
| `skopeo` not found | `brew install skopeo` (Apple path) |
| Host has the image, sbx does not | Expected until import succeeds |
| Can’t find toolkit data | Brew share, or `export SBX_TREE=/path/to/checkout` |
| Wrong toolchain versions | Ensure `mise.toml` exists; agent or `sbx exec … mise install`; then `mise ls` |
| Removed pin still on PATH | `mise install && mise prune -y`; fresh `bash -l -c '…'` if env looks stale |
| Downloads blocked | `sbx policy log` / kit allowlist; Cursor package: `downloads.cursor.com` |
| `cursor-agent` missing after local load | Rebuild `kit-core` then `kit-cursor` |

## License

See repository license file if present; otherwise all rights reserved by the author until stated otherwise.
