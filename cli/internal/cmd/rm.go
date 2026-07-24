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
		keepState bool
		force     bool
	)

	cmd := &cobra.Command{
		Use:   "rm <agent|sandbox-name> [project-dir]",
		Short: "Remove an sbx sandbox (optionally export portable state first)",
		Long: `Resolve the sandbox from a catalog agent + project binding (or an explicit
sbx name), optionally pack state to ~/.local/share/sbx-kit/profiles/<id>/,
then run sbx rm.

Without --keep-state, warns that /home/agent workplace state will be lost.`,
		Example: `  sbx-kit rm cursor .
  sbx-kit rm cursor . --keep-state
  sbx-kit rm sbxk-cursor-deadbeef --keep-state --force`,
		Args: cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectDir := "."
			if len(args) > 1 {
				projectDir = args[1]
			}
			rs, err := resolveSandboxArg(args[0], projectDir)
			if err != nil {
				return err
			}

			r := sbxutil.Default()
			exists, err := r.Exists(rs.SandboxName)
			if err != nil {
				return err
			}
			if !exists {
				return fmt.Errorf("sandbox %q not found in sbx ls (stale binding?)", rs.SandboxName)
			}

			if keepState {
				if err := statexfer.Export(r, rs.SandboxName, rs.ProfileID); err != nil {
					return err
				}
			} else {
				fmt.Printf("==> warning: removing %s without --keep-state; VM /home/agent state will be deleted\n", rs.SandboxName)
			}

			if err := r.Rm(rs.SandboxName, force); err != nil {
				return err
			}

			if rs.AgentName != "" && rs.ProjectDir != "" {
				_ = binding.Delete(rs.ProjectDir, rs.AgentName)
			}
			fmt.Printf("==> removed %s\n", rs.SandboxName)
			return nil
		},
	}

	cmd.Flags().BoolVar(&keepState, "keep-state", false, "export portable state to host XDG profile before rm")
	cmd.Flags().BoolVar(&force, "force", false, "pass --force to sbx rm")
	return cmd
}
