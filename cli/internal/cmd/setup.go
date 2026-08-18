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
		Short: "Configure the local recipes tree",
		Long: `Writes the tree path to ~/.config/sbx-kit/config.yaml.

The tree is a directory of catalogs (one-level children). Each child is a
catalog: recipes/agents.yaml, kits/, images/. Git is optional.

With no --tree, setup asks for the path (default ~/sbx-kit-recipes).
SBX_KIT_TREE overrides the configured path for one shell.`,
		Example: `  sbx-kit setup
  sbx-kit setup --tree ~/sbx-kit-recipes`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := tree
			if dir == "" {
				dir = os.Getenv(toolkit.TreeEnv)
			}
			if dir == "" {
				if !stdinIsTTY() {
					return fmt.Errorf("pass --tree <dir> (or run setup from a terminal)")
				}
				def := toolkit.DefaultTree()
				if existing, err := toolkit.ConfiguredTree(); err == nil && existing != "" {
					def = existing
				}
				var err error
				dir, err = promptLine(cmd.InOrStdin(), cmd.OutOrStdout(), "Tree path:", def)
				if err != nil {
					return err
				}
				if dir == "" {
					return fmt.Errorf("tree path is required")
				}
			}
			abs, err := toolkit.WriteTree(dir)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "tree: %s\n", abs)
			fmt.Fprintf(cmd.OutOrStdout(), "add a catalog:  sbx-kit catalog add <git-url>\n")
			fmt.Fprintf(cmd.OutOrStdout(), "optional: export %s=%s\n", toolkit.TreeEnv, abs)
			return nil
		},
	}
	cmd.Flags().StringVar(&tree, "tree", "", "path to the recipes tree (parent of catalogs)")
	return cmd
}
