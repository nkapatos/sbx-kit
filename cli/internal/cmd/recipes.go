package cmd

import (
	"fmt"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/nkapatos/sbx-kit/cli/internal/catalog"
	"github.com/nkapatos/sbx-kit/cli/internal/recipeverify"
	"github.com/nkapatos/sbx-kit/cli/internal/sbxutil"
)

func newRecipesCmd() *cobra.Command {
	ls := newRecipesLsCmd()
	cmd := &cobra.Command{
		Use:   "recipes",
		Short: "List and manage recipe content",
		Long:  `Recipe ids are <dir>/<name>. ` + recipeverify.Describe(),
		Example: `  sbx-kit recipes
  sbx-kit recipes create mine
  sbx-kit recipes skill --cursor
  sbx-kit recipes ls
  sbx-kit recipes verify
  sbx-kit recipes verify kits mine
  sbx-kit recipes image ls`,
		RunE: ls.RunE,
	}
	cmd.AddCommand(ls)
	cmd.AddCommand(newRecipesCreateCmd())
	cmd.AddCommand(newRecipesSkillCmd())
	cmd.AddCommand(newRecipesVerifyCmd())
	cmd.AddCommand(newImageCmd())
	return cmd
}

func newRecipesLsCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "ls",
		Short:   "List recipes in the catalog",
		Aliases: []string{"list"},
		Args:    cobra.NoArgs,
		RunE:    runRecipesList,
	}
}

func runRecipesList(cmd *cobra.Command, args []string) error {
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
}

func newRecipesVerifyCmd() *cobra.Command {
	var skipKits bool

	cmd := &cobra.Command{
		Use:   "verify [id]",
		Short: "Verify recipe manifests (and kits via sbx)",
		Long:  recipeverify.Describe(),
		Example: `  sbx-kit recipes verify
  sbx-kit recipes verify mine/cursor
  sbx-kit recipes verify --skip-kits mine/cursor
  sbx-kit recipes verify kits mine`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			catalogRoot, err := requireToolkitRoot()
			if err != nil {
				return err
			}
			id := ""
			if len(args) == 1 {
				id = args[0]
			}
			return recipeverify.VerifyRecipe(catalogRoot, id, verifyOpts(cmd, skipKits))
		},
	}
	cmd.Flags().BoolVar(&skipKits, "skip-kits", false, "recipe manifest only; do not run sbx kit verify")

	kits := &cobra.Command{
		Use:   "kits [dir]",
		Short: "Verify catalog kits with sbx",
		Long:  "Kit checks are done by sbx. sbx-kit runs sbx kit verify on each kit under kits/.",
		Example: `  sbx-kit recipes verify kits
  sbx-kit recipes verify kits mine`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			catalogRoot, err := requireToolkitRoot()
			if err != nil {
				return err
			}
			dir := ""
			if len(args) == 1 {
				dir = args[0]
			}
			opts := verifyOpts(cmd, false)
			opts.SkipKits = false
			return recipeverify.VerifyKits(catalogRoot, dir, opts)
		},
	}
	cmd.AddCommand(kits)
	return cmd
}

func verifyOpts(cmd *cobra.Command, skipKits bool) recipeverify.Options {
	r := sbxutil.Default()
	r.Out = UI().Out
	return recipeverify.Options{
		Out:      cmd.OutOrStdout(),
		SkipKits: skipKits,
		Runner: recipeverify.SbxKitRunner{
			ProbeVersion: r.ProbeVersion,
			KitVerifyFn:  r.KitVerify,
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
