# Hermes kit

Thin **sandbox kit** for the [Hermes](https://hermes-agent.nousresearch.com/) agent.

- **Image:** `local/sbx-hermes-mise-docker:latest` (agent baked; `--skip-browser`)
- **Mixins:** `mise-workspace`, `agent-workspace` (via catalog)
- **This kit:** network + agentContext — **no** `commands.install`

```bash
sbx-kit template load --engine docker shell-mise-docker
sbx-kit template load --engine docker hermes-mise-docker
sbx-kit run --agent hermes --yes
```

Provider secrets and first-run `hermes setup` / `hermes model` are host-side.
Watch `sbx policy log` for missing domains.
