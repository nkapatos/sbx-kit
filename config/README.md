# Resource defaults

sbx applies CPU, memory, and disk at **create** — not in the image.
`sbx-kit run` loads one of these files (`--resources` / `SBX_RESOURCES_PROFILE`).

| File | When |
| --- | --- |
| [`resources-remote-llm.env`](resources-remote-llm.env) | Default. External LLM APIs. |
| [`resources-local-llm.env`](resources-local-llm.env) | Local model needs more RAM/disk. |

Values live in the env files. Recreate the sandbox to pick up new defaults.
