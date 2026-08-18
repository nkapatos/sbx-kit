package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/nkapatos/sbx-kit/cli/internal/toolkit"
)

func newSetupCmd() *cobra.Command {
	var tree string
	cmd := &cobra.Command{
		Use:   "setup",
		Short: "Point sbx-kit at a local recipes/kits/images tree",
		Long: `Writes the tree path to ~/.config/sbx-kit/config.yaml.

SBX_KIT_TREE overrides that path for one shell. The tree is always a local
directory (recipes/agents.yaml, kits/, images/).`,
		Example: `  sbx-kit setup --tree ~/src/sbx-kit-tree`,
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := tree
			if dir == "" {
				dir = os.Getenv(toolkit.TreeEnv)
			}
			if dir == "" {
				return fmt.Errorf("pass --tree <dir> (or set %s)", toolkit.TreeEnv)
			}
			abs, err := toolkit.WriteTree(dir)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "tree: %s\n", abs)
			fmt.Fprintf(cmd.OutOrStdout(), "optional: export %s=%s\n", toolkit.TreeEnv, abs)
			return nil
		},
	}
	cmd.Flags().StringVar(&tree, "tree", "", "path to the recipes/kits/images tree")
	return cmd
}
