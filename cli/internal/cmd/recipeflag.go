package cmd

import (
	"github.com/spf13/cobra"
)

func addRecipeFlag(cmd *cobra.Command, recipe *string, usage string) {
	if usage == "" {
		usage = "recipe <source>/<name>"
	}
	cmd.Flags().StringVar(recipe, "recipe", "", usage)
}
