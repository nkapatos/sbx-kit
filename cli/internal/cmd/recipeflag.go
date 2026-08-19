package cmd

import (
	"github.com/spf13/cobra"
)

func addRecipeFlag(cmd *cobra.Command, recipe *string) {
	cmd.Flags().StringVar(recipe, "recipe", "", "recipe <dir>/<name> (with --path)")
}

// addTargetFlags adds the standard way to pick a sandbox for lifecycle commands.
func addTargetFlags(cmd *cobra.Command, recipe, path, name *string) {
	addRecipeFlag(cmd, recipe)
	cmd.Flags().StringVar(path, "path", ".", "project directory")
	cmd.Flags().StringVar(name, "name", "", "existing sandbox name (from sbx ls)")
}
