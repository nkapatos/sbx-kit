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
		Short: "Create a recipe sandbox or attach an existing one",
		Long: `Lifecycle helper for custom template/kit recipes (examples ship here; bring your own tree).

Intents (do not mix):

  sbx-kit run
      Attach the sole sandbox bound to the current directory.

  sbx-kit run --agent <recipe> [--path <dir>]
      CREATE only: prompt (or --yes) to create from the catalog recipe.
      If that project+recipe sandbox already exists, errors — use --name or bare run.

  sbx-kit run --name <sandbox>
      ATTACH only from anywhere (no --path / create flags).

--agent is a recipe id (template + kits + sbx_agent), not "attach".
Pass-through after -- goes to sbx:
  sbx-kit run --agent cursor --yes -- --memory 8g`,
		Example: `  sbx-kit run
  sbx-kit run --agent cursor --yes
  sbx-kit run --agent cursor --path ~/my-project --yes
  sbx-kit run --agent cursor --yes --restore-state
  sbx-kit run --name sbxk-cursor-deadbeef
  sbx-kit run --agent cursor --yes -- --memory 8g`,
		Args: cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			extra := extractPassthrough(os.Args)
			if len(args) > 0 && len(extra) == 0 {
				return fmt.Errorf("unexpected arguments %v\n  use: sbx-kit run\n       sbx-kit run --agent <recipe> [--path <dir>]\n       sbx-kit run --name <sandbox>", args)
			}

			if sandboxName != "" && agent != "" {
				return fmt.Errorf("use either --name (attach) or --agent (create recipe), not both")
			}

			// --- attach by name ---
			if sandboxName != "" {
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
				return runAttachOnly(sandboxName, restoreState)
			}

			// --- create by recipe ---
			if agent != "" {
				if projectPath == "" {
					projectPath = "."
				}
				return runCreateRecipe(agent, projectPath, resourcesProfile, clone, restoreState, yes, extra)
			}

			// --- bare run: attach sole cwd binding ---
			if projectPath == "" {
				projectPath = "."
			}
			if cmd.Flags().Changed("path") {
				// allow --path with bare run to attach sole binding for that project
			}
			if clone || yes || resourcesProfile != "" {
				return fmt.Errorf("bare run is attach-only; pass --agent <recipe> to create (with --yes/--clone/--resources)")
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
	cmd.Flags().StringVar(&sandboxName, "name", "", "sandbox id (ATTACH only)")
	cmd.Flags().StringVar(&resourcesProfile, "resources", "", "resource profile at create (remote-llm|local-llm)")
	cmd.Flags().BoolVar(&clone, "clone", false, "create-time: sandbox clone mode")
	cmd.Flags().BoolVar(&restoreState, "restore-state", false, "import host profile archive before attach/create")
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "create without confirmation")

	return cmd
}

func runCreateRecipe(agent, projectPath, resourcesProfile string, clone, restoreState, yes bool, extra []string) error {
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
		ConfirmCreate:    !yes,
	}
	if !yes {
		opts.ConfirmFn = promptCreate
	}
	_, err = run.Sbx(opts)
	return err
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
		fmt.Printf("==> attaching %s  label=%s  recipe=%s  project=%s\n",
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

func promptCreate(agent, path, name string) (bool, error) {
	fmt.Printf("Create sandbox for recipe=%s path=%s\nCreate %s? [y/N] ", agent, path, name)
	line, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return false, err
	}
	ans := strings.TrimSpace(strings.ToLower(line))
	return ans == "y" || ans == "yes", nil
}
