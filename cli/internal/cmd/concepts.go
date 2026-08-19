package cmd

import (
	"github.com/spf13/cobra"
)

func newConceptsCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "concepts",
		Short:   "How sbx-kit sits on sbx",
		Aliases: []string{"about"},
		Long:    conceptsText(),
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			UI().Printf("%s", conceptsText())
		},
	}
}

func conceptsText() string {
	return `sbx-kit is a convenience layer on Docker sbx. It does not replace sbx.

Glossary
  catalog     path sbx-kit manages (sbx-kit setup; default ~/sbx-kit-catalog)
  directory   subdirectory with recipes/ (and optional kits/, images/)
  recipe      <dir>/<name> — sbx kind + kits + optional image
  overlay     CLI-owned files in the box (/etc/sbx-kit/context.md is the index)
  kit         optional create-time YAML sbx applies at sandbox create (schema owned by sbx)
              kit agentContext is user land; it may point at overlay docs
  image       custom Dockerfile or registry tag (sbx-kit recipes image)
  box         one sbx sandbox bound to a host project path

Most commands need sbx-kit setup first, or a cwd inside your catalog.

Workflow
  sbx-kit setup
  sbx-kit catalog add <url>
  sbx-kit recipes
  sbx-kit box run <dir>/<name> --yes
  sbx-kit box bindings

See also: sbx-kit --help
`
}
