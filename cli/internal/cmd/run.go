package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/spf13/cobra"

	"github.com/nkapatos/sbx-kit/cli/internal/binding"
	"github.com/nkapatos/sbx-kit/cli/internal/catalog"
	"github.com/nkapatos/sbx-kit/cli/internal/resources"
	"github.com/nkapatos/sbx-kit/cli/internal/run"
	"github.com/nkapatos/sbx-kit/cli/internal/sbxutil"
	"github.com/nkapatos/sbx-kit/cli/internal/statexfer"
)

func newRunCmd() *cobra.Command {
	var (
		agent            string
		projectPath      string
		sandboxName      string
		resourcesProfile string
		clone            bool
		restoreState     bool
		yes              bool
	)

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Attach or create a catalog agent sandbox for a project",
		Long: `Project-scoped launcher for custom templates/kits.

  sbx-kit run --agent cursor
      Attach the bound sandbox for the current directory, or prompt to create.

  sbx-kit run --agent cursor --path ~/proj
      Same for an explicit project path.

  sbx-kit run --name sbxk-cursor-deadbeef
      Attach an existing sandbox from anywhere (no create).

--agent selects the catalog recipe (template + kits) used at create time.
--name is attach-only. Opaque sandbox ids stay internal; see sbx-kit status.

Pass-through: anything after -- goes to sbx unchanged:
  sbx-kit run --agent cursor -- --memory 8g`,
		Example: `  sbx-kit run --agent cursor
  sbx-kit run --agent cursor --path ~/my-project
  sbx-kit run --agent cursor --yes
  sbx-kit run --agent cursor --clone --restore-state
  sbx-kit run --name sbxk-cursor-deadbeef
  sbx-kit run --agent cursor -- --memory 8g`,
		Args: cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			extra := extractPassthrough(os.Args)
			if len(args) > 0 && len(extra) == 0 {
				return fmt.Errorf("unexpected arguments %v\n  use: sbx-kit run --agent <name> [--path <dir>]\n  or:  sbx-kit run --name <sandbox>", args)
			}

			if sandboxName != "" && agent != "" {
				return fmt.Errorf("use either --name (attach) or --agent (project recipe), not both")
			}
			if sandboxName != "" {
				return runAttachOnly(sandboxName, restoreState)
			}

			if projectPath == "" {
				projectPath = "."
			}

			if agent == "" {
				resolved, err := resolveSoleBinding(projectPath)
				if err != nil {
					return err
				}
				agent = resolved
			}

			root, err := requireToolkitRoot()
			if err != nil {
				return err
			}
			cat, err := catalog.Load(filepath.Join(root, "config", "agents.yaml"))
			if err != nil {
				return err
			}
			ag, ok := cat.Agents[agent]
			if !ok {
				return fmt.Errorf("unknown agent %q (try: sbx-kit agents)", agent)
			}
			if ag.Stub {
				return fmt.Errorf("agent %q is still a stub in config/agents.yaml", agent)
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

			kits := ag.Kits
			if len(kits) == 0 {
				kits = cat.Defaults.Kits
			}
			kitPaths := make([]string, 0, len(kits))
			for _, k := range kits {
				kitPaths = append(kitPaths, filepath.Join(root, "kits", k))
			}

			if clone && !containsFlag(extra, "--clone") {
				extra = append([]string{"--clone"}, extra...)
			}

			overrideKey := "SBX_" + strings.ToUpper(agent) + "_TEMPLATE"
			opts := run.Opts{
				Root:             root,
				AgentCatalogName: agent,
				SbxAgent:         ag.SbxAgent,
				ImageName:        ag.ImageName,
				TemplateFallback: ag.TemplateFallback,
				TemplateOverride: os.Getenv(overrideKey),
				KitPaths:         kitPaths,
				ProjectDir:       projectPath,
				Extra:            extra,
				Resources:        res,
				ResourcesProfile: profile,
				RestoreState:     restoreState,
				ConfirmCreate:    !yes,
			}
			if !yes {
				opts.ConfirmFn = promptCreate
			}
			_, err = run.Sbx(opts)
			return err
		},
	}

	cmd.Flags().StringVar(&agent, "agent", "", "catalog agent recipe (create/attach by project binding)")
	cmd.Flags().StringVar(&projectPath, "path", ".", "project directory")
	cmd.Flags().StringVar(&sandboxName, "name", "", "attach existing sandbox by name (no create)")
	cmd.Flags().StringVar(&resourcesProfile, "resources", "", "resource profile (remote-llm|local-llm); default from catalog")
	cmd.Flags().BoolVar(&clone, "clone", false, "sandbox clone mode (isolates the host working tree)")
	cmd.Flags().BoolVar(&restoreState, "restore-state", false, "import host profile archive into the sandbox before attach")
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "create without confirmation when no sandbox exists")

	return cmd
}

func runAttachOnly(name string, restoreState bool) error {
	r := sbxutil.Default()
	if _, err := r.LookPath(); err != nil {
		return fmt.Errorf("sbx not found on PATH")
	}
	exists, err := r.Exists(name)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("sandbox %q not found; try: sbx-kit status  or  sbx ls", name)
	}

	profileID := name
	if rec, err := binding.GetBySandbox(name); err == nil && rec != nil {
		profileID = rec.ProfileID
		fmt.Printf("==> attaching %s  label=%s  agent=%s  project=%s\n",
			name, binding.Label(rec), rec.Agent, rec.ProjectDir)
	} else {
		fmt.Printf("==> attaching %s (no sbx-kit binding)\n", name)
	}

	if restoreState {
		has, err := statexfer.HasArchive(profileID)
		if err != nil {
			return err
		}
		if !has {
			fmt.Printf("==> warning: --restore-state set but no archive for profile %s\n", profileID)
		} else if err := statexfer.Import(r, name, profileID); err != nil {
			return err
		}
	}
	return r.RunInteractive("run", "--name", name)
}

func resolveSoleBinding(projectPath string) (string, error) {
	abs, err := filepath.Abs(projectPath)
	if err != nil {
		return "", err
	}
	recs, err := binding.ListForProject(abs)
	if err != nil {
		return "", err
	}
	switch len(recs) {
	case 0:
		return "", fmt.Errorf("no sandbox bound to %s\n  create with: sbx-kit run --agent <name> --path %s\n  list:       sbx-kit agents  /  sbx-kit status", abs, projectPath)
	case 1:
		return recs[0].Agent, nil
	default:
		var b strings.Builder
		fmt.Fprintf(&b, "multiple agents bound to %s; pass --agent explicitly:\n", abs)
		for _, rec := range recs {
			fmt.Fprintf(&b, "  --agent %s  (sandbox %s)\n", rec.Agent, rec.SandboxName)
		}
		return "", fmt.Errorf("%s", strings.TrimSuffix(b.String(), "\n"))
	}
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

// promptCreate is the default interactive confirm used by run.Sbx.
func promptCreate(agent, path, name string) (bool, error) {
	fmt.Printf("No sandbox for agent=%s path=%s\nCreate %s? [y/N] ", agent, path, name)
	line, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return false, err
	}
	ans := strings.TrimSpace(strings.ToLower(line))
	return ans == "y" || ans == "yes", nil
}
