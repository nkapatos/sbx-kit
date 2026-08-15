# Resource defaults

sbx applies CPU, memory, and disk at **sandbox create/run** time — not inside
the template image. `sbx-kit run` loads one of these files (see also
`SBX_RESOURCES_PROFILE` / `--resources`).

| File | When |
| --- | --- |
| [`resources-remote-llm.env`](resources-remote-llm.env) | Default. Agent talks to external LLM APIs. |
| [`resources-local-llm.env`](resources-local-llm.env) | Follow-up. Local model in/near the VM needs more RAM/disk. |

```bash
# default (remote)
sbx-kit run cursor --yes

# opt into local-LLM profile
sbx-kit run cursor --yes --resources local-llm
# or: SBX_RESOURCES_PROFILE=local-llm sbx-kit run cursor --yes

# one-off overrides
SBX_MEMORY=8g SBX_CPUS=6 sbx-kit run cursor --yes
```

Existing sandboxes keep the resources they were created with; recreate to pick up new defaults.
