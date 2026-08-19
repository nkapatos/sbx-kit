package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nkapatos/sbx-kit/cli/internal/catalog"
	"github.com/nkapatos/sbx-kit/cli/internal/initproj"
)

func newProjectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project",
		Short: "Host project helpers (not box lifecycle)",
		Long:  `Commands that modify the host project repo, not sbx-kit itself.`,
		RunE:  func(cmd *cobra.Command, args []string) error { return cmd.Help() },
	}
	cmd.AddCommand(newProjectReadmeCmd())
	return cmd
}

func newProjectReadmeCmd() *cobra.Command {
	var recipe string

	cmd := &cobra.Command{
		Use:   "readme [project-dir]",
		Short: "Add a Docker Sandbox section to the project README",
		Long:  `Writes or updates "## Docker Sandbox" using a catalog recipe.`,
		Example: `  sbx-kit project readme --recipe mine/shell .
  sbx-kit project readme --recipe mine/cursor ~/my-project`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectDir := "."
			if len(args) > 0 {
				projectDir = args[0]
			}
			if recipe == "" {
				return fmt.Errorf("pass --recipe <dir>/<name> (see sbx-kit recipes)")
			}
			catalogRoot, err := requireToolkitRoot()
			if err != nil {
				return err
			}
			src, manifest, _, err := catalog.Lookup(catalogRoot, recipe)
			if err != nil {
				return err
			}
			_, recName, err := catalog.ParseID(recipe)
			if err != nil {
				return err
			}
			return initproj.Run(initproj.Opts{
				Root:       src.Root,
				Agent:      recName,
				RecipeID:   recipe,
				ProjectDir: projectDir,
				Manifest:   manifest,
			})
		},
	}

	addRecipeFlag(cmd, &recipe)
	return cmd
}
