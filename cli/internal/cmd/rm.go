package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nkapatos/sbx-kit/cli/internal/binding"
	"github.com/nkapatos/sbx-kit/cli/internal/sbxutil"
	"github.com/nkapatos/sbx-kit/cli/internal/statexfer"
)

func newRmCmd() *cobra.Command {
	var (
		recipe    string
		path      string
		name      string
		keepState bool
		force     bool
	)

	cmd := &cobra.Command{
		Use:   "rm",
		Short: "Remove a sandbox (optionally keep portable state)",
		Long: `Resolve via --recipe/--path or --name, optionally export state to the host
vault, then sbx rm.

Without --keep-state, /home/agent workplace state is discarded with the box.`,
		Example: `  sbx-kit rm --recipe shell --keep-state
  sbx-kit rm --recipe cursor --path ~/proj --keep-state
  sbx-kit rm --name my-project --keep-state --force`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			rs, err := resolveFlags(recipe, path, name)
			if err != nil {
				return err
			}

			r := sbxutil.Default()
			exists, err := r.Exists(rs.SandboxName)
			if err != nil {
				return err
			}
			if !exists {
				return fmt.Errorf("sandbox %q not found in sbx ls (stale binding?); try: sbx-kit status", rs.SandboxName)
			}

			if keepState {
				if err := statexfer.Export(r, rs.SandboxName, rs.ProfileID); err != nil {
					return err
				}
			} else {
				fmt.Println("==> warning: removing without --keep-state discards portable /home/agent state")
			}

			if err := r.Rm(rs.SandboxName, force); err != nil {
				return err
			}
			_ = binding.Delete(rs.ProjectDir, rs.AgentName)
			return nil
		},
	}

	addRecipeFlag(cmd, &recipe, "catalog recipe (via project binding)")
	cmd.Flags().StringVar(&path, "path", ".", "project directory")
	cmd.Flags().StringVar(&name, "name", "", "existing sandbox name")
	cmd.Flags().BoolVar(&keepState, "keep-state", false, "export portable state to host vault before rm")
	cmd.Flags().BoolVar(&force, "force", false, "pass --force to sbx rm")
	return cmd
}
