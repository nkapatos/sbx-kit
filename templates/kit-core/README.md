# kit-core

Docker `FROM` parent (`debian:bookworm-slim`). **Not** imported into sbx.

OS, utils, sbx glue, mise **binary**. Languages → mise. Preference CLIs → kits.
Agent CLIs → child images. Docker Engine → later `-docker` / VPS.

Children: [kit-shell](../kit-shell/), [kit-cursor](../kit-cursor/).
Package list and layer cache: `Dockerfile`.
