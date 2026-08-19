# Contributing

## Develop

CLI sources live in `cli/`. From that directory:

```bash
go test ./...
go vet ./...
```

`version.Version` is `dev` until CI or Homebrew sets it with `-ldflags`.

## Release

Do not bump `Formula/sbx-kit.rb` URLs by hand. Tag and publish from GitHub Actions:

1. Merge the work you want on the default branch.
2. Run the **Release** workflow (`workflow_dispatch`). Optional input: `0.2.0`.
3. Empty input lets git-cliff pick the next version from conventional commits since the last `v*.*.*` tag.
4. The script tags, writes the changelog, creates the GitHub Release, and pushes the Formula to the tap.

Host smoke before tagging (needs Docker `sbx` on PATH for box commands):

```bash
sbx-kit setup --catalog ~/sbx-kit-catalog
sbx-kit recipes create smoke
sbx-kit recipes verify --skip-kits smoke/shell
sbx-kit box run smoke/shell --yes
sbx-kit box check --recipe smoke/shell
sbx-kit box bindings
```

Confirm `sbx-kit version` reports an sbx that can check kits (`sbx kit validate`).

## Tap

The Homebrew formula in this repo is the source; the tap is updated by the release workflow (`HOMEBREW_TAP_TOKEN`).
