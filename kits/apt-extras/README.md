# Apt extras

Optional **mixin** for small apt packages inside the sandbox (Linux VM — use
apt here, not host Homebrew).

Default set is light (`tree`, `zip`, archive helpers, etc.). Use this kit (or a
team fork) for **preference CLIs** such as `gh` / `glab` — they are **not**
baked into `kit-shell`. Do **not** put Go/Node/Rust/Python toolchains here —
that is mise + project `mise.toml`.

```bash
# Catalog: add apt-extras to a recipe's kits list, or:
sbx kit add <sandbox> --kit /path/to/apt-extras
```

Playwright / browsers / cloud CLIs belong in dedicated kits, not this one.
