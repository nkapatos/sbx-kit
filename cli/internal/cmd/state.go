package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/nkapatos/sbx-kit/cli/internal/binding"
	"github.com/nkapatos/sbx-kit/cli/internal/sbxutil"
	"github.com/nkapatos/sbx-kit/cli/internal/statexfer"
	"github.com/nkapatos/sbx-kit/cli/internal/xdg"
)

func newStateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "state",
		Short: "Export or import portable sandbox state to the host XDG vault",
	}
	cmd.AddCommand(newStateExportCmd())
	cmd.AddCommand(newStateImportCmd())
	return cmd
}

func newStateExportCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "export <agent|sandbox-name> [project-dir]",
		Short: "Pack VM state and copy to ~/.local/share/sbx-kit/profiles/<id>/",
		Args:  cobra.RangeArgs(1, 2),
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
				return fmt.Errorf("sandbox %q not found", rs.SandboxName)
			}
			return statexfer.Export(r, rs.SandboxName, rs.ProfileID)
		},
	}
}

func newStateImportCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "import <agent|sandbox-name> [project-dir]",
		Short: "Copy host profile archive into the sandbox and unpack",
		Args:  cobra.RangeArgs(1, 2),
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
				return fmt.Errorf("sandbox %q not found", rs.SandboxName)
			}
			return statexfer.Import(r, rs.SandboxName, rs.ProfileID)
		},
	}
}

func newStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status [project-dir]",
		Short: "Show sbx-kit bindings and whether sandboxes still exist",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := xdg.Ensure(); err != nil {
				return err
			}
			share, _ := xdg.ShareDir()
			state, _ := xdg.StateDir()
			fmt.Printf("share: %s\nstate: %s\n\n", share, state)

			projectFilter := ""
			if len(args) == 1 {
				abs, err := filepath.Abs(args[0])
				if err != nil {
					return err
				}
				projectFilter = abs
			}

			recs, err := binding.List()
			if err != nil {
				return err
			}
			r := sbxutil.Default()
			rows, lsErr := r.Ls()
			alive := map[string]bool{}
			if lsErr == nil {
				for _, s := range rows {
					alive[s.Name] = true
				}
			} else {
				fmt.Printf("warning: sbx ls failed: %v\n\n", lsErr)
			}

			shown := 0
			for _, rec := range recs {
				if projectFilter != "" && rec.ProjectDir != projectFilter {
					continue
				}
				status := "missing"
				if alive[rec.SandboxName] {
					status = "present"
				}
				fmt.Printf("%s  agent=%s  sandbox=%s  profile=%s  sbx=%s\n",
					rec.ProjectDir, rec.Agent, rec.SandboxName, rec.ProfileID, status)
				shown++
			}
			if shown == 0 {
				fmt.Println("(no bindings)")
			}
			return nil
		},
	}
}
