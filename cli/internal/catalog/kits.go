package catalog

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// RequireKitPath reports whether a kit directory exists with spec.yaml.
func RequireKitPath(kitPath string) error {
	spec := filepath.Join(kitPath, "spec.yaml")
	st, err := os.Stat(spec)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("missing kit %q (no spec.yaml)", kitPath)
		}
		return err
	}
	if st.IsDir() {
		return fmt.Errorf("invalid kit %q (spec.yaml is a directory)", kitPath)
	}
	return nil
}

// ListKitPaths returns kit directories under dirRoot/kits that contain spec.yaml.
func ListKitPaths(dirRoot string) ([]string, error) {
	kitsRoot := filepath.Join(dirRoot, "kits")
	ents, err := os.ReadDir(kitsRoot)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var out []string
	for _, e := range ents {
		if !e.IsDir() || e.Name() == "" || e.Name()[0] == '.' {
			continue
		}
		p := filepath.Join(kitsRoot, e.Name())
		if err := RequireKitPath(p); err != nil {
			continue
		}
		out = append(out, p)
	}
	sort.Strings(out)
	return out, nil
}

// ListAllKitPaths returns every kit path under all directories in the catalog.
func ListAllKitPaths(catalogRoot string) ([]string, error) {
	dirs, err := List(catalogRoot)
	if err != nil {
		return nil, err
	}
	seen := map[string]struct{}{}
	var out []string
	for _, d := range dirs {
		paths, err := ListKitPaths(d.Root)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", d.Name, err)
		}
		for _, p := range paths {
			if _, ok := seen[p]; ok {
				continue
			}
			seen[p] = struct{}{}
			out = append(out, p)
		}
	}
	sort.Strings(out)
	return out, nil
}
