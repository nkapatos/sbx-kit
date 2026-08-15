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
		Short:         "Recipes and kits for Docker AI Sandboxes (Hub or local templates) + portable state",
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
	return `sbx-kit sits on top of the Docker sbx CLI: recipes, kit placement, and
portable agent state. Use official Hub templates (sbx pulls them) or local
images you build — same commands either way.

  Hub / official   recipe without a local image → sbx stock agent + your kits
  Local build      sbx-kit template load … then a recipe with image_name/fallback

This tree ships a few example recipes. Point SBX_TREE at your own
config/kits/templates when you want a different catalog.

macOS install:  brew tap nkapatos/sbx-kit https://github.com/nkapatos/sbx-kit && brew install sbx-kit
Docs:           docs/homebrew.md, docs/cli-tooling.md, docs/product-scope.md

Host vault (created on demand):
  ~/.local/share/sbx-kit/profiles/   portable state archives
  ~/.local/state/sbx-kit/            project↔recipe bindings

Day-to-day:
  sbx-kit agents
  sbx-kit template ls                     # → sbx template ls
  sbx-kit run --agent shell-hub --yes     # Hub shell + deepseek-creds trial
  sbx-kit check                           # diagnostics + sbx secret ls
  sbx-kit status                          # recipe↔sandbox bindings
  sbx-kit run --agent cursor-hub --yes
  sbx-kit run --agent cursor --yes        # after local template load
  sbx-kit run
  sbx-kit run --name my-project
  sbx-kit rm --agent shell-hub --keep-state
  sbx-kit upgrade --agent shell-hub
  sbx-kit init --agent shell-hub .
  sbx-kit template load --engine docker kit-core
  sbx-kit template load --engine docker kit-cursor`
}

func requireToolkitRoot() (string, error) {
	root, err := toolkit.Root()
	if err != nil {
		return "", fmt.Errorf("%w\n  tip: brew install places data in share/sbx-kit; for a git checkout set SBX_TREE", err)
	}
	return root, nil
}
