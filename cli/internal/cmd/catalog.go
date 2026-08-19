package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/nkapatos/sbx-kit/cli/internal/catalog"
	"github.com/nkapatos/sbx-kit/cli/internal/gitdir"
)

func newCatalogCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "catalog",
		Short: "Manage recipe directories in the catalog",
		Long: `Add, list, check, or update directories that hold recipes.

  sbx-kit catalog add <url>       add a directory
  sbx-kit catalog ls              list directories
  sbx-kit catalog status [dir]    check for upstream updates
  sbx-kit catalog update [dir]    update directories

Configure the catalog path with sbx-kit setup.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(newCatalogAddCmd())
	cmd.AddCommand(newCatalogLsCmd())
	cmd.AddCommand(newCatalogStatusCmd())
	cmd.AddCommand(newCatalogUpdateCmd())
	return cmd
}

func newCatalogAddCmd() *cobra.Command {
	var name string
	cmd := &cobra.Command{
		Use:   "add <url>",
		Short: "Add a recipe directory",
		Example: `  sbx-kit catalog add https://github.com/example/recipes.git
  sbx-kit catalog add https://github.com/example/recipes.git --name mine`,
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
				return fmt.Errorf("could not derive directory name from %s; pass --name", url)
			}
			dest := filepath.Join(catalogRoot, destName)
			if _, err := os.Stat(dest); err == nil {
				return fmt.Errorf("directory %q already exists at %s", destName, dest)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "==> add %s\n", destName)
			if err := gitdir.Clone(url, dest); err != nil {
				return err
			}
			if !catalog.IsDir(dest) {
				fmt.Fprintf(cmd.OutOrStdout(), "warning: %s/recipes/agents.yaml not found yet\n", destName)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "directory name in the catalog (default: repo name)")
	return cmd
}

func newCatalogLsCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "ls",
		Short:   "List recipe directories",
		Aliases: []string{"list"},
		Args:    cobra.NoArgs,
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
			fmt.Fprintln(w, "DIR\tORIGIN")
			for _, d := range dirs {
				origin := "local"
				if gitdir.IsRepo(d.Root) {
					if u := gitdir.RemoteURL(d.Root); u != "" {
						origin = u
					} else {
						origin = "git"
					}
				}
				fmt.Fprintf(w, "%s\t%s\n", d.Name, origin)
			}
			return w.Flush()
		},
	}
}

func newCatalogStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status [dir]",
		Short: "Check recipe directories for upstream updates",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			filter := ""
			if len(args) == 1 {
				filter = args[0]
			}
			return runCatalogSync(cmd, filter, false)
		},
	}
}

func newCatalogUpdateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update [dir]",
		Short: "Update recipe directories",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			filter := ""
			if len(args) == 1 {
				filter = args[0]
			}
			return runCatalogSync(cmd, filter, true)
		},
	}
}

func runCatalogSync(cmd *cobra.Command, filter string, pull bool) error {
	catalogRoot, err := requireToolkitRoot()
	if err != nil {
		return err
	}
	dirs, err := catalog.List(catalogRoot)
	if err != nil {
		return err
	}
	dirs, err = catalog.FilterDirs(dirs, filter)
	if err != nil {
		return err
	}
	if len(dirs) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "(no directories)")
		return nil
	}

	out := cmd.OutOrStdout()
	w := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	if pull {
		fmt.Fprintln(w, "DIR\tRESULT")
	} else {
		fmt.Fprintln(w, "DIR\tSTATE")
	}
	for _, d := range dirs {
		if !gitdir.IsRepo(d.Root) {
			fmt.Fprintf(w, "%s\tlocal\n", d.Name)
			continue
		}
		if pull {
			if err := gitdir.Pull(d.Root); err != nil {
				fmt.Fprintf(w, "%s\terror: %v\n", d.Name, err)
			} else {
				fmt.Fprintf(w, "%s\tupdated\n", d.Name)
			}
			continue
		}
		if err := gitdir.Fetch(d.Root); err != nil {
			fmt.Fprintf(w, "%s\terror: %v\n", d.Name, err)
			continue
		}
		st, err := gitdir.Status(d.Root)
		if err != nil {
			fmt.Fprintf(w, "%s\terror: %v\n", d.Name, err)
			continue
		}
		fmt.Fprintf(w, "%s\t%s\n", d.Name, summarizeSyncState(st))
	}
	return w.Flush()
}

func summarizeSyncState(status string) string {
	line := status
	if i := strings.IndexByte(status, '\n'); i >= 0 {
		line = status[:i]
	}
	line = strings.TrimSpace(line)
	switch {
	case strings.Contains(line, "[behind") && strings.Contains(line, "[ahead"):
		return "diverged"
	case strings.Contains(line, "[behind"):
		return "updates available"
	case strings.Contains(line, "[ahead"):
		return "ahead of upstream"
	case strings.Contains(status, "\n??"), strings.Contains(status, "\n M"), strings.Contains(status, "\nM"):
		return "local changes"
	case line == "##" || line == "":
		return "unknown"
	default:
		return "up to date"
	}
}
