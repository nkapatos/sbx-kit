package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/nkapatos/sbx-kit/cli/internal/catalog"
	"github.com/nkapatos/sbx-kit/cli/internal/kitcreds"
	"github.com/nkapatos/sbx-kit/cli/internal/sbxutil"
)

func newCheckCmd() *cobra.Command {
	var agent, path, name string

	cmd := &cobra.Command{
		Use:   "check",
		Short: "Diagnostics for a recipe/sandbox (bindings, declared secrets, sbx secret ls)",
		Long: `Convenience diagnostics. Resolves the sandbox from --name, --agent/--path,
or the sole binding for the current directory, then:

  • shows binding / recipe identity
  • lists credential services declared by the recipe kits
  • runs sbx secret ls --sandbox <name> when the box exists (else sbx secret ls)

Does not store or validate secrets — sbx owns that.`,
		Example: `  sbx-kit check
  sbx-kit check --name my-project
  sbx-kit check --agent shell-hub --path ~/proj`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCheck(agent, path, name)
		},
	}

	cmd.Flags().StringVar(&agent, "agent", "", "catalog recipe (with --path)")
	cmd.Flags().StringVar(&path, "path", ".", "project directory")
	cmd.Flags().StringVar(&name, "name", "", "friendly sbx sandbox name")
	return cmd
}

func runCheck(agent, path, name string) error {
	rs, err := resolveFlags(agent, path, name)
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
	recipe := rs.AgentName
	if len(kitPaths) == 0 && rs.AgentName != "" && rs.Root != "" {
		// resolveSandboxArg by name alone may omit kits — reload from catalog.
		cat, err := catalog.Load(filepath.Join(rs.Root, "config", "agents.yaml"))
		if err == nil {
			if ag, ok := cat.Agents[rs.AgentName]; ok {
				kits := ag.Kits
				if len(kits) == 0 {
					kits = cat.Defaults.Kits
				}
				kitPaths = make([]string, 0, len(kits))
				for _, k := range kits {
					kitPaths = append(kitPaths, filepath.Join(rs.Root, "kits", k))
				}
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
	} else if recipe == "" {
		fmt.Println()
		fmt.Println("==> no recipe kits loaded (name-only binding); skip declared-secret scan")
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
