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

func newSourceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "source",
		Short: "Add, list, or fetch sources in the catalog",
		Long: `Sources are subdirectories of the catalog (sbx-kit setup).

  sbx-kit source add <git-url>   git clone into the catalog
  sbx-kit source ls              sources (git or local)
  sbx-kit source fetch           git fetch && git status (no pull)

Local dirs without .git are sources too; fetch skips them.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(newSourceAddCmd())
	cmd.AddCommand(newSourceLsCmd())
	cmd.AddCommand(newSourceFetchCmd())
	return cmd
}

func newSourceAddCmd() *cobra.Command {
	var name string
	cmd := &cobra.Command{
		Use:   "add <git-url>",
		Short: "Clone a git remote into the catalog as a source",
		Long: `git clone <url> into the catalog. Does not pull later; use source fetch
to see updates, then pull or edit by hand.

The directory name defaults to the repo name (bar from …/bar.git).`,
		Example: `  sbx-kit source add https://github.com/example/recipes.git
  sbx-kit source add https://github.com/example/recipes.git --name mine`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			catalogRoot, err := requireToolkitRoot()
			if err != nil {
				return err
			}
			url := args[0]
			destName := name
			if destName == "" {
				destName = gitdir.DirName(url)
			}
			if destName == "" || destName == "." || destName == ".." {
				return fmt.Errorf("could not derive source name from %s; pass --name", url)
			}
			dest := filepath.Join(catalogRoot, destName)
			if _, err := os.Stat(dest); err == nil {
				return fmt.Errorf("source %q already exists at %s", destName, dest)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "==> clone %s → %s\n", url, dest)
			if err := gitdir.Clone(url, dest); err != nil {
				return err
			}
			if !catalog.IsDir(dest) {
				fmt.Fprintf(cmd.OutOrStdout(), "warning: %s has no recipes/agents.yaml yet\n", dest)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "source: %s\n", destName)
			return nil
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "directory name under the catalog (default: repo name)")
	return cmd
}

func newSourceLsCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "ls",
		Short:   "List sources in the catalog",
		Aliases: []string{"list"},
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			catalogRoot, err := requireToolkitRoot()
			if err != nil {
				return err
			}
			srcs, err := catalog.List(catalogRoot)
			if err != nil {
				return err
			}
			if len(srcs) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "(no sources)")
				fmt.Fprintln(cmd.OutOrStdout(), "add one:  sbx-kit source add <git-url>")
				return nil
			}
			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "SOURCE\tORIGIN")
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

func newSourceFetchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "fetch",
		Short: "git fetch and git status for every git source",
		Long: `For each source with a .git dir: git fetch, then git status -sb.
Does not pull. Local (non-git) sources are listed and skipped.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			catalogRoot, err := requireToolkitRoot()
			if err != nil {
				return err
			}
			srcs, err := catalog.List(catalogRoot)
			if err != nil {
				return err
			}
			if len(srcs) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "(no sources)")
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
