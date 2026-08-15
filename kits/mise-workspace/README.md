# Mise workspace

Mixin for **custom** images that ship `/usr/local/bin/mise`. Not for official
Hub recipes (those already have toolchains). Does not run `mise install` or pin
languages — the project `mise.toml` owns that. Never sets
`environment.variables.PATH`.

Attach via a recipe (`sbx-kit recipes`), not by hand.
