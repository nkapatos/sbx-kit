package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/nkapatos/sbx-kit/cli/internal/toolkit"
	"github.com/nkapatos/sbx-kit/cli/internal/version"
)

const (
	groupCatalog      = "catalog"
	groupRecipes      = "recipes"
	groupBox          = "box"
	groupProject      = "project"
	groupExperimental = "experimental"
	groupOther        = "other"
)

// NewRoot builds the sbx-kit command tree.
func NewRoot() *cobra.Command {
	root := &cobra.Command{
		Use:           "sbx-kit",
		Short:         "Recipes, kits, and custom images on top of Docker sbx",
		Long:          longHelp(),
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       version.Version,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	root.SetVersionTemplate("{{.Version}}\n")
	root.SetOut(os.Stdout)
	root.SetErr(os.Stderr)

	root.AddGroup(
		&cobra.Group{ID: groupCatalog, Title: "Catalog:"},
		&cobra.Group{ID: groupRecipes, Title: "Recipes:"},
		&cobra.Group{ID: groupBox, Title: "Box:"},
		&cobra.Group{ID: groupProject, Title: "Project:"},
		&cobra.Group{ID: groupExperimental, Title: "Experimental:"},
		&cobra.Group{ID: groupOther, Title: "Other:"},
	)

	setup := newSetupCmd()
	setup.GroupID = groupCatalog
	catalog := newCatalogCmd()
	catalog.GroupID = groupCatalog

	recipes := newRecipesCmd()
	recipes.GroupID = groupRecipes

	box := newBoxCmd()
	box.GroupID = groupBox

	project := newProjectCmd()
	project.GroupID = groupProject

	experimental := newExperimentalCmd()
	experimental.GroupID = groupExperimental

	concepts := newConceptsCmd()
	concepts.GroupID = groupOther
	ver := newVersionCmd()
	ver.GroupID = groupOther

	root.AddCommand(
		setup,
		catalog,
		recipes,
		box,
		project,
		experimental,
		concepts,
		ver,
	)

	return root
}

func longHelp() string {
	return `Convenience layer on Docker sbx: recipes, kits, and custom images.

Catalog
  sbx-kit setup
  sbx-kit catalog add | ls | status | update

Recipes
  sbx-kit recipes
  sbx-kit recipes create | skill
  sbx-kit recipes verify [id] | verify kits [dir]
  sbx-kit recipes image ls | load | pull

Box
  sbx-kit box run <dir>/<name> --yes
  sbx-kit box bindings | check | upgrade | rm | state

Project
  sbx-kit project readme --recipe <dir>/<name>

Experimental (stubs)
  sbx-kit experimental spec

Glossary: sbx-kit concepts
Default catalog: ~/sbx-kit-catalog (sbx-kit setup)
Recipe id: <dir>/<name>`
}

func requireToolkitRoot() (string, error) {
	root, err := toolkit.Root()
	if err != nil {
		return "", fmt.Errorf("%w\n  tip: sbx-kit setup", err)
	}
	return root, nil
}
