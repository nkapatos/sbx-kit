package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/nkapatos/sbx-kit/cli/internal/catalog"
	"github.com/nkapatos/sbx-kit/cli/internal/gitdir"
)

func newCatalogCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "catalog",
		Short: "Add, list, or fetch catalogs in the recipes tree",
		Long: `A catalog is a one-level child of the tree (sbx-kit setup).

  sbx-kit catalog add <git-url>   git clone into the tree
  sbx-kit catalog ls              children (git or local)
  sbx-kit catalog fetch           git fetch && git status (no pull)

Local dirs with no .git are catalogs too; fetch skips them.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(newCatalogAddCmd())
	cmd.AddCommand(newCatalogLsCmd())
	cmd.AddCommand(newCatalogFetchCmd())
	return cmd
}

func newCatalogAddCmd() *cobra.Command {
	var name string
	cmd := &cobra.Command{
		Use:   "add <git-url>",
		Short: "Clone a git remote into the tree as a catalog",
		Long: `git clone <url> into the tree. Does not pull later; use catalog fetch
to see updates, then pull or edit by hand.

The directory name defaults to the repo name (bar from …/bar.git).`,
		Example: `  sbx-kit catalog add https://github.com/example/recipes.git
  sbx-kit catalog add https://github.com/example/recipes.git --name mine`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			tree, err := requireToolkitRoot()
			if err != nil {
				return err
			}
			url := args[0]
			destName := name
			if destName == "" {
				destName = gitdir.DirName(url)
			}
			if destName == "" || destName == "." || destName == ".." {
				return fmt.Errorf("could not derive catalog name from %s; pass --name", url)
			}
			dest := filepath.Join(tree, destName)
			if _, err := os.Stat(dest); err == nil {
				return fmt.Errorf("catalog %q already exists at %s", destName, dest)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "==> clone %s → %s\n", url, dest)
			if err := gitdir.Clone(url, dest); err != nil {
				return err
			}
			if !catalog.IsDir(dest) {
				fmt.Fprintf(cmd.OutOrStdout(), "warning: %s has no recipes/agents.yaml yet\n", dest)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "catalog: %s\n", destName)
			return nil
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "directory name under the tree (default: repo name)")
	return cmd
}

func newCatalogLsCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "ls",
		Short:   "List catalogs in the tree",
		Aliases: []string{"list"},
		Args:    cobra.NoArgs,
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
			fmt.Fprintln(w, "CATALOG\tSOURCE")
			for _, s := range srcs {
				src := "local"
				if gitdir.IsRepo(s.Root) {
					if u := gitdir.RemoteURL(s.Root); u != "" {
						src = u
					} else {
						src = "git"
					}
				}
				fmt.Fprintf(w, "%s\t%s\n", s.Name, src)
			}
			return w.Flush()
		},
	}
}

func newCatalogFetchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "fetch",
		Short: "git fetch and git status for every git catalog",
		Long: `For each catalog with a .git dir: git fetch, then git status -sb.
Does not pull. Local (non-git) catalogs are listed and skipped.`,
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
				return nil
			}
			out := cmd.OutOrStdout()
			for i, s := range srcs {
				if i > 0 {
					fmt.Fprintln(out)
				}
				if !gitdir.IsRepo(s.Root) {
					fmt.Fprintf(out, "==> %s (local)\n    not a git checkout\n", s.Name)
					continue
				}
				fmt.Fprintf(out, "==> %s\n", s.Name)
				if err := gitdir.Fetch(s.Root); err != nil {
					fmt.Fprintf(out, "    fetch: %v\n", err)
				} else {
					fmt.Fprintln(out, "    fetch: ok")
				}
				st, err := gitdir.Status(s.Root)
				if err != nil {
					fmt.Fprintf(out, "    status: %v\n", err)
					continue
				}
				fmt.Fprintf(out, "    %s\n", st)
			}
			return nil
		},
	}
}
