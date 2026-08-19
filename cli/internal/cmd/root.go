package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/nkapatos/sbx-kit/cli/internal/toolkit"
	"github.com/nkapatos/sbx-kit/cli/internal/version"
)

const (
	groupCatalog = "catalog"
	groupSandbox = "sandbox"
	groupOther   = "other"
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
		&cobra.Group{ID: groupSandbox, Title: "Sandbox:"},
		&cobra.Group{ID: groupOther, Title: "Other:"},
	)

	setup := newSetupCmd()
	setup.GroupID = groupCatalog
	catalog := newCatalogCmd()
	catalog.GroupID = groupCatalog
	recipes := newRecipesCmd()
	recipes.GroupID = groupCatalog

	run := newRunCmd()
	run.GroupID = groupSandbox
	bindings := newBindingsCmd()
	bindings.GroupID = groupSandbox
	check := newCheckCmd()
	check.GroupID = groupSandbox
	upgrade := newUpgradeCmd()
	upgrade.GroupID = groupSandbox
	rm := newRmCmd()
	rm.GroupID = groupSandbox
	state := newStateCmd()
	state.GroupID = groupSandbox

	image := newImageCmd()
	image.GroupID = groupOther
	init := newInitCmd()
	init.GroupID = groupOther
	concepts := newConceptsCmd()
	concepts.GroupID = groupOther
	ver := newVersionCmd()
	ver.GroupID = groupOther

	root.AddCommand(
		setup,
		catalog,
		recipes,
		run,
		bindings,
		check,
		upgrade,
		rm,
		state,
		image,
		init,
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
  sbx-kit recipes

Sandbox
  sbx-kit run <dir>/<name> --yes
  sbx-kit run --name <sandbox>
  sbx-kit bindings | check | upgrade | rm | state

Other
  sbx-kit image ls | load | pull
  sbx-kit init

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
