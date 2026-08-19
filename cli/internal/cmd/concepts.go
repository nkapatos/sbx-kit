package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newConceptsCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "concepts",
		Short:   "How sbx-kit sits on sbx",
		Aliases: []string{"about"},
		Long:    conceptsText(),
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprint(cmd.OutOrStdout(), conceptsText())
		},
	}
}

func conceptsText() string {
	return `sbx-kit is a convenience layer on Docker sbx. It does not replace sbx.

Glossary
  catalog     path sbx-kit manages (sbx-kit setup; default ~/sbx-kit-catalog)
  directory   subdirectory with recipes/ (and optional kits/, images/)
  recipe      <dir>/<name> — sbx kind + kits + optional image
  kit         create-time YAML sbx applies at sandbox create
  image       custom Dockerfile or registry tag (sbx-kit recipes image)

Most commands need sbx-kit setup first, or a cwd inside your catalog.

Workflow
  sbx-kit setup
  sbx-kit catalog add <url>
  sbx-kit catalog ls
  sbx-kit catalog status
  sbx-kit catalog update
  sbx-kit recipes
  sbx-kit recipes verify
  sbx-kit recipes skill --cursor   # agent skill for sbx-kit
  sbx-kit box run <dir>/<name> --yes
  sbx-kit box bindings

Recipe checks: sbx-kit recipes verify. Kit checks: sbx (via sbx-kit recipes verify kits).
Kit schema: sbx and SPEC-v2 — not sbx-kit.
Parked recipe spec: sbx-kit experimental spec.

See also: sbx-kit --help
`
}
