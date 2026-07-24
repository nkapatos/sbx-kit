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
	"github.com/nkapatos/sbx-kit/cli/internal/sbxname"
	"github.com/nkapatos/sbx-kit/cli/internal/sbxutil"
	"github.com/nkapatos/sbx-kit/cli/internal/statexfer"
)

func newRunCmd() *cobra.Command {
	var (
		agent            string
		projectPath      string
		attachName       string
		sandboxName      string
		resourcesProfile string
		clone            bool
		restoreState     bool
		yes              bool
	)

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Create a recipe sandbox or attach an existing one",
		Long: `Lifecycle helper for custom template/kit recipes (examples ship here; bring your own tree).

Intents (do not mix):

  sbx-kit run
      Attach the sole sandbox bound to the current directory.

  sbx-kit run --agent <recipe> [--path <dir>]
      CREATE only. Friendly sbx name defaults to the project dirname (like sbx).
      Interactive: ask for a custom name, then confirm create.
      --yes: skip prompts; use dirname (or --sandbox-name).
      Opaque profile id stays in the host vault — not the sbx name.

  sbx-kit run --name <sandbox>
      ATTACH only by friendly sbx name (what sbx ls shows).

Pass-through after -- goes to sbx:
  sbx-kit run --agent cursor --yes -- --memory 8g`,
		Example: `  sbx-kit run
  sbx-kit run --agent cursor --yes
  sbx-kit run --agent cursor --sandbox-name my-mocks --yes
  sbx-kit run --agent cursor --path ~/my-project --yes --restore-state
  sbx-kit run --name my-project
  sbx-kit run --agent cursor --yes -- --memory 8g`,
		Args: cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			extra := extractPassthrough(os.Args)
			if len(args) > 0 && len(extra) == 0 {
				return fmt.Errorf("unexpected arguments %v\n  use: sbx-kit run\n       sbx-kit run --agent <recipe> [--path <dir>]\n       sbx-kit run --name <sandbox>", args)
			}

			if attachName != "" && agent != "" {
				return fmt.Errorf("use either --name (attach) or --agent (create recipe), not both")
			}
			if attachName != "" && sandboxName != "" {
				return fmt.Errorf("--sandbox-name is create-only; use --name to attach")
			}

			if attachName != "" {
				if cmd.Flags().Changed("path") {
					return fmt.Errorf("--name is attach-only; do not pass --path")
				}
				if clone {
					return fmt.Errorf("--name is attach-only; --clone is a create-time flag")
				}
				if yes {
					return fmt.Errorf("--name is attach-only; --yes is only for create")
				}
				if resourcesProfile != "" {
					return fmt.Errorf("--name is attach-only; --resources applies at create")
				}
				return runAttachOnly(attachName, restoreState)
			}

			if agent != "" {
				if projectPath == "" {
					projectPath = "."
				}
				if sandboxName != "" && !sbxname.Valid(sandboxName) {
					return fmt.Errorf("invalid --sandbox-name %q", sandboxName)
				}
				return runCreateRecipe(agent, projectPath, sandboxName, resourcesProfile, clone, restoreState, yes, extra)
			}

			if projectPath == "" {
				projectPath = "."
			}
			if clone || yes || resourcesProfile != "" || sandboxName != "" {
				return fmt.Errorf("bare run is attach-only; pass --agent <recipe> to create")
			}
			rec, err := soleBindingRecord(projectPath)
			if err != nil {
				return err
			}
			return runAttachOnly(rec.SandboxName, restoreState)
		},
	}

	cmd.Flags().StringVar(&agent, "agent", "", "catalog recipe id (CREATE only)")
	cmd.Flags().StringVar(&projectPath, "path", ".", "project directory (create, or bare-run attach filter)")
	cmd.Flags().StringVar(&attachName, "name", "", "friendly sbx sandbox name (ATTACH only)")
	cmd.Flags().StringVar(&sandboxName, "sandbox-name", "", "create-time friendly sbx name (default: project dirname)")
	cmd.Flags().StringVar(&resourcesProfile, "resources", "", "resource profile at create (remote-llm|local-llm)")
	cmd.Flags().BoolVar(&clone, "clone", false, "create-time: sandbox clone mode")
	cmd.Flags().BoolVar(&restoreState, "restore-state", false, "import host profile archive after create / before attach")
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "create without prompts (dirname or --sandbox-name)")

	return cmd
}

func runCreateRecipe(agent, projectPath, sandboxName, resourcesProfile string, clone, restoreState, yes bool, extra []string) error {
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
		return fmt.Errorf("unknown recipe %q (try: sbx-kit agents)", agent)
	}
	if ag.Stub {
		return fmt.Errorf("recipe %q is still a stub in config/agents.yaml", agent)
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
		CreateOnly:       true,
		SandboxName:      sandboxName,
		ConfirmCreate:    !yes,
	}
	if !yes {
		opts.ConfirmFn = promptCreate
		if sandboxName == "" {
			opts.NamePromptFn = promptSandboxName
		}
		opts.StaleArchiveFn = promptStaleArchive
	}
	_, err = run.Sbx(opts)
	return err
}

func runAttachOnly(name string, restoreState bool) error {
	r := sbxutil.Default()
	if _, err := r.LookPath(); err != nil {
		return fmt.Errorf("sbx not found on PATH")
	}

	sbxName := name
	profileID := name
	if rec, err := binding.GetBySandbox(name); err != nil {
		return err
	} else if rec != nil {
		sbxName = rec.SandboxName
		profileID = rec.ProfileID
		fmt.Printf("==> attaching %s  recipe=%s  project=%s  profile=%s\n",
			sbxName, rec.Agent, rec.ProjectDir, profileID)
	}

	exists, err := r.Exists(sbxName)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("sandbox %q not found; try: sbx-kit status  or  sbx ls", sbxName)
	}
	if rec, _ := binding.GetBySandbox(name); rec == nil {
		fmt.Printf("==> attaching %s (no sbx-kit binding)\n", sbxName)
	}

	if restoreState {
		has, err := statexfer.HasArchive(profileID)
		if err != nil {
			return err
		}
		if !has {
			fmt.Printf("==> warning: --restore-state set but no archive for profile %s\n", profileID)
		} else if err := statexfer.Import(r, sbxName, profileID); err != nil {
			return err
		}
	}
	return r.RunInteractive("run", "--name", sbxName)
}

func soleBindingRecord(projectPath string) (*binding.Record, error) {
	abs, err := filepath.Abs(projectPath)
	if err != nil {
		return nil, err
	}
	recs, err := binding.ListForProject(abs)
	if err != nil {
		return nil, err
	}
	switch len(recs) {
	case 0:
		return nil, fmt.Errorf("no sandbox bound to %s\n  create:  sbx-kit run --agent <recipe> --path %s\n  list:    sbx-kit agents  /  sbx-kit status", abs, projectPath)
	case 1:
		rec := recs[0]
		return &rec, nil
	default:
		var b strings.Builder
		fmt.Fprintf(&b, "multiple recipes bound to %s; attach with --name:\n", abs)
		for _, rec := range recs {
			fmt.Fprintf(&b, "  --name %s  (recipe %s)\n", rec.SandboxName, rec.Agent)
		}
		return nil, fmt.Errorf("%s", strings.TrimSuffix(b.String(), "\n"))
	}
}

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

func promptSandboxName(defaultName string) (string, error) {
	fmt.Printf("Sandbox name? [%s] ", defaultName)
	line, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return "", err
	}
	ans := strings.TrimSpace(line)
	if ans == "" {
		return defaultName, nil
	}
	if !sbxname.Valid(ans) {
		return "", fmt.Errorf("invalid sandbox name %q", ans)
	}
	return ans, nil
}

func promptCreate(agent, path, name string) (bool, error) {
	fmt.Printf("Create sandbox name=%s recipe=%s\n  path=%s\nCreate? [y/N] ", name, agent, path)
	line, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return false, err
	}
	ans := strings.TrimSpace(strings.ToLower(line))
	return ans == "y" || ans == "yes", nil
}

func promptStaleArchive(profileID string) (restore, discard bool, err error) {
	fmt.Printf("Saved state exists for profile %s\n  [r]estore  [d]iscard  [n] abort  [r/d/N] ", profileID)
	line, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return false, false, err
	}
	ans := strings.TrimSpace(strings.ToLower(line))
	switch ans {
	case "r", "restore", "y", "yes":
		return true, false, nil
	case "d", "discard":
		return false, true, nil
	default:
		return false, false, nil
	}
}
