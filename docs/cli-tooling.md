# Toolkit CLI (`sbx-kit`)

Convenience layer on Docker `sbx`. Same words as sbx for **kind**, **template**,
and **kit**. **Recipe** is the only extra idea.

```bash
sbx-kit --help
sbx-kit concepts
sbx-kit recipes
sbx-kit version          # includes required sbx range
```

Install: [homebrew.md](homebrew.md). Scope: [product-scope.md](product-scope.md).

`--recipe` stays on lifecycle commands that have no positional
(`rm`, `upgrade`, `check`, `init`). Prefer `sbx-kit run cursor`.

Escape hatch for the sbx version gate: `SBX_KIT_SKIP_SBX_CHECK=1`.

## Develop from a checkout

```bash
cd cli
go build -ldflags "-X github.com/nkapatos/sbx-kit/cli/internal/version.Version=dev" \
  -o ../bin/sbx-kit ./cmd/sbx-kit
export SBX_TREE=/path/to/sbx-kit
../bin/sbx-kit concepts
```
