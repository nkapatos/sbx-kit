package catalog

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Source is one bundle under the catalog (recipes/, kits/, images/).
type Source struct {
	Name string
	Root string
}

// IsDir reports whether dir looks like a source.
func IsDir(dir string) bool {
	st, err := os.Stat(File(dir))
	return err == nil && !st.IsDir()
}

// List returns sources under the catalog, sorted by name. Hidden dirs are skipped.
func List(catalogRoot string) ([]Source, error) {
	ents, err := os.ReadDir(catalogRoot)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var out []Source
	for _, e := range ents {
		if !e.IsDir() || strings.HasPrefix(e.Name(), ".") {
			continue
		}
		root := filepath.Join(catalogRoot, e.Name())
		if !IsDir(root) {
			continue
		}
		out = append(out, Source{Name: e.Name(), Root: root})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}

// Lookup loads <source>/<name> from the catalog.
func Lookup(catalogRoot, recipeID string) (src Source, manifest *Manifest, agent Agent, err error) {
	sourceName, rec, err := ParseID(recipeID)
	if err != nil {
		return Source{}, nil, Agent{}, err
	}
	root := filepath.Join(catalogRoot, sourceName)
	if !IsDir(root) {
		return Source{}, nil, Agent{}, fmt.Errorf("unknown recipe %q (try: sbx-kit recipes)", recipeID)
	}
	manifest, err = Load(File(root))
	if err != nil {
		return Source{}, nil, Agent{}, err
	}
	agent, ok := manifest.Agents[rec]
	if !ok {
		return Source{}, nil, Agent{}, fmt.Errorf("unknown recipe %q (try: sbx-kit recipes)", recipeID)
	}
	return Source{Name: sourceName, Root: root}, manifest, agent, nil
}
