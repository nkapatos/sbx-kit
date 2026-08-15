package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newConceptsCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "concepts",
		Short:   "How sbx and sbx-kit fit together",
		Aliases: []string{"about"},
		Long:    conceptsText(),
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprint(cmd.OutOrStdout(), conceptsText())
		},
	}
}

func conceptsText() string {
	return `sbx and sbx-kit use the same words for the same things.

  agent      (sbx)     What you run: shell, cursor, … Each boots from a template image.
  template   (sbx)     The image behind an agent (Hub/official, or one you load).
  kit        (sbx)     Create-time customization (network, credentials, startup, …).

  recipe     (sbx-kit) Named shortcut: which agent (+ optional template) + which kits.
                       So you do not pass kit paths by hand on every create.

sbx-kit sits on top of sbx: recipes, kit placement, portable state, and thin
helpers (template ls, check). It does not replace sbx secret, sbx run, or Hub pulls.

Typical flow:
  sbx-kit recipes              # catalog shortcuts
  sbx-kit agents               # what sbx can run + custom templates in view
  sbx-kit run --recipe …       # create from a recipe (or --name to attach)
  sbx-kit check                # binding + declared secrets + sbx secret ls
  sbx-kit status               # project ↔ sandbox bindings

See also: docs/cli-tooling.md
`
}
