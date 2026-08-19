package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/nkapatos/sbx-kit/cli/internal/catalog"
	"github.com/nkapatos/sbx-kit/cli/internal/kitcreds"
	"github.com/nkapatos/sbx-kit/cli/internal/sbxutil"
)

func newCheckCmd() *cobra.Command {
	var recipe, path, name string

	cmd := &cobra.Command{
		Use:   "check",
		Short: "Sandbox/recipe info and sbx secret ls",
		Long: `Resolve via --name, --recipe/--path, or the sole cwd binding, then show
identity, credential services declared by recipe kits, and run
sbx secret ls (--sandbox when the box exists).`,
		Example: `  sbx-kit check
  sbx-kit check --name my-project
  sbx-kit check --recipe mine/shell --path ~/proj`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCheck(recipe, path, name)
		},
	}

	addTargetFlags(cmd, &recipe, &path, &name)
	return cmd
}

func runCheck(recipeID, path, name string) error {
	rs, err := resolveFlags(recipeID, path, name)
	if err != nil {
		return err
	}

	fmt.Printf("==> check\n")
	fmt.Printf("  sandbox=%s\n", rs.SandboxName)
	if rs.AgentName != "" {
		fmt.Printf("  recipe=%s\n", rs.AgentName)
	}
	if rs.ProfileID != "" {
		fmt.Printf("  profile=%s\n", rs.ProfileID)
	}
	if rs.ProjectDir != "" {
		fmt.Printf("  project=%s\n", rs.ProjectDir)
	}
	if rs.SbxAgent != "" {
		fmt.Printf("  sbx_agent=%s\n", rs.SbxAgent)
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
			return err
		}
		fmt.Println()
		if len(needs) == 0 {
			fmt.Println("==> recipe kits declare no credential services")
		} else {
			fmt.Println("==> credential services declared by recipe kits")
			for _, n := range needs {
				line := fmt.Sprintf("  %s", n.Service)
				if len(n.Envs) > 0 {
					line += fmt.Sprintf("  env=%v", n.Envs)
				}
				if n.KitName != "" {
					line += fmt.Sprintf("  kit=%s", n.KitName)
				}
				fmt.Println(line)
			}
			fmt.Println("  set on host:  sbx secret set <service>")
		}
	} else if rs.AgentName == "" {
		fmt.Println()
		fmt.Println("==> no recipe kits loaded; skip declared-secret scan")
	}

	r := sbxutil.Default()
	if _, err := r.LookPath(); err != nil {
		fmt.Printf("\n==> sbx not on PATH; skip secret ls\n")
		return nil
	}

	exists, err := r.Exists(rs.SandboxName)
	if err != nil {
		fmt.Printf("\n==> warning: could not query sandboxes: %v\n", err)
		exists = false
	}

	fmt.Println()
	if exists {
		fmt.Printf("==> sbx secret ls --sandbox %s\n", rs.SandboxName)
		if err := r.RunInteractive("secret", "ls", "--sandbox", rs.SandboxName); err != nil {
			fmt.Fprintf(os.Stderr, "sbx secret ls failed: %v\n", err)
		}
	} else {
		fmt.Printf("==> sandbox %s not in sbx ls; showing global secrets\n", rs.SandboxName)
		fmt.Println("==> sbx secret ls")
		if err := r.RunInteractive("secret", "ls"); err != nil {
			fmt.Fprintf(os.Stderr, "sbx secret ls failed: %v\n", err)
		}
	}
	return nil
}
