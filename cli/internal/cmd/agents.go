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
		Short:   "List catalog recipes (Hub or local template + kits)",
		Aliases: []string{"ls", "list", "recipes"},
		Long: `Lists recipes from config/agents.yaml under the toolkit root.

SOURCE:
  hub     no local image — sbx uses the stock agent template (official path)
  local   recipe pins image_name / template_fallback (build or registry tag)

Add your own recipes to experiment with kits on any sbx agent.`,
		Args: cobra.NoArgs,
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
			fmt.Fprintln(w, "RECIPE\tSBX_AGENT\tSOURCE\tIMAGE\tKITS\tSTATUS")
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
				source, image := recipeSource(a)
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
					name, a.SbxAgent, source, image, strings.Join(kits, ","), status)
			}
			return w.Flush()
		},
	}
}

func recipeSource(a catalog.Agent) (source, image string) {
	if a.ImageName == "" && a.TemplateFallback == "" {
		return "hub", "-"
	}
	image = a.ImageName
	if image == "" {
		image = a.TemplateFallback
	}
	return "local", image
}
