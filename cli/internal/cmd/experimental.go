package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/nkapatos/sbx-kit/cli/internal/experimental"
	"github.com/nkapatos/sbx-kit/cli/internal/recipespec"
)

func newExperimentalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "experimental",
		Short: "Parked recipe spec and skill helpers",
		Long: `Work in progress that is not part of the stable CLI surface.

Recipe and kit verification live under sbx-kit recipes verify.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(newExperimentalSpecCmd())
	cmd.AddCommand(newExperimentalSkillCmd())
	return cmd
}

func newExperimentalSpecCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "spec",
		Short: "Show recipe spec status (parked)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintln(cmd.OutOrStdout(), recipespec.Status)
			return experimental.ErrNotReady{Feature: "recipe spec", Track: "spec"}
		},
	}
}

func newExperimentalSkillCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "skill",
		Short: "Print path to the agent skill doc",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := skillDocPath()
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), path)
			return nil
		},
	}
}

func skillDocPath() (string, error) {
	if wd, err := os.Getwd(); err == nil {
		p := filepath.Join(wd, "docs", "sbx-kit-skill.md")
		if st, err := os.Stat(p); err == nil && !st.IsDir() {
			return p, nil
		}
	}
	return filepath.Join("docs", "sbx-kit-skill.md"), nil
}
