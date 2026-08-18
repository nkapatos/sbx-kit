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

  sbx run <kind>     sandbox kind (cursor, shell, …) plus optional -t image
  sbx template ls    images already imported into the sbx engine
  sbx kit            create-time YAML sbx applies at sandbox create

  recipe             (sbx-kit) <catalog>/<name>: kind + kits + optional image
  catalog            (sbx-kit) one-level child of the tree (git clone or local dir)
  image              (sbx-kit) our Dockerfiles / registry tags — build or pull,
                     then import so they show up in sbx template ls

sbx's first argument is the sandbox kind. A recipe picks that for you.
Mixin kits stack on a Hub kind. A sandbox kit (pi) *is* the kind — the recipe's
sbx_agent matches the kit name, and the shell image lives in the kit (or a
recipe -t pin). Credentials live on the kit that needs them (many services
per kit); the host stores values with sbx secret set (any you use).

Catalog default kit: agent-workspace (portable state).

  sbx-kit setup
  sbx-kit catalog add <git-url>
  sbx-kit catalog fetch
  sbx-kit recipes
  sbx-kit run <catalog>/<name> --yes
  sbx-kit image ls|load|pull   custom images (not sbx template ls)
  sbx-kit check | status

See also: sbx-kit --help
`
}
