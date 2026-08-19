package cmd

import (
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/nkapatos/sbx-kit/cli/internal/binding"
	"github.com/nkapatos/sbx-kit/cli/internal/ui"
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
			filter := ""
			if cmd.Flags().Changed("path") {
				abs, err := filepath.Abs(path)
				if err != nil {
					return err
				}
				filter = abs
			}
			res, err := computeBindings(filter)
			if err != nil {
				return err
			}
			renderBindings(UI(), res)
			return nil
		},
	}
	cmd.Flags().StringVar(&path, "path", ".", "filter to this project directory")
	return cmd
}

// BindingRow is one project↔sandbox mapping.
type BindingRow struct {
	Project string
	Name    string
	Recipe  string
	Profile string
	Status  string
}

// BindingsResult is the data for box bindings.
type BindingsResult struct {
	Share  string
	State  string
	Rows   []BindingRow
	LsWarn string
}

func computeBindings(projectFilter string) (*BindingsResult, error) {
	if err := xdg.Ensure(); err != nil {
		return nil, err
	}
	share, _ := xdg.ShareDir()
	state, _ := xdg.StateDir()
	res := &BindingsResult{Share: share, State: state}

	recs, err := binding.List()
	if err != nil {
		return nil, err
	}
	r := sbxRunner()
	alive := map[string]string{}
	if rows, lsErr := r.Ls(); lsErr != nil {
		res.LsWarn = lsErr.Error()
	} else {
		for _, s := range rows {
			alive[s.Name] = s.Status
		}
	}

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
		res.Rows = append(res.Rows, BindingRow{
			Project: rec.ProjectDir,
			Name:    rec.SandboxName,
			Recipe:  rec.Agent,
			Profile: rec.ProfileID,
			Status:  sbxState,
		})
	}
	return res, nil
}

func renderBindings(w *ui.Writer, res *BindingsResult) {
	w.Detail("share", res.Share)
	w.Detail("state", res.State)
	w.Println()
	if res.LsWarn != "" {
		w.Warn("sbx ls failed: " + res.LsWarn)
		w.Println()
	}
	if len(res.Rows) == 0 {
		w.Empty("bindings", "")
		return
	}
	rows := make([][]string, 0, len(res.Rows))
	for _, r := range res.Rows {
		rows = append(rows, []string{r.Project, r.Name, r.Recipe, r.Profile, r.Status})
	}
	_ = w.Table([]string{"PROJECT", "NAME", "RECIPE", "PROFILE", "SBX"}, rows)
}
