package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nkapatos/sbx-kit/cli/internal/template"
)

func newTemplateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "template",
		Short: "Build and import local sandbox template images",
		Long:  `Maintainer helpers for local template images. Hub recipes skip this — sbx pulls the stock agent template.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(newTemplateLoadCmd())
	return cmd
}

func newTemplateLoadCmd() *cobra.Command {
	var engine string

	cmd := &cobra.Command{
		Use:   "load <template-name-or-path> [image-tag]",
		Short: "Build a template and import it into sbx",
		Long: `Build a local template directory (Dockerfile, or bake.env → sibling _bake for
external trees), then import via sbx template load.

Not needed for Hub recipes (e.g. cursor-hub): sbx already has the stock agent image.

Engines:
  docker      Docker Desktop or Colima (docker CLI)
  container   Apple container + skopeo (OCI → docker-archive)

Template names resolve under the toolkit root (Brew share or SBX_TREE):
  sbx-kit template load --engine docker kit-core
  sbx-kit template load --engine docker kit-cursor

Not supported: OrbStack, Podman.`,
		Example: `  sbx-kit template load --engine docker kit-core
  sbx-kit template load --engine docker kit-cursor
  sbx-kit template load --engine container templates/kit-core
  sbx-kit template load --engine docker kit-cursor local/sbx-kit-cursor:dev`,
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
