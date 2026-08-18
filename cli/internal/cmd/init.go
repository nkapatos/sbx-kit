package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nkapatos/sbx-kit/cli/internal/catalog"
	"github.com/nkapatos/sbx-kit/cli/internal/initproj"
)

func newInitCmd() *cobra.Command {
	var recipe string

	cmd := &cobra.Command{
		Use:   "init [project-dir]",
		Short: "Add a Docker Sandbox section to the project README",
		Long:  `Writes or updates a short "## Docker Sandbox" section using a catalog recipe.`,
		Example: `  sbx-kit init --recipe mine/shell .
  sbx-kit init --recipe mine/cursor ~/my-project`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectDir := "."
			if len(args) > 0 {
				projectDir = args[0]
			}
			if recipe == "" {
				return fmt.Errorf("pass --recipe <catalog>/<name> (see sbx-kit recipes)")
			}
			tree, err := requireToolkitRoot()
			if err != nil {
				return err
			}
			src, cat, _, err := catalog.Lookup(tree, recipe)
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
				Catalog:    cat,
			})
		},
	}

	addRecipeFlag(cmd, &recipe, "catalog recipe <catalog>/<name>")
	return cmd
}
