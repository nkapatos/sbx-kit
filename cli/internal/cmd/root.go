package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/nkapatos/sbx-kit/cli/internal/toolkit"
	"github.com/nkapatos/sbx-kit/cli/internal/version"
)

// NewRoot builds the sbx-kit command tree.
func NewRoot() *cobra.Command {
	root := &cobra.Command{
		Use:           "sbx-kit",
		Short:         "Recipes and lifecycle helpers on top of Docker sbx",
		Long:          longHelp(),
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       version.Version,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	root.SetVersionTemplate("{{.Version}}\n")
	root.SetOut(os.Stdout)
	root.SetErr(os.Stderr)

	root.AddCommand(newConceptsCmd())
	root.AddCommand(newRecipesCmd())
	root.AddCommand(newAgentsCmd())
	root.AddCommand(newRunCmd())
	root.AddCommand(newRmCmd())
	root.AddCommand(newUpgradeCmd())
	root.AddCommand(newStateCmd())
	root.AddCommand(newStatusCmd())
	root.AddCommand(newCheckCmd())
	root.AddCommand(newInitCmd())
	root.AddCommand(newTemplateCmd())
	root.AddCommand(newVersionCmd())

	return root
}

func longHelp() string {
	return `Convenience layer on the Docker sbx CLI: recipes (agent + kits), lifecycle,
and portable state. Same words as sbx for agent, template, and kit.

  sbx-kit concepts     short wiring guide
  sbx-kit recipes      catalog shortcuts
  sbx-kit agents       sbx agents + custom templates in view
  sbx-kit run --recipe <id> --yes
  sbx-kit run --name <sandbox>
  sbx-kit check | status
  sbx-kit template ls | template load --engine docker <name>

  template load builds/imports unpublished images; recipes may also pin
  registry tags once images are published.

Host vault: ~/.local/share/sbx-kit/profiles/  and  ~/.local/state/sbx-kit/
Docs: docs/cli-tooling.md  ·  Override tree: SBX_TREE=`
}

func requireToolkitRoot() (string, error) {
	root, err := toolkit.Root()
	if err != nil {
		return "", fmt.Errorf("%w\n  tip: brew install places data in share/sbx-kit; for a git checkout set SBX_TREE", err)
	}
	return root, nil
}
