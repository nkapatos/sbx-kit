package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/nkapatos/sbx-kit/cli/internal/experimental"
	"github.com/nkapatos/sbx-kit/cli/internal/recipeverify"
	"github.com/nkapatos/sbx-kit/cli/internal/recipespec"
)

func newExperimentalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "experimental",
		Short: "Parked features (stubs; not production-ready)",
		Long: `Entry points for recipe spec, verify, skill, and kit v2 work in progress.

These commands exist so the CLI shape stays stable while implementation is parked.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(newExperimentalVerifyCmd())
	cmd.AddCommand(newExperimentalSpecCmd())
	cmd.AddCommand(newExperimentalSkillCmd())
	cmd.AddCommand(newExperimentalKitV2Cmd())
	return cmd
}

func newExperimentalVerifyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verify",
		Short: "Verify recipes or kits (stub)",
		RunE:  func(cmd *cobra.Command, args []string) error { return cmd.Help() },
	}
	cmd.AddCommand(newExperimentalVerifyRecipeCmd())
	cmd.AddCommand(newExperimentalVerifyKitCmd())
	return cmd
}

func newExperimentalVerifyRecipeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "recipe [id]",
		Short: "Verify recipe manifest and references (stub)",
		Long:  recipeverify.Describe(),
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			catalogRoot, err := requireToolkitRoot()
			if err != nil {
				return err
			}
			id := ""
			if len(args) == 1 {
				id = args[0]
			}
			return recipeverify.VerifyRecipe(catalogRoot, id)
		},
	}
}

func newExperimentalVerifyKitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "kit [dir]",
		Short: "Verify kits in a recipe directory via sbx (stub)",
		Long:  recipeverify.Describe(),
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			catalogRoot, err := requireToolkitRoot()
			if err != nil {
				return err
			}
			dir := ""
			if len(args) == 1 {
				dir = args[0]
			}
			return recipeverify.VerifyKits(catalogRoot, dir)
		},
	}
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

func newExperimentalKitV2Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "kit-v2",
		Short: "Kit schema v2 migration notes (stub)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintln(cmd.OutOrStdout(), `Kit schema: sbx-kit catalogs may contain schemaVersion "1" kits until rewritten to sbx v2.
Migration authority: sbx CLI and upstream kit spec — not sbx-kit.
Planned: hints and verify integration via experimental verify kit.`)
			return experimental.ErrNotReady{Feature: "kit v2 migration", Track: "kit-v2"}
		},
	}
}

func skillDocPath() (string, error) {
	// Prefer repo-relative docs/ when developing from source tree.
	if wd, err := os.Getwd(); err == nil {
		p := filepath.Join(wd, "docs", "sbx-kit-skill.md")
		if st, err := os.Stat(p); err == nil && !st.IsDir() {
			return p, nil
		}
	}
	// Installed/built binary: doc may sit next to module root in source installs only.
	return filepath.Join("docs", "sbx-kit-skill.md"), nil
}
