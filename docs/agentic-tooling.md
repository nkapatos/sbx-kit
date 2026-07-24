# Agentic sandbox tooling

**Audience:** contributors and users of this toolkit.  
**Goal:** one **shared bake** (`templates/_bake`) + thin per-runtime images on official Docker sbx starters, plus runtime-agnostic mixin kits (`mise-workspace`, later `agent-workspace`).

This file is the source of truth for *what goes where*. Consumer READMEs are secondary.

---

## Architecture (do not couple kit to one template)

```text
docker/sandbox-templates:<runtime>-docker     # official starter (bake.env BASE_IMAGE)
        │
        ▼
   templates/_bake                            # CLIs, compilers, UX, mise binary (no activate)
        │
        ├── cursor-mise-docker
        ├── opencode-mise-docker
        └── shell-mise-docker                 # Pi / Hermes / generic shell
                │
                ├── kits/mise-workspace       # mixin: dirs, activate, allowlists, agentContext
                ├── kits/agent-workspace      # mixin: portable state + sbx-kit-state + agentContext
                └── kits/hermes|pi            # sandbox kits on shell-mise (stubs)
                        │
                        ▼
                project mise.toml + .cursor/  # language pins + agent config
```

| Concern | Lives in | Notes |
| --- | --- | --- |
| Agent runtime binary / entrypoint | Official sbx starter (`cursor-agent`, `opencode`, `shell`, …) | Do not reinvent agents in the bake layer |
| Always-on agent CLIs, compilers, non-interactive UX, **mise binary** | **Shared template bake** (`agent-base` / Dockerfile fragment) | Same package set for every runtime image |
| Per-runtime image tag | Thin Dockerfile `FROM docker/sandbox-templates:<runtime>-docker` then apply shared bake | One bake definition, N FROM lines |
| Mise data dirs, shims activate, registry allowlist, mise agentContext | **`kits/mise-workspace` mixin** | Must work on cursor **and** opencode **and** shell without renaming |
| Portable agent state + `sbx-kit-state` pack/unpack | **`kits/agent-workspace` mixin** | Manifest/layout in kit (not bake); host CLI only `exec` + `cp` |
| Go/Node/Rust/Python versions, linters, `air`, … | **Project `mise.toml`** | Never bake project pins into the image |
| Cloud CLIs, Playwright, jj, k8s | **Optional thin kits** | Not part of mise-workspace |

**Hard rules**

1. `mise-workspace` is **not** Cursor-specific. No cursor-only paths, entrypoints, or `aiFilename` assumptions in the kit.
2. Do **not** set `environment.variables.PATH` in the kit (breaks/fights sandbox PATH management). Activate + `/mise/shims` via persistent shell env only.
3. Do **not** run `mise install` in kit `commands.install` (create-time WORKDIR is unreliable). Agent/`sbx exec` installs on first session.
4. Never put bash completion scripts in `/etc/sandbox-persistent.sh`.
5. Sandbox = agent workplace; host = human workplace (no fzf/IDE/browser in base).

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
- fzf, zoxide, starship, full editors, browsers
- kubectl/helm, aws/gcloud/azure CLIs, Playwright/Chromium, Terraform → separate kits
- Anything agent-runtime-specific (Cursor settings, OpenCode config, Pi/Hermes binaries) unless that thin image’s job is to add that agent on `shell`

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
SBX_RESOURCES_PROFILE=local-llm sbx-kit run cursor .   # follow-up
SBX_MEMORY=8g sbx-kit run cursor .
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
sbx-kit run cursor . -- --clone
# or: sbx-kit run cursor . --clone
```

## Status

| Item | Status |
| --- | --- |
| `templates/_bake` + cursor thin image | Done |
| `mise-workspace` (activate, prune, allowlists) | Done |
| Resource profiles (`remote-llm` / `local-llm`) | Done |
| Git workspace docs (direct / `--clone`) | Done |
| `sbx-kit` CLI + Homebrew tap | Done ([cli-tooling.md](cli-tooling.md), [homebrew.md](homebrew.md)) |
| `sbx-kit init` + `template load` | Done (bash `bin/` / `scripts/` removed) |
| Host XDG vault + sandbox name bindings | Done (`run` injects `--name`; `status` / `rm --keep-state` / `upgrade`) |
| `agent-workspace` portable state + `sbx-kit-state` | Done (kit-owned manifest/helper; bake floor optional later) |
| opencode / shell thin images | Stub (`bake.env` ready) |
| hermes, pi kits | Stub |
| Kit add/remove UX / CLI upgrades via kits | Planned |
