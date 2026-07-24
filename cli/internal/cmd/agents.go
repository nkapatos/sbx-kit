package cmd

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/nkapatos/sbx-kit/cli/internal/catalog"
)

func newAgentsCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "agents",
		Short:   "List catalog recipes (example recipes in this tree; bring your own)",
		Aliases: []string{"ls", "list", "recipes"},
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := requireToolkitRoot()
			if err != nil {
				return err
			}
			cat, err := catalog.Load(filepath.Join(root, "config", "agents.yaml"))
			if err != nil {
				return err
			}

			names := make([]string, 0, len(cat.Agents))
			for name := range cat.Agents {
				names = append(names, name)
			}
			sort.Strings(names)

			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "RECIPE\tSBX_AGENT\tIMAGE\tKITS\tSTATUS")
			for _, name := range names {
				a := cat.Agents[name]
				status := "ready"
				if a.Stub {
					status = "stub"
				}
				kits := a.Kits
				if len(kits) == 0 {
					kits = cat.Defaults.Kits
				}
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
					name, a.SbxAgent, a.ImageName, strings.Join(kits, ","), status)
			}
			return w.Flush()
		},
	}
}
