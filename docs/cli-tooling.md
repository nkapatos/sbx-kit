# Toolkit CLI (`sbx-kit`)

Catalog-driven Go CLI for Docker AI Sandboxes templates, kits, and resource profiles.

**Install (macOS):** [homebrew.md](homebrew.md) — `brew tap nkapatos/sbx-kit …`  
**Go module:** `github.com/nkapatos/sbx-kit/cli`

## Commands

```text
sbx-kit agents
sbx-kit run <agent> [project-dir] [-- sbx-args...]
sbx-kit init [--agent <name>] [project-dir]
sbx-kit template load --engine <docker|container> <name-or-path> [image-tag]
sbx-kit version
```

## Catalog

Agents are declared in [`config/agents.yaml`](../config/agents.yaml). Defaults pull resource profile + kits; each agent sets the sbx agent name, image name, and kit list.

Resource profiles: [`config/resources-remote-llm.env`](../config/resources-remote-llm.env), [`config/resources-local-llm.env`](../config/resources-local-llm.env).

## Develop from a checkout

```bash
cd cli
go build -ldflags "-X github.com/nkapatos/sbx-kit/cli/internal/version.Version=dev" \
  -o ../bin/sbx-kit ./cmd/sbx-kit
export SBX_TREE=/path/to/sbx-kit
../bin/sbx-kit agents
../bin/sbx-kit run cursor .
../bin/sbx-kit template load --engine docker cursor-mise-docker
../bin/sbx-kit init --agent cursor /tmp/demo
```

Or: `go install github.com/nkapatos/sbx-kit/cli/cmd/sbx-kit@latest`
