package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nkapatos/sbx-kit/cli/internal/experimental"
	"github.com/nkapatos/sbx-kit/cli/internal/recipespec"
)

func newExperimentalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "experimental",
		Short: "Parked recipe spec helpers",
		Long: `Work in progress that is not part of the stable CLI surface.

Agent skill: sbx-kit recipes skill
Recipe verify: sbx-kit recipes verify`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(newExperimentalSpecCmd())
	return cmd
}

func newExperimentalSpecCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "spec",
		Short: "Show recipe spec status (parked)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintln(cmd.OutOrStdout(), recipespec.Status)
			fmt.Fprintln(cmd.OutOrStdout(), "Agent skill: sbx-kit recipes skill")
			return experimental.ErrNotReady{Feature: "recipe spec", Track: "spec"}
		},
	}
}
