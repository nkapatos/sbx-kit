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
		Short: "List catalog recipes (kind + kits shortcuts)",
		Long: `Recipes are sbx-kit shortcuts: an sbx kind, optional custom image, and kits.

IDs are <catalog>/<name> (catalog = one-level child of the tree).

Mixin kits stack on a Hub kind (cursor, shell). A sandbox kit IS the kind:
sbx_agent must match the kit name (pi → sbx run pi --kit …/pi), not shell.
Catalog defaults (agent-workspace) are always attached.

SOURCE:
  stock    no image pin (Hub kind, or sandbox kit owns sandbox.image)
  custom   recipe pins image_name / template_fallback (local/… or a registry tag)`,
		Example: `  sbx-kit recipes
  sbx-kit run mine/shell --yes`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			tree, err := requireToolkitRoot()
			if err != nil {
				return err
			}
			srcs, err := catalog.List(tree)
			if err != nil {
				return err
			}
			if len(srcs) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "(no catalogs)")
				fmt.Fprintln(cmd.OutOrStdout(), "add one:  sbx-kit catalog add <git-url>")
				return nil
			}

			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "RECIPE\tSBX_AGENT\tSOURCE\tIMAGE\tKITS\tSTATUS")
			any := false
			for _, src := range srcs {
				cat, err := catalog.Load(catalog.File(src.Root))
				if err != nil {
					return err
				}
				names := make([]string, 0, len(cat.Agents))
				for name := range cat.Agents {
					names = append(names, name)
				}
				sort.Strings(names)
				for _, name := range names {
					any = true
					a := cat.Agents[name]
					status := "ready"
					if a.Stub {
						status = "stub"
					}
					kits := catalog.ResolveKits(a.Kits, cat.Defaults.Kits)
					source, image := recipeSource(a)
					id := catalog.JoinID(src.Name, name)
					fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
						id, a.SbxAgent, source, image, strings.Join(kits, ","), status)
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

func recipeSource(a catalog.Agent) (source, image string) {
	if a.ImageName == "" && a.TemplateFallback == "" {
		return "stock", "-"
	}
	image = a.ImageName
	if image == "" {
		image = a.TemplateFallback
	}
	return "custom", image
}
