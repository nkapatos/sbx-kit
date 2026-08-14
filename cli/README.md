# sbx-kit CLI

Go sources for the toolkit CLI. See [docs/cli-tooling.md](../docs/cli-tooling.md).

**macOS teammates:** [docs/homebrew.md](../docs/homebrew.md).

```bash
cd cli
go build -ldflags "-X github.com/nkapatos/sbx-kit/cli/internal/version.Version=dev" \
  -o ../bin/sbx-kit ./cmd/sbx-kit
export SBX_TREE=/path/to/this/repo
../bin/sbx-kit agents
../bin/sbx-kit run --agent cursor --yes
../bin/sbx-kit run --agent cursor --yes --clone
../bin/sbx-kit template load --engine docker kit-core
../bin/sbx-kit template load --engine docker kit-cursor
../bin/sbx-kit init --agent cursor .
```

Brew installs the binary plus `share/sbx-kit/{config,kits,…}` so `SBX_TREE` is not required for normal use.
