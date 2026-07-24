package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/nkapatos/sbx-kit/cli/internal/run"
	"github.com/nkapatos/sbx-kit/cli/internal/sbxutil"
	"github.com/nkapatos/sbx-kit/cli/internal/statexfer"
)

func newUpgradeCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "upgrade <agent> [project-dir]",
		Short: "Export state, recreate sandbox with current catalog kits/template, restore",
		Long: `Blessed path when templates/kits change:

  1. sbx-kit-state pack + sbx cp → host profile
  2. sbx rm
  3. sbx create with current catalog recipe + same --name
  4. restore archive
  5. sbx run --name (attach)

Requires agent-workspace kit so sbx-kit-state is available.`,
		Example: `  sbx-kit upgrade cursor .
  sbx-kit upgrade cursor . --force`,
		Args: cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			agentName := args[0]
			projectDir := "."
			if len(args) > 1 {
				projectDir = args[1]
			}
			rs, err := resolveFromAgent(agentName, projectDir)
			if err != nil {
				return err
			}

			r := sbxutil.Default()
			exists, err := r.Exists(rs.SandboxName)
			if err != nil {
				return err
			}
			if exists {
				if err := statexfer.Export(r, rs.SandboxName, rs.ProfileID); err != nil {
					return err
				}
				if err := r.Rm(rs.SandboxName, force); err != nil {
					return err
				}
			} else {
				fmt.Printf("==> no running sandbox %s; will create fresh and restore if archive exists\n", rs.SandboxName)
			}

			extra := extractPassthrough(os.Args)
			// Force the same name through passthrough.
			extra = append([]string{"--name", rs.SandboxName}, stripName(extra)...)

			overrideKey := "SBX_" + strings.ToUpper(agentName) + "_TEMPLATE"
			_, err = run.Sbx(run.Opts{
				Root:             rs.Root,
				AgentCatalogName: agentName,
				SbxAgent:         rs.SbxAgent,
				ImageName:        rs.ImageName,
				TemplateFallback: rs.TemplateFB,
				TemplateOverride: os.Getenv(overrideKey),
				KitPaths:         rs.KitPaths,
				ProjectDir:       rs.ProjectDir,
				Extra:            extra,
				Resources:        rs.Resources,
				ResourcesProfile: rs.ResProfile,
				RestoreState:     true,
				Runner:           r,
			})
			return err
		},
	}

	cmd.Flags().BoolVar(&force, "force", false, "pass --force to sbx rm")
	cmd.Flags().SetInterspersed(false)
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
