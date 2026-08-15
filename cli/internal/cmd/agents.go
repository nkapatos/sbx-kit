package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/spf13/cobra"

	"github.com/nkapatos/sbx-kit/cli/internal/catalog"
	"github.com/nkapatos/sbx-kit/cli/internal/sbxutil"
)

func newAgentsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "agents",
		Short: "Show sbx agents and custom templates in view",
		Long: `Inventory aligned with sbx terminology:

  • tries sbx agents (or similar) when available
  • lists sbx agents referenced by recipes in this toolkit tree
  • points at sbx-kit template ls for loaded / custom images

For composed shortcuts (agent + kits), use: sbx-kit recipes`,
		Example: `  sbx-kit agents
  sbx-kit recipes
  sbx-kit template ls`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			out := cmd.OutOrStdout()

			r := sbxutil.Default()
			if _, err := r.LookPath(); err != nil {
				fmt.Fprintln(out, "==> sbx not on PATH; skipping live agent list")
			} else {
				fmt.Fprintln(out, "==> sbx agents (pass-through)")
				if err := r.RunInteractive("agents"); err != nil {
					fmt.Fprintf(os.Stderr, "note: sbx agents failed (%v); continuing with catalog view\n", err)
				}
				fmt.Fprintln(out)
			}

			root, err := requireToolkitRoot()
			if err != nil {
				return err
			}
			cat, err := catalog.Load(filepath.Join(root, "config", "agents.yaml"))
			if err != nil {
				return err
			}

			sbxAgents := map[string]struct{}{}
			var local []string
			for id, a := range cat.Agents {
				if a.SbxAgent != "" {
					sbxAgents[a.SbxAgent] = struct{}{}
				}
				if a.ImageName != "" || a.TemplateFallback != "" {
					tag := a.TemplateFallback
					if tag == "" {
						tag = a.ImageName
					}
					local = append(local, fmt.Sprintf("%s  (recipe %s, sbx agent %s)", tag, id, a.SbxAgent))
				}
			}

			names := make([]string, 0, len(sbxAgents))
			for n := range sbxAgents {
				names = append(names, n)
			}
			sort.Strings(names)
			fmt.Fprintln(out, "==> sbx agents used by recipes in this tree")
			for _, n := range names {
				fmt.Fprintf(out, "  %s\n", n)
			}

			fmt.Fprintln(out)
			fmt.Fprintln(out, "==> custom / local templates pinned by recipes")
			if len(local) == 0 {
				fmt.Fprintln(out, "  (none)")
			} else {
				sort.Strings(local)
				for _, line := range local {
					fmt.Fprintf(out, "  %s\n", line)
				}
			}
			fmt.Fprintln(out)
			fmt.Fprintln(out, "Loaded images:  sbx-kit template ls")
			fmt.Fprintln(out, "Recipes:        sbx-kit recipes")
			return nil
		},
	}
}
