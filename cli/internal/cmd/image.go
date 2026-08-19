package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/nkapatos/sbx-kit/cli/internal/catalog"
	"github.com/nkapatos/sbx-kit/cli/internal/template"
)

func newImageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "image",
		Short: "List, build, or pull custom images",
		Long: `Manage custom images under source images/ directories.

  sbx-kit image ls                         Dockerfiles in every source
  sbx-kit image load --engine docker <source>/<name> [tag]
  sbx-kit image pull [--engine docker] <registry/tag>

Not sbx template ls (engine import store). See sbx-kit concepts.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(newImageLsCmd())
	cmd.AddCommand(newImageLoadCmd())
	cmd.AddCommand(newImagePullCmd())
	return cmd
}

func newImageLsCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "ls",
		Short:   "List custom Dockerfiles in sources",
		Long:    `Names are <source>/<image>. Not sbx template ls.`,
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
			fmt.Fprintln(w, "NAME\tDEFAULT_TAG\tROLE")
			n := 0
			for _, src := range srcs {
				imgs, err := template.ListLocal(src.Root)
				if err != nil {
					return err
				}
				for _, img := range imgs {
					n++
					fmt.Fprintf(w, "%s\t%s\t%s\n", catalog.JoinID(src.Name, img.Name), img.ImageTag, img.Role)
				}
			}
			if err := w.Flush(); err != nil {
				return err
			}
			if n == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "(no custom images under images/)")
			}
			fmt.Fprintln(cmd.OutOrStdout())
			fmt.Fprintln(cmd.OutOrStdout(), "Imported into sbx: sbx template ls")
			return nil
		},
	}
}

func newImageLoadCmd() *cobra.Command {
	var engine string

	cmd := &cobra.Command{
		Use:   "load <source>/<name-or-path> [image-tag]",
		Short: "Build a local Dockerfile and import it into sbx",
		Long: `Build then sbx template load. Names are <source>/<image>.
See sbx-kit image ls for ids.`,
		Example: `  sbx-kit image load --engine docker mine/kit-shell
  sbx-kit image load --engine docker mine/kit-cursor local/sbx-kit-cursor:dev`,
		Args: cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if engine == "" {
				return fmt.Errorf("--engine is required (docker|container)")
			}
			catalogRoot, err := requireToolkitRoot()
			if err != nil {
				return err
			}
			root, nameOrPath, err := resolveImageRef(catalogRoot, args[0])
			if err != nil {
				return err
			}
			tag := ""
			if len(args) > 1 {
				tag = args[1]
			}
			return template.Load(template.LoadOpts{
				Root:       root,
				Engine:     engine,
				NameOrPath: nameOrPath,
				ImageTag:   tag,
			})
		},
	}

	cmd.Flags().StringVar(&engine, "engine", "", "build engine: docker | container (required)")
	_ = cmd.MarkFlagRequired("engine")
	return cmd
}

func resolveImageRef(catalogRoot, ref string) (sourceRoot, nameOrPath string, err error) {
	if st, err := os.Stat(ref); err == nil && st.IsDir() {
		abs, err := filepath.Abs(ref)
		if err != nil {
			return "", "", err
		}
		root := abs
		for {
			if catalog.IsDir(root) {
				return root, abs, nil
			}
			parent := filepath.Dir(root)
			if parent == root {
				return "", abs, nil
			}
			root = parent
		}
	}
	sourceName, img, err := catalog.ParseID(ref)
	if err != nil {
		return "", "", fmt.Errorf("image id is <source>/<name> (got %q; try: sbx-kit image ls)", ref)
	}
	srcRoot := filepath.Join(catalogRoot, sourceName)
	if !catalog.IsDir(srcRoot) {
		return "", "", fmt.Errorf("unknown source %q (try: sbx-kit source ls)", sourceName)
	}
	return srcRoot, img, nil
}

func newImagePullCmd() *cobra.Command {
	var engine string

	cmd := &cobra.Command{
		Use:   "pull <registry/tag>",
		Short: "Pull a registry image and import it into sbx",
		Long:  `docker pull, then sbx template load.`,
		Example: `  sbx-kit image pull ghcr.io/example/sbx-kit-cursor:latest
  sbx-kit image pull --engine docker ghcr.io/example/sbx-kit-shell:latest`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return template.Pull(template.PullOpts{
				Engine:   engine,
				ImageTag: args[0],
			})
		},
	}

	cmd.Flags().StringVar(&engine, "engine", "docker", "pull engine (docker)")
	return cmd
}
