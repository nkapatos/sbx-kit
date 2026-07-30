# Pi kit

Thin **sandbox kit** for the [Pi](https://pi.dev/) coding agent.

- **Image:** `local/sbx-pi-mise-docker:latest` (agent baked in the template)
- **Mixins:** `mise-workspace`, `agent-workspace` (via catalog)
- **This kit:** credentials, network, agentContext — **no** `commands.install`
- **Out of scope:** Oh My Pi (`omp`) and other forks (remote registry / `SBX_TREE`)

```bash
sbx-kit template load --engine docker shell-mise-docker
sbx-kit template load --engine docker pi-mise-docker
sbx-kit run --agent pi --yes
```

Register Anthropic (or other provider) secrets per Docker sbx docs. Iterate
missing domains with `sbx policy log`.
