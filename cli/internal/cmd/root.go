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
		Short:         "Compose Docker AI Sandboxes templates, kits, and agents",
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
	root.AddCommand(newInitCmd())
	root.AddCommand(newTemplateCmd())
	root.AddCommand(newVersionCmd())

	return root
}

func longHelp() string {
	return `sbx-kit launches coding agents in Docker AI Sandboxes using a catalog of
templates, mixin kits, and resource profiles.

macOS install:  brew tap nkapatos/sbx-kit https://github.com/nkapatos/sbx-kit && brew install sbx-kit
Docs:           docs/homebrew.md, docs/cli-tooling.md

Day-to-day:
  sbx-kit agents
  sbx-kit run cursor .
  sbx-kit run cursor . --clone
  sbx-kit init --agent cursor .
  sbx-kit template load --engine docker cursor-mise-docker`
}

func requireToolkitRoot() (string, error) {
	root, err := toolkit.Root()
	if err != nil {
		return "", fmt.Errorf("%w\n  tip: brew install places data in share/sbx-kit; for a git checkout set SBX_TREE", err)
	}
	return root, nil
}
