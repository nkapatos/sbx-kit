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
	var agent, path, name string
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Pack VM state and copy to ~/.local/share/sbx-kit/profiles/<id>/",
		Args:  cobra.NoArgs,
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
				return fmt.Errorf("sandbox %q not found", rs.SandboxName)
			}
			return statexfer.Export(r, rs.SandboxName, rs.ProfileID)
		},
	}
	cmd.Flags().StringVar(&agent, "agent", "", "catalog recipe (resolve via project binding)")
	cmd.Flags().StringVar(&path, "path", ".", "project directory")
	cmd.Flags().StringVar(&name, "name", "", "sandbox name (alternative to --agent)")
	return cmd
}

func newStateImportCmd() *cobra.Command {
	var agent, path, name string
	cmd := &cobra.Command{
		Use:   "import",
		Short: "Copy host profile archive into the sandbox and unpack",
		Args:  cobra.NoArgs,
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
				return fmt.Errorf("sandbox %q not found", rs.SandboxName)
			}
			return statexfer.Import(r, rs.SandboxName, rs.ProfileID)
		},
	}
	cmd.Flags().StringVar(&agent, "agent", "", "catalog recipe (resolve via project binding)")
	cmd.Flags().StringVar(&path, "path", ".", "project directory")
	cmd.Flags().StringVar(&name, "name", "", "sandbox name (alternative to --agent)")
	return cmd
}

func newStatusCmd() *cobra.Command {
	var path string
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show sbx-kit bindings and whether sandboxes still exist",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := xdg.Ensure(); err != nil {
				return err
			}
			share, _ := xdg.ShareDir()
			state, _ := xdg.StateDir()
			fmt.Printf("share: %s\nstate: %s\n\n", share, state)

			projectFilter := ""
			if cmd.Flags().Changed("path") {
				abs, err := filepath.Abs(path)
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
			alive := map[string]string{}
			if lsErr == nil {
				for _, s := range rows {
					alive[s.Name] = s.Status
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
				if st, ok := alive[rec.SandboxName]; ok {
					if st == "" {
						status = "present"
					} else {
						status = st
					}
				}
				fmt.Printf("%s  label=%s  recipe=%s  sandbox=%s  profile=%s  sbx=%s\n",
					rec.ProjectDir, binding.Label(&rec), rec.Agent, rec.SandboxName, rec.ProfileID, status)
				shown++
			}
			if shown == 0 {
				fmt.Println("(no bindings)")
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&path, "path", ".", "filter bindings to this project directory")
	return cmd
}

// resolveFlags picks a sandbox from --name or --agent/--path.
func resolveFlags(agent, path, name string) (*resolvedSandbox, error) {
	if name != "" && agent != "" {
		return nil, fmt.Errorf("use either --name or --agent, not both")
	}
	if name != "" {
		return resolveSandboxArg(name, path)
	}
	if agent == "" {
		if path == "" {
			path = "."
		}
		recs, err := binding.ListForProject(path)
		if err != nil {
			return nil, err
		}
		switch len(recs) {
		case 0:
			return nil, fmt.Errorf("no binding for path %s; pass --agent or --name (see sbx-kit status)", path)
		case 1:
			agent = recs[0].Agent
		default:
			return nil, fmt.Errorf("multiple agents bound to %s; pass --agent explicitly", path)
		}
	}
	if path == "" {
		path = "."
	}
	return resolveFromAgent(agent, path)
}
