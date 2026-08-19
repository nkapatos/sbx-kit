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
	"github.com/nkapatos/sbx-kit/cli/internal/kitcreds"
	"github.com/nkapatos/sbx-kit/cli/internal/resources"
	"github.com/nkapatos/sbx-kit/cli/internal/run"
	"github.com/nkapatos/sbx-kit/cli/internal/sbxname"
	"github.com/nkapatos/sbx-kit/cli/internal/sbxutil"
	"github.com/nkapatos/sbx-kit/cli/internal/statexfer"
)

func newRunCmd() *cobra.Command {
	var (
		projectPath      string
		attachName       string
		sandboxName      string
		resourcesProfile string
		clone            bool
		restoreState     bool
		yes              bool
	)

	cmd := &cobra.Command{
		Use:   "run [<dir>/<name>]",
		Short: "Create from a recipe or attach a sandbox",
		Long: `Create or attach — pick one mode:

Attach
  sbx-kit run                          sole binding at --path
  sbx-kit run --name <sandbox>         attach by sbx ls name
  Flags: --restore-state

Create
  sbx-kit run <dir>/<name> [--path <dir>]
  Flags: --path --sandbox-name --yes --resources --clone --restore-state

Pass-through after -- goes to sbx.`,
		Example: `  sbx-kit run mine/cursor --yes
  sbx-kit run mine/shell --yes
  sbx-kit run mine/kit-cursor --yes
  sbx-kit run mine/shell --sandbox-name ds-creds --yes
  sbx-kit run --name my-project
  sbx-kit run mine/cursor --yes -- --memory 8g`,
		Args: cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			extra := extractPassthrough(os.Args)
			pos := positionalArgs(args, extra)
			if len(pos) > 1 {
				return fmt.Errorf("unexpected arguments %v\n  use: sbx-kit run\n       sbx-kit run <dir>/<name> [--path <dir>]\n       sbx-kit run --name <sandbox>", pos)
			}

			recipeID := ""
			if len(pos) == 1 {
				recipeID = pos[0]
			}

			if attachName != "" && recipeID != "" {
				return fmt.Errorf("use either --name (attach) or a recipe (create), not both")
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

			if recipeID != "" {
				if projectPath == "" {
					projectPath = "."
				}
				if sandboxName != "" && !sbxname.Valid(sandboxName) {
					return fmt.Errorf("invalid --sandbox-name %q", sandboxName)
				}
				return runCreateRecipe(recipeID, projectPath, sandboxName, resourcesProfile, clone, restoreState, yes, extra)
			}

			if projectPath == "" {
				projectPath = "."
			}
			if clone || yes || resourcesProfile != "" || sandboxName != "" {
				return fmt.Errorf("bare run is attach-only; pass a recipe to create (sbx-kit run <dir>/<name> --yes)")
			}
			rec, err := soleBindingRecord(projectPath)
			if err != nil {
				return err
			}
			return runAttachOnly(rec.SandboxName, restoreState)
		},
	}

	cmd.Flags().StringVar(&projectPath, "path", ".", "[create/attach] project directory")
	cmd.Flags().StringVar(&attachName, "name", "", "[attach] sandbox name from sbx ls")
	cmd.Flags().StringVar(&sandboxName, "sandbox-name", "", "[create] sandbox name (default: project dirname)")
	cmd.Flags().StringVar(&resourcesProfile, "resources", "", "[create] resource profile (remote-llm|local-llm)")
	cmd.Flags().BoolVar(&clone, "clone", false, "[create] pass --clone to sbx")
	cmd.Flags().BoolVar(&restoreState, "restore-state", false, "restore host profile archive")
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "[create] skip prompts")

	return cmd
}

func runCreateRecipe(recipeID, projectPath, sandboxName, resourcesProfile string, clone, restoreState, yes bool, extra []string) error {
	catalogRoot, err := requireToolkitRoot()
	if err != nil {
		return err
	}
	src, manifest, ag, err := catalog.Lookup(catalogRoot, recipeID)
	if err != nil {
		return err
	}
	if ag.Stub {
		return fmt.Errorf("recipe %q is still a stub in %s", recipeID, catalog.File(src.Root))
	}

	profile := resourcesProfile
	if profile == "" {
		profile = manifest.Defaults.Resources
	}
	if profile == "" {
		profile = "remote-llm"
	}
	res, err := resources.Load(src.Root, profile)
	if err != nil {
		return err
	}

	kits := catalog.ResolveKits(ag.Kits, manifest.Defaults.Kits)
	kitPaths := catalog.KitPaths(src.Root, kits)

	if clone && !containsFlag(extra, "--clone") {
		extra = append([]string{"--clone"}, extra...)
	}

	if needs, err := kitcreds.ScanSpecs(kitPaths); err != nil {
		return err
	} else if len(needs) > 0 {
		kitcreds.PrintHints(os.Stdout, recipeID, needs)
		fmt.Println()
	}

	overrideKey := templateOverrideEnv(recipeID)
	opts := run.Opts{
		Root:             src.Root,
		RecipeID:         recipeID,
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
		return fmt.Errorf("sandbox %q not found; try: sbx-kit bindings  or  sbx ls", sbxName)
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
			fmt.Printf("==> warning: --restore-state set but no archive at profile %s\n", profileID)
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
		return nil, fmt.Errorf("no sandbox bound to %s\n  create:  sbx-kit run <dir>/<name> --path %s --yes\n  list:    sbx-kit recipes  /  sbx-kit bindings", abs, projectPath)
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

func extractPassthrough(argv []string) []string {
	for i, a := range argv {
		if a == "--" && i+1 < len(argv) {
			return append([]string{}, argv[i+1:]...)
		}
	}
	return nil
}

// positionalArgs drops cobra args that came from after `--`.
func positionalArgs(args, extra []string) []string {
	if len(extra) == 0 || len(args) < len(extra) {
		return args
	}
	tail := args[len(args)-len(extra):]
	for i := range extra {
		if tail[i] != extra[i] {
			return args
		}
	}
	return args[:len(args)-len(extra)]
}

func containsFlag(args []string, flag string) bool {
	return slices.Contains(args, flag)
}
