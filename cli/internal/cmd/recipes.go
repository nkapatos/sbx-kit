package cmd

import (
	"fmt"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/nkapatos/sbx-kit/cli/internal/catalog"
)

func newRecipesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "recipes",
		Short: "List recipes in the catalog",
		Long:  `Recipe ids are <dir>/<name>. See sbx-kit concepts.`,
		Example: `  sbx-kit recipes
  sbx-kit run mine/cursor --yes`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			catalogRoot, err := requireToolkitRoot()
			if err != nil {
				return err
			}
			dirs, err := catalog.List(catalogRoot)
			if err != nil {
				return err
			}
			if len(dirs) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "(no directories)")
				fmt.Fprintln(cmd.OutOrStdout(), "add one:  sbx-kit catalog add <url>")
				return nil
			}

			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "RECIPE\tSBX_AGENT\tIMAGE\tKITS\tSTATUS")
			any := false
			for _, d := range dirs {
				manifest, err := catalog.Load(catalog.File(d.Root))
				if err != nil {
					return err
				}
				names := make([]string, 0, len(manifest.Agents))
				for name := range manifest.Agents {
					names = append(names, name)
				}
				sort.Strings(names)
				for _, name := range names {
					any = true
					a := manifest.Agents[name]
					status := "ready"
					if a.Stub {
						status = "stub"
					}
					kits := catalog.ResolveKits(a.Kits, manifest.Defaults.Kits)
					image := recipeImage(a)
					id := catalog.JoinID(d.Name, name)
					fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
						id, a.SbxAgent, image, strings.Join(kits, ","), status)
				}
			}
			if err := w.Flush(); err != nil {
				return err
			}
			if !any {
				fmt.Fprintln(cmd.OutOrStdout(), "(no recipes)")
			}
			return nil
		},
	}
}

func recipeImage(a catalog.Agent) string {
	if a.ImageName != "" {
		return a.ImageName
	}
	if a.TemplateFallback != "" {
		return a.TemplateFallback
	}
	return "-"
}
