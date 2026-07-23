package cmd

import (
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/nkapatos/sbx-kit/cli/internal/catalog"
	"github.com/nkapatos/sbx-kit/cli/internal/initproj"
)

func newInitCmd() *cobra.Command {
	var agent string

	cmd := &cobra.Command{
		Use:   "init [project-dir]",
		Short: "Stamp a Docker Sandbox section into a project README",
		Long: `Writes (or updates) a short "## Docker Sandbox" section in project-dir/README.md
so the repo documents how to run it under sbx via sbx-kit.

Uses catalog agent names from config/agents.yaml (default: cursor).`,
		Example: `  sbx-kit init
  sbx-kit init ~/my-project
  sbx-kit init --agent opencode .`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectDir := "."
			if len(args) > 0 {
				projectDir = args[0]
			}
			root, err := requireToolkitRoot()
			if err != nil {
				return err
			}
			cat, err := catalog.Load(filepath.Join(root, "config", "agents.yaml"))
			if err != nil {
				return err
			}
			return initproj.Run(initproj.Opts{
				Root:       root,
				Agent:      agent,
				ProjectDir: projectDir,
				Catalog:    cat,
			})
		},
	}

	cmd.Flags().StringVar(&agent, "agent", "cursor", "catalog agent to document (see sbx-kit agents)")
	return cmd
}
