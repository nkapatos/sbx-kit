# sbx-kit CLI

Go sources. See [docs/cli-tooling.md](../docs/cli-tooling.md) and `sbx-kit concepts`.

```bash
cd cli
go build -ldflags "-X github.com/nkapatos/sbx-kit/cli/internal/version.Version=dev" \
  -o ../bin/sbx-kit ./cmd/sbx-kit
export SBX_TREE=/path/to/this/repo
../bin/sbx-kit concepts
../bin/sbx-kit recipes
../bin/sbx-kit run shell --yes
../bin/sbx-kit check
```
