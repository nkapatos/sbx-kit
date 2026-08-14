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
		Short:         "Compose custom sbx templates/kits and manage portable sandbox state",
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
	root.AddCommand(newInitCmd())
	root.AddCommand(newTemplateCmd())
	root.AddCommand(newVersionCmd())

	return root
}

func longHelp() string {
	return `sbx-kit helps you compose custom Docker AI Sandboxes templates and kits,
give sandboxes stable identity, and export/restore workplace state across
teardown. This repo ships a few example recipes — point SBX_TREE (or a future
external tree) at your own templates/kits/catalog when you outgrow them.

macOS install:  brew tap nkapatos/sbx-kit https://github.com/nkapatos/sbx-kit && brew install sbx-kit
Docs:           docs/homebrew.md, docs/cli-tooling.md, docs/product-scope.md

Host vault (created on demand):
  ~/.local/share/sbx-kit/profiles/   portable state archives
  ~/.local/state/sbx-kit/            project↔recipe bindings

Day-to-day:
  sbx-kit agents
  sbx-kit run --agent cursor --yes
  sbx-kit run
  sbx-kit run --name my-project
  sbx-kit rm --agent cursor --keep-state
  sbx-kit upgrade --agent cursor
  sbx-kit status
  sbx-kit init --agent cursor .
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
