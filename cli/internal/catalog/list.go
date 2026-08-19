package catalog

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Dir is one recipe directory under the catalog (recipes/, kits/, images/).
type Dir struct {
	Name string
	Root string
}

// IsDir reports whether dir looks like a recipe directory.
func IsDir(dir string) bool {
	st, err := os.Stat(File(dir))
	return err == nil && !st.IsDir()
}

// List returns recipe directories under the catalog, sorted by name.
func List(catalogRoot string) ([]Dir, error) {
	ents, err := os.ReadDir(catalogRoot)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var out []Dir
	for _, e := range ents {
		if !e.IsDir() || strings.HasPrefix(e.Name(), ".") {
			continue
		}
		root := filepath.Join(catalogRoot, e.Name())
		if !IsDir(root) {
			continue
		}
		out = append(out, Dir{Name: e.Name(), Root: root})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}

// Lookup loads <dir>/<recipe> from the catalog.
func Lookup(catalogRoot, recipeID string) (d Dir, manifest *Manifest, agent Agent, err error) {
	dirName, rec, err := ParseID(recipeID)
	if err != nil {
		return Dir{}, nil, Agent{}, err
	}
	root := filepath.Join(catalogRoot, dirName)
	if !IsDir(root) {
		return Dir{}, nil, Agent{}, fmt.Errorf("unknown recipe %q (try: sbx-kit recipes)", recipeID)
	}
	manifest, err = Load(File(root))
	if err != nil {
		return Dir{}, nil, Agent{}, err
	}
	agent, ok := manifest.Agents[rec]
	if !ok {
		return Dir{}, nil, Agent{}, fmt.Errorf("unknown recipe %q (try: sbx-kit recipes)", recipeID)
	}
	return Dir{Name: dirName, Root: root}, manifest, agent, nil
}

// FilterDirs returns dirs named filter, or all dirs when filter is empty.
func FilterDirs(dirs []Dir, filter string) ([]Dir, error) {
	if filter == "" {
		return dirs, nil
	}
	for _, d := range dirs {
		if d.Name == filter {
			return []Dir{d}, nil
		}
	}
	return nil, fmt.Errorf("unknown directory %q (try: sbx-kit catalog ls)", filter)
}
