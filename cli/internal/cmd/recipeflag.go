package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// coalesceRecipe returns --recipe, or --agent as alias. Errors if both set differently.
func coalesceRecipe(cmd *cobra.Command, recipe, agentAlias string) (string, error) {
	rSet := cmd.Flags().Changed("recipe")
	aSet := cmd.Flags().Changed("agent")
	if rSet && aSet && recipe != agentAlias {
		return "", fmt.Errorf("use either --recipe or --agent (alias), not both")
	}
	if recipe != "" {
		return recipe, nil
	}
	return agentAlias, nil
}

// addRecipeFlag registers --recipe and deprecated alias --agent on cmd.
func addRecipeFlag(cmd *cobra.Command, recipe, agentAlias *string, usage string) {
	if usage == "" {
		usage = "catalog recipe id"
	}
	cmd.Flags().StringVar(recipe, "recipe", "", usage)
	cmd.Flags().StringVar(agentAlias, "agent", "", "alias for --recipe")
}
