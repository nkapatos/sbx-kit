package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nkapatos/sbx-kit/cli/internal/binding"
	"github.com/nkapatos/sbx-kit/cli/internal/statexfer"
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
	var recipe, path, name string
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Pack VM state into the host profile archive",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			rs, err := resolveFlags(recipe, path, name)
			if err != nil {
				return err
			}
			r := sbxRunner()
			exists, err := r.Exists(rs.SandboxName)
			if err != nil {
				return err
			}
			if !exists {
				return fmt.Errorf("sandbox %q not found", rs.SandboxName)
			}
			return statexfer.Export(r, rs.SandboxName, rs.ProfileID, UI().Out)
		},
	}
	addTargetFlags(cmd, &recipe, &path, &name)
	return cmd
}

func newStateImportCmd() *cobra.Command {
	var recipe, path, name string
	cmd := &cobra.Command{
		Use:   "import",
		Short: "Restore host profile archive into the sandbox",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			rs, err := resolveFlags(recipe, path, name)
			if err != nil {
				return err
			}
			r := sbxRunner()
			exists, err := r.Exists(rs.SandboxName)
			if err != nil {
				return err
			}
			if !exists {
				return fmt.Errorf("sandbox %q not found", rs.SandboxName)
			}
			return statexfer.Import(r, rs.SandboxName, rs.ProfileID, UI().Out)
		},
	}
	addTargetFlags(cmd, &recipe, &path, &name)
	return cmd
}

// resolveFlags picks a sandbox from --name or --recipe/--path.
func resolveFlags(recipe, path, name string) (*resolvedSandbox, error) {
	if name != "" && recipe != "" {
		return nil, fmt.Errorf("use either --name or --recipe, not both")
	}
	if name != "" {
		return resolveSandboxArg(name, path)
	}
	if recipe == "" {
		if path == "" {
			path = "."
		}
		recs, err := binding.ListForProject(path)
		if err != nil {
			return nil, err
		}
		switch len(recs) {
		case 0:
			return nil, fmt.Errorf("no binding for path %s; pass --recipe <dir>/<name> or --name (see sbx-kit box bindings)", path)
		case 1:
			recipe = recs[0].Agent
		default:
			return nil, fmt.Errorf("multiple recipes bound to %s; pass --recipe or --name explicitly", path)
		}
	}
	if path == "" {
		path = "."
	}
	return resolveFromAgent(recipe, path)
}
