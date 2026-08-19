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
  catalog    host path from sbx-kit setup (holds sources)
  source     subdirectory with recipes/, kits/, images/
  recipe     <source>/<name> — sbx kind + kits + optional image
  kit        create-time YAML sbx applies at sandbox create
  image      custom Dockerfile or registry tag (sbx-kit image, not sbx template ls)

sbx terms
  sbx run <kind>       sandbox kind plus optional -t image
  sbx template ls      images imported into the sbx engine
  sbx kit              create-time YAML at sandbox create

Workflow
  sbx-kit setup
  sbx-kit source add <git-url>
  sbx-kit source fetch
  sbx-kit recipes
  sbx-kit run <source>/<name> --yes
  sbx-kit image ls|load|pull
  sbx-kit check | status

See also: sbx-kit --help
`
}
