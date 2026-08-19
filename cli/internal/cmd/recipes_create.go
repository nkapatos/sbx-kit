package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/nkapatos/sbx-kit/cli/internal/recipecreate"
)

func newRecipesCreateCmd() *cobra.Command {
	var (
		recipeName string
		sbxAgent   string
		resources  string
		kits       []string
		noAgentsMD bool
		force      bool
	)

	cmd := &cobra.Command{
		Use:   "create <dir>",
		Short: "Scaffold a new catalog directory bundle",
		Long: `Creates <catalog>/<dir>/ with recipes/agents.yaml, kits/, images/,
and AGENTS.md for AI agents working on this bundle.

Kit schema is owned by sbx — see AGENTS.md links and sbx-kit recipes skill.`,
		Example: `  sbx-kit recipes create mine
  sbx-kit recipes create team --recipe cursor --sbx-agent cursor
  sbx-kit recipes create mine --kit agent-workspace --kit mise-workspace`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			catalogRoot, err := requireToolkitRoot()
			if err != nil {
				return err
			}
			return recipecreate.Create(recipecreate.CreateOpts{
				CatalogRoot: catalogRoot,
				DirName:     args[0],
				RecipeName:  recipeName,
				SbxAgent:    sbxAgent,
				DefaultKits: kits,
				Resources:   resources,
				WriteAgents: !noAgentsMD,
				Force:       force,
				Out:         UI().Out,
			})
		},
	}

	cmd.Flags().StringVar(&recipeName, "recipe", "shell", "first agent name in agents.yaml")
	cmd.Flags().StringVar(&sbxAgent, "sbx-agent", "shell", "sbx agent kind for the starter recipe")
	cmd.Flags().StringVar(&resources, "resources", "remote-llm", "defaults.resources profile name")
	cmd.Flags().StringArrayVar(&kits, "kit", nil, "defaults.kits entry (repeatable; default: agent-workspace)")
	cmd.Flags().BoolVar(&noAgentsMD, "no-agents-md", false, "skip AGENTS.md")
	cmd.Flags().BoolVar(&force, "force", false, "overwrite manifest and AGENTS.md if present")
	return cmd
}

func newRecipesSkillCmd() *cobra.Command {
	var (
		dirName string
		output  string
		cursor  bool
	)

	cmd := &cobra.Command{
		Use:   "skill",
		Short: "Render sbx-kit agent skill markdown",
		Long: `Prints or writes a Cursor-style SKILL.md for sbx-kit.

Includes sbx-kit commands and curated upstream sbx doc links (sandboxes,
customize, integrations, kit spec v2). Kit schema authority stays with sbx.`,
		Example: `  sbx-kit recipes skill
  sbx-kit recipes skill --dir mine --output ./SKILL.md
  sbx-kit recipes skill --cursor
  sbx-kit recipes skill --output ~/.cursor/skills/sbx-kit/SKILL.md`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			catalogRoot, _ := requireToolkitRoot()
			body, err := recipecreate.RenderSkill(recipecreate.SkillOpts{
				CatalogRoot: catalogRoot,
				DirName:     dirName,
			})
			if err != nil {
				return err
			}

			path := output
			if cursor {
				path = filepath.Join(".cursor", "skills", "sbx-kit", "SKILL.md")
			}
			if path == "" || path == "-" {
				fmt.Fprint(cmd.OutOrStdout(), body)
				return nil
			}
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				return err
			}
			if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "wrote %s\n", path)
			return nil
		},
	}

	cmd.Flags().StringVar(&dirName, "dir", "", "catalog directory name to mention in the skill")
	cmd.Flags().StringVarP(&output, "output", "o", "", "write file (- for stdout)")
	cmd.Flags().BoolVar(&cursor, "cursor", false, "write .cursor/skills/sbx-kit/SKILL.md")
	return cmd
}
