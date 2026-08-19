package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/nkapatos/sbx-kit/cli/internal/run"
	"github.com/nkapatos/sbx-kit/cli/internal/statexfer"
)

func newUpgradeCmd() *cobra.Command {
	var (
		recipe string
		path   string
		name   string
		force  bool
	)

	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Recreate sandbox from current recipe, keeping state",
		Long: `When the recipe's template or kits change:

  1. pack state to the host profile
  2. sbx rm
  3. create from the current recipe (same sandbox name)
  4. restore archive and attach

Requires the agent-workspace kit (sbx-kit-state).`,
		Example: `  sbx-kit box upgrade --recipe mine/shell
  sbx-kit box upgrade --recipe mine/cursor --path ~/proj --force
  sbx-kit box upgrade --name my-project --force`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			rs, err := resolveFlags(recipe, path, name)
			if err != nil {
				return err
			}
			if rs.AgentName == "" {
				return fmt.Errorf("could not resolve a recipe; pass --recipe or use a bound --name")
			}

			r := sbxRunner()
			exists, err := r.Exists(rs.SandboxName)
			if err != nil {
				return err
			}
			if exists {
				if err := statexfer.Export(r, rs.SandboxName, rs.ProfileID, UI().Out); err != nil {
					return err
				}
				if err := r.Rm(rs.SandboxName, force); err != nil {
					return err
				}
			} else {
				UI().Header("no running sandbox " + rs.SandboxName + "; will create fresh and restore if archive exists")
			}

			extra := extractPassthrough(os.Args)
			overrideKey := templateOverrideEnv(rs.AgentName)
			_, err = run.Sbx(run.Opts{
				Root:             rs.Root,
				RecipeID:         rs.AgentName,
				SbxAgent:         rs.SbxAgent,
				ImageName:        rs.ImageName,
				TemplateFallback: rs.TemplateFB,
				TemplateOverride: os.Getenv(overrideKey),
				KitPaths:         rs.KitPaths,
				ProjectDir:       rs.ProjectDir,
				Extra:            stripName(extra),
				Resources:        rs.Resources,
				ResourcesProfile: rs.ResProfile,
				RestoreState:     true,
				CreateOnly:       true,
				ConfirmCreate:    false,
				SandboxName:      rs.SandboxName,
				ProfileID:        rs.ProfileID,
				Runner:           r,
				Out:              UI().Out,
			})
			return err
		},
	}

	addTargetFlags(cmd, &recipe, &path, &name)
	cmd.Flags().BoolVar(&force, "force", false, "pass --force to sbx rm")
	return cmd
}

func stripName(args []string) []string {
	out := make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "--name" && i+1 < len(args) {
			i++
			continue
		}
		if strings.HasPrefix(a, "--name=") {
			continue
		}
		out = append(out, a)
	}
	return out
}
