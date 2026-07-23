package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/spf13/cobra"

	"github.com/nkapatos/sbx-kit/cli/internal/catalog"
	"github.com/nkapatos/sbx-kit/cli/internal/resources"
	"github.com/nkapatos/sbx-kit/cli/internal/run"
)

func newRunCmd() *cobra.Command {
	var (
		resourcesProfile string
		clone            bool
	)

	cmd := &cobra.Command{
		Use:   "run <agent> [project-dir]",
		Short: "Run an agent sandbox from the catalog",
		Long: `Resolve the agent recipe (template image + kits + resources) and exec sbx.

Pass-through: anything after -- goes to sbx run unchanged:
  sbx-kit run cursor . -- --name my-sandbox

--clone on this command injects sbx's --clone (workspace isolation).`,
		Example: `  sbx-kit run cursor
  sbx-kit run cursor ~/my-project
  sbx-kit run cursor . --clone
  sbx-kit run cursor . -- --name feature
  SBX_MEMORY=8g sbx-kit run cursor .`,
		Args: cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			agentName := args[0]
			projectDir := "."
			if len(args) > 1 {
				projectDir = args[1]
			}

			root, err := requireToolkitRoot()
			if err != nil {
				return err
			}

			cat, err := catalog.Load(filepath.Join(root, "config", "agents.yaml"))
			if err != nil {
				return err
			}
			agent, ok := cat.Agents[agentName]
			if !ok {
				return fmt.Errorf("unknown agent %q (try: sbx-kit agents)", agentName)
			}
			if agent.Stub {
				return fmt.Errorf("agent %q is still a stub in config/agents.yaml", agentName)
			}

			profile := resourcesProfile
			if profile == "" {
				profile = cat.Defaults.Resources
			}
			if profile == "" {
				profile = "remote-llm"
			}
			res, err := resources.Load(root, profile)
			if err != nil {
				return err
			}

			kits := agent.Kits
			if len(kits) == 0 {
				kits = cat.Defaults.Kits
			}
			kitPaths := make([]string, 0, len(kits))
			for _, k := range kits {
				kitPaths = append(kitPaths, filepath.Join(root, "kits", k))
			}

			extra := extractPassthrough(os.Args)
			if clone && !containsFlag(extra, "--clone") {
				extra = append([]string{"--clone"}, extra...)
			}

			overrideKey := "SBX_" + strings.ToUpper(agentName) + "_TEMPLATE"
			return run.Sbx(run.Opts{
				Root:             root,
				SbxAgent:         agent.SbxAgent,
				ImageName:        agent.ImageName,
				TemplateFallback: agent.TemplateFallback,
				TemplateOverride: os.Getenv(overrideKey),
				KitPaths:         kitPaths,
				ProjectDir:       projectDir,
				Extra:            extra,
				Resources:        res,
				ResourcesProfile: profile,
			})
		},
	}

	cmd.Flags().StringVar(&resourcesProfile, "resources", "", "resource profile (remote-llm|local-llm); default from catalog")
	cmd.Flags().BoolVar(&clone, "clone", false, "sandbox clone mode (isolates the host working tree)")
	// Keep unknown flags after `--` for sbx; do not parse them as cobra flags.
	cmd.Flags().SetInterspersed(false)
	cmd.DisableFlagsInUseLine = false

	return cmd
}

// extractPassthrough returns argv after the first bare "--".
func extractPassthrough(argv []string) []string {
	for i, a := range argv {
		if a == "--" && i+1 < len(argv) {
			return append([]string{}, argv[i+1:]...)
		}
	}
	return nil
}

func containsFlag(args []string, flag string) bool {
	return slices.Contains(args, flag)
}
