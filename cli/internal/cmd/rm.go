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
		agent     string
		path      string
		name      string
		keepState bool
		force     bool
	)

	cmd := &cobra.Command{
		Use:   "rm",
		Short: "Remove a sandbox (optionally export portable state first)",
		Long: `Resolve the sandbox from --agent/--path (recipe binding) or --name, optionally
pack state to ~/.local/share/sbx-kit/profiles/<id>/, then run sbx rm.

Without --keep-state, warns that /home/agent workplace state will be lost.
With --keep-state, export best-effort waits if status is still running, then
checkpoints SQLite WALs inside the VM before packing.`,
		Example: `  sbx-kit rm --agent cursor --keep-state
  sbx-kit rm --agent cursor --path ~/proj --keep-state
  sbx-kit rm --name sbxk-cursor-deadbeef --keep-state --force`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			rs, err := resolveFlags(agent, path, name)
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
				fmt.Printf("==> warning: removing %s without --keep-state; VM /home/agent state will be deleted\n", rs.SandboxName)
			}

			if err := r.Rm(rs.SandboxName, force); err != nil {
				return err
			}

			if rs.AgentName != "" && rs.ProjectDir != "" {
				_ = binding.Delete(rs.ProjectDir, rs.AgentName)
			}
			fmt.Printf("==> removed %s (%s)\n", rs.SandboxName, binding.Label(&binding.Record{ProjectDir: rs.ProjectDir, SandboxName: rs.SandboxName}))
			return nil
		},
	}

	cmd.Flags().StringVar(&agent, "agent", "", "catalog recipe (resolve via project binding)")
	cmd.Flags().StringVar(&path, "path", ".", "project directory")
	cmd.Flags().StringVar(&name, "name", "", "sandbox id (no create)")
	cmd.Flags().BoolVar(&keepState, "keep-state", false, "export portable state to host XDG profile before rm")
	cmd.Flags().BoolVar(&force, "force", false, "pass --force to sbx rm")
	return cmd
}
