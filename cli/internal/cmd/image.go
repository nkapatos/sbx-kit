package cmd

import (
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/nkapatos/sbx-kit/cli/internal/template"
)

func newImageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "image",
		Short: "List, build, or pull custom images",
		Long: `Manage this tree's custom images (images/ Dockerfiles).

  sbx-kit image ls                         Dockerfiles we ship
  sbx-kit image load --engine docker <name> [tag]
  sbx-kit image pull [--engine docker] <registry/tag>

This is not sbx template ls. That command lists images already imported
into the sbx engine. After load or pull, run sbx template ls to confirm.

Stock recipes (cursor, shell) do not need a custom image.`,
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
		Short:   "List custom Dockerfiles in this toolkit tree",
		Long:    `Lists images/ directories that have a Dockerfile (or bake.env). Not sbx template ls.`,
		Aliases: []string{"list"},
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := requireToolkitRoot()
			if err != nil {
				return err
			}
			imgs, err := template.ListLocal(root)
			if err != nil {
				return err
			}
			if len(imgs) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "(no custom images under images/)")
				fmt.Fprintln(cmd.OutOrStdout(), "Engine store: sbx template ls")
				return nil
			}
			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "NAME\tDEFAULT_TAG\tROLE")
			for _, img := range imgs {
				fmt.Fprintf(w, "%s\t%s\t%s\n", img.Name, img.ImageTag, img.Role)
			}
			if err := w.Flush(); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout())
			fmt.Fprintln(cmd.OutOrStdout(), "ROLE parent = Docker FROM only (not imported). load = sbx template load.")
			fmt.Fprintln(cmd.OutOrStdout(), "Imported into sbx: sbx template ls")
			return nil
		},
	}
}

func newImageLoadCmd() *cobra.Command {
	var engine string

	cmd := &cobra.Command{
		Use:   "load <name-or-path> [image-tag]",
		Short: "Build a local Dockerfile and import it into sbx",
		Long: `Build a local template directory (Dockerfile, or bake.env → sibling _bake),
then import via sbx template load.

Not needed for stock recipes (cursor, shell): sbx already has the Hub image.

Engines:
  docker      Docker Desktop or Colima (docker CLI)
  container   Apple container + skopeo (OCI → docker-archive)

Names resolve under the toolkit root:
  sbx-kit image load --engine docker kit-shell
  sbx-kit image load --engine docker kit-cursor

kit-core is a Docker FROM parent (also the intended VPS host floor later).
It is docker-built automatically; do not import it into sbx.

Not supported: OrbStack, Podman.`,
		Example: `  sbx-kit image load --engine docker kit-shell
  sbx-kit image load --engine docker kit-cursor
  sbx-kit image load --engine container images/kit-shell
  sbx-kit image load --engine docker kit-cursor local/sbx-kit-cursor:dev`,
		Args: cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if engine == "" {
				return fmt.Errorf("--engine is required (docker|container)")
			}
			root, err := requireToolkitRoot()
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
				NameOrPath: args[0],
				ImageTag:   tag,
			})
		},
	}

	cmd.Flags().StringVar(&engine, "engine", "", "build engine: docker | container (required)")
	_ = cmd.MarkFlagRequired("engine")
	return cmd
}

func newImagePullCmd() *cobra.Command {
	var engine string

	cmd := &cobra.Command{
		Use:   "pull <registry/tag>",
		Short: "Pull a registry image and import it into sbx",
		Long: `docker pull, save as docker-archive, then sbx template load so the tag
appears in sbx template ls.

sbx cannot pull arbitrary registries easily; this is the workaround.

Engine is docker (Docker Desktop / Colima). Apple container is not supported
for pull yet.`,
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
