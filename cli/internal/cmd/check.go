package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/nkapatos/sbx-kit/cli/internal/catalog"
	"github.com/nkapatos/sbx-kit/cli/internal/kitcreds"
	"github.com/nkapatos/sbx-kit/cli/internal/sbxutil"
	"github.com/nkapatos/sbx-kit/cli/internal/ui"
)

func newCheckCmd() *cobra.Command {
	var recipe, path, name string

	cmd := &cobra.Command{
		Use:   "check",
		Short: "Sandbox/recipe info and sbx secret ls",
		Long: `Resolve via --name, --recipe/--path, or the sole cwd binding, then show
identity, credential services declared by recipe kits, and run
sbx secret ls (--sandbox when the box exists).`,
		Example: `  sbx-kit box check
  sbx-kit box check --name my-project
  sbx-kit box check --recipe mine/shell --path ~/proj`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCheck(recipe, path, name)
		},
	}

	addTargetFlags(cmd, &recipe, &path, &name)
	return cmd
}

// CheckResult is the data for box check. A TUI can consume this without rendering.
type CheckResult struct {
	SandboxName string
	RecipeID    string
	Profile     string
	Project     string
	SbxAgent    string
	CredNeeds   []kitcreds.Need
	SkipCreds   bool
	SbxOnPath   bool
	SbxExists   bool
	LsErr       error
}

func computeCheck(recipeID, path, name string) (*CheckResult, error) {
	rs, err := resolveFlags(recipeID, path, name)
	if err != nil {
		return nil, err
	}
	res := &CheckResult{
		SandboxName: rs.SandboxName,
		RecipeID:    rs.AgentName,
		Profile:     rs.ProfileID,
		Project:     rs.ProjectDir,
		SbxAgent:    rs.SbxAgent,
	}

	kitPaths := rs.KitPaths
	if len(kitPaths) == 0 && rs.AgentName != "" {
		catalogRoot := rs.Catalog
		if catalogRoot == "" {
			catalogRoot, _ = requireToolkitRoot()
		}
		if catalogRoot != "" {
			if src, manifest, ag, err := catalog.Lookup(catalogRoot, rs.AgentName); err == nil {
				kits := catalog.ResolveKits(ag.Kits, manifest.Defaults.Kits)
				kitPaths = catalog.KitPaths(src.Root, kits)
			}
		}
	}
	if len(kitPaths) > 0 {
		needs, err := kitcreds.ScanSpecs(kitPaths)
		if err != nil {
			return nil, err
		}
		res.CredNeeds = needs
	} else if rs.AgentName == "" {
		res.SkipCreds = true
	}

	r := sbxRunner()
	if _, err := r.LookPath(); err != nil {
		res.SbxOnPath = false
		return res, nil
	}
	res.SbxOnPath = true
	exists, err := r.Exists(rs.SandboxName)
	if err != nil {
		res.LsErr = err
		return res, nil
	}
	res.SbxExists = exists
	return res, nil
}

func renderCheck(w *ui.Writer, res *CheckResult) {
	w.Header("check")
	w.Detail("sandbox", res.SandboxName)
	w.Detail("recipe", res.RecipeID)
	w.Detail("profile", res.Profile)
	w.Detail("project", res.Project)
	w.Detail("sbx_agent", res.SbxAgent)
	w.Println()

	switch {
	case res.SkipCreds:
		w.Header("no recipe kits loaded; skip declared-secret scan")
	case len(res.CredNeeds) == 0:
		w.Header("recipe kits declare no credential services")
	default:
		w.Header("credential services declared by recipe kits")
		rows := make([][]string, 0, len(res.CredNeeds))
		for _, n := range res.CredNeeds {
			rows = append(rows, []string{n.Service, strings.Join(n.Envs, ","), n.KitName})
		}
		_ = w.Table([]string{"SERVICE", "ENV", "KIT"}, rows)
		w.Println("  set on host:  sbx secret set <service>")
	}

	if !res.SbxOnPath {
		w.Println()
		w.Header("sbx not on PATH; skip secret ls")
		return
	}
	if res.LsErr != nil {
		w.Warn(fmt.Sprintf("could not query sandboxes: %v", res.LsErr))
	}
}

func runCheck(recipeID, path, name string) error {
	res, err := computeCheck(recipeID, path, name)
	if err != nil {
		return err
	}
	renderCheck(UI(), res)
	if !res.SbxOnPath {
		return nil
	}

	r := sbxRunner()
	w := UI()
	w.Println()
	if res.SbxExists {
		w.Header("sbx secret ls --sandbox " + res.SandboxName)
		if err := r.RunInteractive("secret", "ls", "--sandbox", res.SandboxName); err != nil {
			w.Warn("sbx secret ls failed: " + err.Error())
		}
		return nil
	}
	w.Header("sandbox " + res.SandboxName + " not in sbx ls; showing global secrets")
	w.Header("sbx secret ls")
	if err := r.RunInteractive("secret", "ls"); err != nil {
		w.Warn("sbx secret ls failed: " + err.Error())
	}
	return nil
}

func sbxRunner() *sbxutil.Runner {
	r := sbxutil.Default()
	r.Out = UI().Out
	return r
}
