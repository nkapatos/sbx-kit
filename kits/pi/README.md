# Pi kit (stub)

**Sandbox kit** (`kind: sandbox`) for Pi agents.

- **Image:** `local/sbx-shell-mise-docker:latest` (shared bake on `shell-docker`)
- **Mixins to add:** `mise-workspace` (and later `agent-workspace`)

No separate Pi Docker image — same machine as Hermes/shell; this kit owns entrypoint/install/auth.

Fill in step 4. See Docker’s [Build your own agent kit](https://docs.docker.com/ai/sandboxes/customize/build-an-agent/).
