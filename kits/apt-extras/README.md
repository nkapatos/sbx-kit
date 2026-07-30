# Apt extras

Optional **mixin** for small Ubuntu packages inside the sandbox (the VM is
Linux — use apt here, not host Homebrew formulae).

Default set is intentionally light (`tree`, `zip`, archives helpers, etc.).
Fork or replace this kit for team needs. Do **not** put Go/Node/Rust/Python
toolchains here — that is mise + project `mise.toml`.

```bash
# Catalog: add apt-extras to a recipe's kits list, or:
sbx kit add <sandbox> --kit /path/to/apt-extras
```

Playwright / browsers / cloud CLIs belong in dedicated kits, not this one.
