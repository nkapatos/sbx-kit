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
		Short:         "Recipes, kits, and custom images on top of Docker sbx",
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

	root.AddCommand(newSetupCmd())
	root.AddCommand(newCatalogCmd())
	root.AddCommand(newConceptsCmd())
	root.AddCommand(newRecipesCmd())
	root.AddCommand(newRunCmd())
	root.AddCommand(newRmCmd())
	root.AddCommand(newUpgradeCmd())
	root.AddCommand(newStateCmd())
	root.AddCommand(newStatusCmd())
	root.AddCommand(newCheckCmd())
	root.AddCommand(newInitCmd())
	root.AddCommand(newImageCmd())
	root.AddCommand(newVersionCmd())

	return root
}

func longHelp() string {
	return `Convenience layer on Docker sbx: recipes (kind + kits + optional image),
portable state, and custom images (build locally or pull a registry, then
import into sbx).

  sbx-kit setup
  sbx-kit catalog add <git-url>
  sbx-kit catalog ls | fetch
  sbx-kit recipes
  sbx-kit run mine/cursor --yes
  sbx-kit run --name <sandbox>
  sbx-kit check | status
  sbx-kit image ls | load | pull

Stock recipes call sbx with no -t. Custom recipes pin an image (-t).
sbx template ls is the engine import store — use that, not image ls.

Host vault: ~/.local/share/sbx-kit/profiles/  and  ~/.local/state/sbx-kit/
Tree: sbx-kit setup  (override: SBX_KIT_TREE)
Recipe ids: <catalog>/<name>`
}

func requireToolkitRoot() (string, error) {
	root, err := toolkit.Root()
	if err != nil {
		return "", fmt.Errorf("%w\n  tip: sbx-kit setup", err)
	}
	return root, nil
}
