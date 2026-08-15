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
		Short: "Export or import portable workplace state",
	}
	cmd.AddCommand(newStateExportCmd())
	cmd.AddCommand(newStateImportCmd())
	return cmd
}

func newStateExportCmd() *cobra.Command {
	var recipe, agentAlias, path, name string
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Pack VM state into the host profile archive",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			recipeID, err := coalesceRecipe(cmd, recipe, agentAlias)
			if err != nil {
				return err
			}
			rs, err := resolveFlags(recipeID, path, name)
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
	addRecipeFlag(cmd, &recipe, &agentAlias, "catalog recipe (via project binding)")
	cmd.Flags().StringVar(&path, "path", ".", "project directory")
	cmd.Flags().StringVar(&name, "name", "", "existing sandbox name")
	return cmd
}

func newStateImportCmd() *cobra.Command {
	var recipe, agentAlias, path, name string
	cmd := &cobra.Command{
		Use:   "import",
		Short: "Restore host profile archive into the sandbox",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			recipeID, err := coalesceRecipe(cmd, recipe, agentAlias)
			if err != nil {
				return err
			}
			rs, err := resolveFlags(recipeID, path, name)
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
	addRecipeFlag(cmd, &recipe, &agentAlias, "catalog recipe (via project binding)")
	cmd.Flags().StringVar(&path, "path", ".", "project directory")
	cmd.Flags().StringVar(&name, "name", "", "existing sandbox name")
	return cmd
}

func newStatusCmd() *cobra.Command {
	var path string
	cmd := &cobra.Command{
		Use:   "status",
		Short: "List recipe↔sandbox bindings",
		Long:  `Shows project bindings and whether each sandbox still appears in sbx ls.`,
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
				fmt.Printf("%s  name=%s  recipe=%s  profile=%s  sbx=%s\n",
					rec.ProjectDir, rec.SandboxName, rec.Agent, rec.ProfileID, status)
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
			return nil, fmt.Errorf("no binding for path %s; pass --recipe <id> or --name (see sbx-kit status)", path)
		case 1:
			agent = recs[0].Agent
		default:
			return nil, fmt.Errorf("multiple recipes bound to %s; pass --recipe or --name explicitly", path)
		}
	}
	if path == "" {
		path = "."
	}
	return resolveFromAgent(agent, path)
}
