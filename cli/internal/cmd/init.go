package cmd

import (
	"path/filepath"

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
		Example: `  sbx-kit init
  sbx-kit init ~/my-project
  sbx-kit init --recipe shell .`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectDir := "."
			if len(args) > 0 {
				projectDir = args[0]
			}
			recipeID := recipe
			if recipeID == "" {
				recipeID = "shell"
			}
			root, err := requireToolkitRoot()
			if err != nil {
				return err
			}
			cat, err := catalog.Load(filepath.Join(root, "config", "agents.yaml"))
			if err != nil {
				return err
			}
			return initproj.Run(initproj.Opts{
				Root:       root,
				Agent:      recipeID,
				ProjectDir: projectDir,
				Catalog:    cat,
			})
		},
	}

	addRecipeFlag(cmd, &recipe, "catalog recipe to document (see sbx-kit recipes)")
	return cmd
}
