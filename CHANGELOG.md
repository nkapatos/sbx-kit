# Changelog

All notable changes to this project are generated from conventional commits.

## [unreleased]

### Bug Fixes

- do not leave a GitHub tag when tap publish fails

## [0.1.1] - 2026-08-18

### CI

- run Actions on Node 24

## [0.1.0] - 2026-08-18

### Bug Fixes

- harden pi mixin install on kit-core
- install Pi with npm on Hub, mise node+pnpm on kit-core
- fill kit-core OS gaps agents hit on a lean floor

### CI

- add tag-synced release workflow and Go build checks
- test the release HEAD before tagging

### Documentation

- align CLI help and docs with recipe create/attach vision
- mark kit-core harden done; note kit-cursor rebuild
- custom templates are local build or registry pin
- kit-shell is the kit basis; kit-core is the image parent
- let CLI and recipes be the live reference
- keep the root README abstract
- split the user README from a contributor stub

### Features

- update the cli to manage agent sbx state files for teardown/migration/updates
- flag-first run UX with safer portable state export
- split run intents into create recipe vs attach by name
- use friendly sbx names with opaque profile ids in the vault
- add lean kit-core and kit-cursor templates
- harden kit-core floor for sbx agent workplaces
- apply kit-cursor box report to core, cursor, and state export
- Hub shell-hub trial with deepseek-creds and secret hints
- default agent-workspace; Pi sandbox kit on shell and kit-core
- Hub vs custom agentContext via floor.md probes
- point at a local tree and ship brew darwin binaries
- catalog trees, Homebrew tap, and sbx 0.38 floor

### Miscellaneous

- drop pi create-time install; keep thin pi creds mixin
- drop pi creds mixin; reserve baked agents for kit-core
- move session and product docs into the sandbox portable store
- add conventional commits and semver tag rules
- keep agent harness files out of the git tree
- move recipes, kits, and images out of this repo

### Other

- Initial public release of sbx-kit.
- checkpoint pi/hermes images, bake nvim, kits and sbxcompat
- park broken pi/hermes sbx recipes; keep shell-mise floor
- ai slope festival checkpoint

### Refactor

- make Pi a shell mixin; drop pi/hermes image templates
- drop official-base bake templates; retarget recipes to kit-core
- split kit-core layers; document host agent-refresh policy
- park VPS deploy; center Hub recipes in the CLI
- replace secrets command with check; soften create hints
- recipes vs agents, --recipe flag, and concepts help
- stock recipes match sbx agent names; kit- for custom
- drop --agent; run <recipe>; image ls/load/pull
- fold provider secrets into the Pi kit; drop deepseek-creds
- rename kit-core to kit-shell
- split kit-core floor from thin kit-shell
- kit-core is a FROM parent, not an sbx image


