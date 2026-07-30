# LSP / editor helpers (mise)

Optional **mixin** for box-level LSPs and editor helpers via `/mise/config.toml`
(`commands.initFiles`, created only if missing).

- Does **not** replace the project `mise.toml` (product/CI pins stay there)
- Does **not** run `mise install` at create time (same rule as mise-workspace)
- Ships a commented starter config — uncomment tools or ship a team kit with pins

```bash
# Add to a recipe's kits: list in config/agents.yaml, or:
sbx kit add <sandbox> --kit /path/to/lsp-mise
```

Pair with `mise-workspace` (activate + allowlists) and, for humans in-box,
baked **neovim**.
