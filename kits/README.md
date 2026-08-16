# Kits

Create-time YAML (`spec.yaml` + optional `files/`). Not images.

This tree authors **schemaVersion `"1"`**. Upstream v2 exists (sbx ≥ 0.36); we
stay on v1 until `sbx-kit version` reports a floor ≥ 0.36 and every spec.yaml
is rewritten.

Which kits attach to which image: **`sbx-kit recipes`**. A sandbox kit **is**
the sbx kind (`name:`). Mixins stack. Credentials live on the kit that needs
them; `sbx-kit run` lists services — set any you use.

| Directory | Kind |
| --- | --- |
| [agent-workspace](agent-workspace/) | mixin (catalog default) |
| [mise-workspace](mise-workspace/) | mixin (custom images with mise) |
| [lsp-mise](lsp-mise/) | mixin (optional) |
| [apt-extras](apt-extras/) | mixin (optional) |
| [pi](pi/) | sandbox |

Add a kit: drop `kits/<name>/spec.yaml`, reference it from `config/agents.yaml`.
