# Agentic sandbox tooling

**Audience:** contributors and users of this toolkit.  
**Goal:** one **shared bake** (`templates/_bake`) + thin per-runtime images on official Docker sbx starters, plus runtime-agnostic mixin kits (`mise-workspace`, `agent-workspace`, optional capability kits).

This file is the source of truth for *what goes where*. Consumer READMEs are secondary.

---

## Architecture (do not couple kit to one template)

```text
docker/sandbox-templates:<runtime>-docker     # official starter (bake.env BASE_IMAGE)
        │
        ▼
   templates/_bake                            # CLIs, compilers, UX, mise, neovim (no activate)
        │
        ├── cursor-mise-docker
        ├── cursor-mise-ide                   # extends cursor-mise; IDE layer (scaffolding)
        ├── opencode-mise-docker
        └── shell-mise-docker                 # generic shell + parent for BYO agent layers
                │
                ├── pi-mise-docker            # Node + official Pi (agent in image)
                └── hermes-mise-docker        # Hermes CLI --skip-browser (agent in image)
                        │
                        ├── kits/mise-workspace
                        ├── kits/agent-workspace
                        ├── kits/lsp-mise | apt-extras   # optional
                        └── kits/hermes|pi               # thin: creds/network/context only
                                │
                                ▼
                        project mise.toml + .cursor/
                        /home/agent/.../portable/
```

| Concern | Lives in | Notes |
| --- | --- | --- |
| Agent runtime binary | Official starter **or** agent thin image (`pi-mise`, `hermes-mise`) | Do not heavy-install agents in kits |
| Always-on agent CLIs, compilers, non-interactive UX, **mise**, **neovim** | **Shared template bake** | Same package set for every runtime image; `EDITOR=true` stays |
| Cursor IDE (GUI) | **`cursor-mise-ide` thin image** | Extends cursor-mise; auth via sbx secrets |
| Per-runtime image tag | Thin `bake.env` → shared bake, or parent Dockerfile | One bake definition, N FROM lines |
| Mise data dirs, shims activate, registry allowlist | **`kits/mise-workspace` mixin** | Runtime-agnostic |
| Portable agent docs/refs + state pack/unpack | **`kits/agent-workspace` mixin** | `…/sbx-kit/portable/`; host vault via CLI |
| Project language versions | **Project `mise.toml`** | Never bake project pins into the image |
| Box LSPs / editor helpers | **`kits/lsp-mise`** (optional) | Global `/mise/config.toml` |
| Extra Linux packages | **`kits/apt-extras`** (optional) | apt in the VM — not host Homebrew |
| Cloud CLIs, Playwright, forks (OMP, …) | **Remote registry / `SBX_TREE`** | Not first-party examples — see [product-scope.md](product-scope.md) |
| Host nvim / `.cursor` / skills copy | **Optional kits** (not sbx-kit core) | Create-time personalization |

**Hard rules**

1. `mise-workspace` is **not** Cursor-specific. No cursor-only paths, entrypoints, or `aiFilename` assumptions in the kit.
2. Do **not** set `environment.variables.PATH` in the kit (breaks/fights sandbox PATH management). Activate + `/mise/shims` via persistent shell env only.
3. Do **not** run heavy agent installs in kit `commands.install` — bake agent binaries into templates. Mixins may use startup for idempotent workplace wiring; sandbox kits for BYO agents may install only if the agent is intentionally not imaged (avoid in this example tree).
4. Never put bash completion scripts in `/etc/sandbox-persistent.sh`.
5. Default recipes stay lean (agent image + mise + portable state). GUI IDE is opt-in (`cursor-mise-ide`). neovim binary in bake is OK; host editor **configs** stay kits.
6. Personalization (dotfiles, skills, forks) is kit/`SBX_TREE`/registry territory — not `sbx-kit` lifecycle commands.

---

## Official bases to support

Build thin images (or one parameterized build) for at least:

| Image tag (suggested) | `FROM` | Used with `sbx run …` |
| --- | --- | --- |
| `…/cursor-mise-docker` | `docker/sandbox-templates:cursor-agent-docker` | `cursor` |
| `…/opencode-mise-docker` | `docker/sandbox-templates:opencode-docker` (confirm exact tag in current Docker docs) | `opencode` |
| `…/shell-mise-docker` | `docker/sandbox-templates:shell-docker` | `shell` — **Pi agents, Hermes**, and other bring-your-own agents |

Prefer the **`-docker`** variants (Docker Engine inside the microVM) unless you explicitly want the lighter non-docker starters.

Confirm current official tags before shipping: [Templates](https://docs.docker.com/ai/sandboxes/customize/templates/), [Kits](https://docs.docker.com/ai/sandboxes/customize/kits/).

---

## Bake into every mise-prepared template (shared fragment)

Install as `root`, then return to `agent`. Pin release versions in build args.

### Packages / binaries

| Item | How | Why |
| --- | --- | --- |
| `fd` | apt `fd-find` + `ln -s fdfind /usr/local/bin/fd` | Fast file find from shell |
| `ast-grep` / `sg` | GitHub release zip (linux gnu, amd64/arm64) | Structural search/replace |
| mikefarah `yq` | GitHub release binary (not apt `yq` — that is a different tool) | YAML/JSON transforms |
| `build-essential`, `cmake`, `pkg-config` | apt | Native deps / wheels / crates |
| `iproute2` (`ss`), `lsof` | apt | Port / publish debugging |
| `git-lfs` | apt + `git lfs install --system` | LFS repos |
| `ca-certificates`, `locales` (`en_US.UTF-8`) | apt | TLS + UTF-8 |
| **`mise` binary** | Install if missing on base (curl install script or release); ensure `/usr/local/bin/mise` | Required so kit only configures, does not install mise |
| **`neovim`** | apt `neovim` | In-box editing and headless/ACP remote-dev; configs stay kits |
| `sqlite3`, `xz-utils` | apt | State WAL checkpoint; Hermes installer prereq |
| Keep from starters | `git`, `gh`, `jq`, `rg`, `curl`, `make`, `python3`, `uv`, `docker` (on `-docker` images) | Already strong on current Cursor sandbox |

### Non-interactive UX (bake)

Env (image `ENV` and/or `/etc/sbx-agent-env.sh` sourced from `/etc/sandbox-persistent.sh`):

```bash
export PAGER=cat
export GIT_PAGER=cat
export EDITOR=true
export VISUAL=true
export LANG=en_US.UTF-8
export LC_ALL=en_US.UTF-8
```

`EDITOR`/`VISUAL` stay non-interactive so agent/git flows never block; humans can still run `nvim` explicitly.

System gitconfig (`/etc/gitconfig`):

```ini
[core]
	pager = cat
	editor = true
[sequence]
	editor = true
[init]
	defaultBranch = main
[push]
	autoSetupRemote = true
[rebase]
	autoStash = true
```

### Explicitly do **not** bake

- Project language versions; asdf/nvm/sdkman; jj
- fzf, zoxide, starship, **GUI IDEs**, browsers (IDE → `cursor-mise-ide`; browsers → kits)
- kubectl/helm, aws/gcloud/azure CLIs, Playwright/Chromium, Terraform → separate kits
- Anything agent-runtime-specific (Cursor settings, OpenCode config, Pi/Hermes binaries) unless that thin image’s / sandbox kit’s job is to add that agent on `shell`
- Host editor configs / skills — optional kits only

---

## Portable docs / refs (agent-workspace)

| Place | Role |
| --- | --- |
| Repo (`AGENTS.md`, `docs/`) | Shared, reviewed team truth |
| `/home/agent/.local/share/sbx-kit/portable/` | Agent-first dumps/refs; survives vault export; empty aside from seeded README is OK |
| Gitignored `ref/` / `internal/` in project | Optional host scratch — not the official contract |

---

## Kit: `mise-workspace` (mixin, runtime-agnostic)

Reference implementation: [`kits/mise-workspace/spec.yaml`](../kits/mise-workspace/spec.yaml). Keep the contract below.

### Must provide

1. **`environment.variables`:** `MISE_DATA_DIR=/mise`, `MISE_CONFIG_DIR=/mise`, `MISE_CACHE_DIR=/mise/cache` only (no `PATH`).
2. **Startup (idempotent):** create `/mise/{cache,shims,installs}`; append activate block to `/etc/sandbox-persistent.sh` guarded by marker `# sbx-mise-workspace` (no completions).
3. **`network.allowedDomains`:** mise/aqua/GitHub releases + common registries (Go/npm/PyPI/crates/Docker Hub). Extend as real projects need; deny stays empty unless required.
4. **`agentContext`:** short imperative mise workflow (trust → install → prune; no apt toolchains; git+gh). Must not mention a single agent product as required.

### Must not provide

- Template image name / `sandbox:` block (this is a **mixin**, not a sandbox kit)
- `mise install` at create time
- Cursor-/OpenCode-/Pi-specific memory filenames (sbx writes kit context via the active sandbox kit’s `aiFilename`)

### Apply (examples)

```bash
sbx run --template <org>/cursor-mise-docker:tag --kit /path/to/mise-workspace cursor .
sbx run --template <org>/opencode-mise-docker:tag --kit /path/to/mise-workspace opencode .
sbx run --template <org>/shell-mise-docker:tag --kit /path/to/mise-workspace shell .
# Pi / Hermes: same shell-mise template + kit; install/launch agent via shell sandbox kit or entrypoint of your choosing
```

---

## Project contract (`mise.toml`)

- Pins only (example: `go = "1.26"`).
- First session: `mise trust` → `mise install` → `mise ls`.
- After removing pins: `mise install && mise prune -y`.
- Prefer fresh login shell if `GOROOT`/`GOBIN` look stale after upgrades.

---

## Resource defaults (create/run time)

sbx does not bake CPU/memory/disk into the image. `sbx-kit run` loads [`config/resources-remote-llm.env`](../config/resources-remote-llm.env) by default (agents calling remote LLM APIs). A follow-up profile lives at [`config/resources-local-llm.env`](../config/resources-local-llm.env).

| Setting | remote-llm | Override |
| --- | --- | --- |
| Memory | `4g` | `SBX_MEMORY` |
| CPUs | `4` | `SBX_CPUS` |
| Root FS | `10g` | `SBX_ROOT_SIZE` / `DOCKER_SANDBOXES_ROOT_SIZE` |
| Docker data | `25g` | `SBX_DOCKER_SIZE` / `DOCKER_SANDBOXES_DOCKER_SIZE` |

```bash
SBX_RESOURCES_PROFILE=local-llm sbx-kit run --agent cursor --yes   # follow-up
SBX_MEMORY=8g sbx-kit run --agent cursor --yes
```

Recreate sandboxes to pick up new disk/memory defaults.

## Git workspace modes

| Mode | Flag | Behavior |
| --- | --- | --- |
| Direct | *(default)* | Host tree R/W at the same path; live sync with the host editor |
| Clone | `--clone` | In-VM private clone; host at `/run/sandbox/source` RO; fetch via `sandbox-<name>` on the host |

**Defaults for this toolkit:** direct mode. Use `--clone` when running **multiple sandboxes** against the same repository or when you want a hard workspace boundary.

**Host worktrees** (`git worktree add <host-path>`) create extra checkouts on the **host** mount — useful so a human and an agent (or two host-visible trees) work in parallel. They are not sandbox-private storage. Do not use `git worktree add /tmp/…` inside the VM expecting isolation from the host `.git`.

```bash
sbx-kit run --agent cursor --clone
# or: sbx-kit run --agent cursor -- --clone
```

## Status

| Item | Status |
| --- | --- |
| `templates/_bake` + cursor thin image | Done (includes neovim) |
| `mise-workspace` (activate, prune, allowlists) | Done |
| Resource profiles (`remote-llm` / `local-llm`) | Done |
| Git workspace docs (direct / `--clone`) | Done |
| `sbx-kit` CLI + Homebrew tap | Done ([cli-tooling.md](cli-tooling.md), [homebrew.md](homebrew.md)) |
| `sbx-kit init` + `template load` | Done (bash `bin/` / `scripts/` removed) |
| Host XDG vault + sandbox name bindings | Done (`run` injects `--name`; `status` / `rm --keep-state` / `upgrade`) |
| `agent-workspace` portable state + `sbx-kit-state` | Done (seeded `portable/README`) |
| opencode / shell thin images | Done (`bake.env` + docs) |
| hermes, pi sandbox kits | Done (thin; agent binary in dedicated templates) |
| `pi-mise-docker`, `hermes-mise-docker` | Done (extend shell-mise) |
| Product scope doc | Done ([product-scope.md](product-scope.md)) |
| `lsp-mise`, `apt-extras` optional mixins | Done |
| `cursor-mise-ide` | Scaffold (IDE install TODO; recipe stub) |
| CI → Hub image publish + brew version tags | Planned |
| Remote recipe registries | Planned |
| Kit add/remove UX polish | Planned |
