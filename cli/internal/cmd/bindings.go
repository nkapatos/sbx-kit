package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/nkapatos/sbx-kit/cli/internal/binding"
	"github.com/nkapatos/sbx-kit/cli/internal/sbxutil"
	"github.com/nkapatos/sbx-kit/cli/internal/xdg"
)

func newBindingsCmd() *cobra.Command {
	var path string
	cmd := &cobra.Command{
		Use:   "bindings",
		Short: "List project ↔ sandbox bindings",
		Long:  `Shows which recipe and sandbox sbx-kit associates with each project.`,
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
				sbxState := "missing"
				if st, ok := alive[rec.SandboxName]; ok {
					if st == "" {
						sbxState = "present"
					} else {
						sbxState = st
					}
				}
				fmt.Printf("%s  name=%s  recipe=%s  profile=%s  sbx=%s\n",
					rec.ProjectDir, rec.SandboxName, rec.Agent, rec.ProfileID, sbxState)
				shown++
			}
			if shown == 0 {
				fmt.Println("(no bindings)")
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&path, "path", ".", "filter to this project directory")
	return cmd
}
