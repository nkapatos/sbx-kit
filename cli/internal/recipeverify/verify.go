package recipeverify

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/nkapatos/sbx-kit/cli/internal/catalog"
	"github.com/nkapatos/sbx-kit/cli/internal/sbxcompat"
)

// KitRunner runs sbx kit verify (tests inject a fake).
type KitRunner interface {
	KitVerify(path string) error
	EnsureKitVerify() error
}

// Options controls recipe verification.
type Options struct {
	Out      io.Writer
	SkipKits bool
	Runner   KitRunner
}

func (o *Options) out() io.Writer {
	if o.Out != nil {
		return o.Out
	}
	return os.Stdout
}

// Describe returns help text for recipes verify commands.
func Describe() string {
	return `Recipe manifests are checked by sbx-kit.
Kit specs are checked by sbx — sbx-kit runs sbx kit verify on kit paths.
Migrate kits with sbx, not sbx-kit.`
}

// VerifyRecipe checks agents.yaml and references; optionally delegates kit verify.
func VerifyRecipe(catalogRoot, id string, opts Options) error {
	targets, err := recipeTargets(catalogRoot, id)
	if err != nil {
		return err
	}
	var errs []string
	kitPaths := map[string]struct{}{}

	for _, t := range targets {
		if err := verifyOneRecipe(t, opts); err != nil {
			errs = append(errs, err.Error())
			continue
		}
		for _, p := range t.kitPaths {
			kitPaths[p] = struct{}{}
		}
	}

	if !opts.SkipKits && len(kitPaths) > 0 {
		if err := verifyKitPaths(sortedKeys(kitPaths), opts); err != nil {
			errs = append(errs, err.Error())
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("verify failed:\n  %s", strings.Join(errs, "\n  "))
	}
	return nil
}

// VerifyKits runs sbx kit verify on kits under catalog directory dir (empty = all).
func VerifyKits(catalogRoot, dir string, opts Options) error {
	dirs, err := catalog.List(catalogRoot)
	if err != nil {
		return err
	}
	dirs, err = catalog.FilterDirs(dirs, dir)
	if err != nil {
		return err
	}
	if len(dirs) == 0 {
		fmt.Fprintln(opts.out(), "(no directories)")
		return nil
	}

	paths := map[string]struct{}{}
	for _, d := range dirs {
		found, err := catalog.ListKitPaths(d.Root)
		if err != nil {
			return err
		}
		for _, p := range found {
			paths[p] = struct{}{}
		}
	}
	if len(paths) == 0 {
		fmt.Fprintln(opts.out(), "(no kits under kits/)")
		return nil
	}
	return verifyKitPaths(sortedKeys(paths), opts)
}

type recipeTarget struct {
	id       string
	dir      catalog.Dir
	manifest *catalog.Manifest
	agent    catalog.Agent
	kitPaths []string
}

func recipeTargets(catalogRoot, id string) ([]recipeTarget, error) {
	if id != "" {
		d, manifest, agent, err := catalog.Lookup(catalogRoot, id)
		if err != nil {
			return nil, err
		}
		kits := catalog.ResolveKits(agent.Kits, manifest.Defaults.Kits)
		return []recipeTarget{{
			id:       id,
			dir:      d,
			manifest: manifest,
			agent:    agent,
			kitPaths: catalog.KitPaths(d.Root, kits),
		}}, nil
	}

	dirs, err := catalog.List(catalogRoot)
	if err != nil {
		return nil, err
	}
	if len(dirs) == 0 {
		fmt.Fprintln(os.Stdout, "(no directories)")
		return nil, nil
	}

	var out []recipeTarget
	for _, d := range dirs {
		manifest, err := catalog.Load(catalog.File(d.Root))
		if err != nil {
			return nil, err
		}
		names := make([]string, 0, len(manifest.Agents))
		for name := range manifest.Agents {
			names = append(names, name)
		}
		sort.Strings(names)
		for _, name := range names {
			agent := manifest.Agents[name]
			kits := catalog.ResolveKits(agent.Kits, manifest.Defaults.Kits)
			out = append(out, recipeTarget{
				id:       catalog.JoinID(d.Name, name),
				dir:      d,
				manifest: manifest,
				agent:    agent,
				kitPaths: catalog.KitPaths(d.Root, kits),
			})
		}
	}
	return out, nil
}

func verifyOneRecipe(t recipeTarget, opts Options) error {
	if strings.TrimSpace(t.agent.SbxAgent) == "" && !t.agent.Stub {
		return fmt.Errorf("recipe %s: missing sbx_agent", t.id)
	}
	for _, kp := range t.kitPaths {
		if err := catalog.RequireKitPath(kp); err != nil {
			return fmt.Errorf("recipe %s: %w", t.id, err)
		}
	}
	fmt.Fprintf(opts.out(), "==> recipe %s: ok\n", t.id)
	return nil
}

func verifyKitPaths(paths []string, opts Options) error {
	if opts.Runner == nil {
		return fmt.Errorf("kit verify: no sbx runner configured")
	}
	if err := opts.Runner.EnsureKitVerify(); err != nil {
		return err
	}

	var errs []string
	for _, p := range paths {
		name := filepath.Base(p)
		fmt.Fprintf(opts.out(), "==> kit %s: verify with sbx\n", name)
		if err := opts.Runner.KitVerify(p); err != nil {
			errs = append(errs, fmt.Sprintf("kit %s: %v\n  run: sbx kit verify %s", name, err, p))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("%s", strings.Join(errs, "\n"))
	}
	return nil
}

func sortedKeys(m map[string]struct{}) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// SbxKitRunner delegates kit verify to the sbx CLI.
type SbxKitRunner struct {
	ProbeVersion func() (string, error)
	KitVerifyFn  func(path string) error
}

func (r SbxKitRunner) EnsureKitVerify() error {
	probe := r.ProbeVersion
	if probe == nil {
		return fmt.Errorf("sbx not available\n  install Docker sbx >= %s", sbxcompat.MinKitVerify)
	}
	return sbxcompat.EnsureFeature(probe, sbxcompat.MinKitVerify, "sbx kit verify")
}

func (r SbxKitRunner) KitVerify(path string) error {
	if r.KitVerifyFn == nil {
		return fmt.Errorf("sbx kit verify not configured")
	}
	return r.KitVerifyFn(path)
}
