package catalog

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Source is one catalog (a one-level child of the tree).
type Source struct {
	Name string
	Root string
}

// IsDir reports whether dir looks like a catalog.
func IsDir(dir string) bool {
	st, err := os.Stat(File(dir))
	return err == nil && !st.IsDir()
}

// List returns catalogs under tree, sorted by name. Hidden dirs are skipped.
func List(tree string) ([]Source, error) {
	ents, err := os.ReadDir(tree)
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
		root := filepath.Join(tree, e.Name())
		if !IsDir(root) {
			continue
		}
		out = append(out, Source{Name: e.Name(), Root: root})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}

// Lookup loads <catalog>/<name> from the tree.
func Lookup(tree, recipeID string) (src Source, c *Catalog, agent Agent, err error) {
	catName, rec, err := ParseID(recipeID)
	if err != nil {
		return Source{}, nil, Agent{}, err
	}
	root := filepath.Join(tree, catName)
	if !IsDir(root) {
		return Source{}, nil, Agent{}, fmt.Errorf("unknown recipe %q (try: sbx-kit recipes)", recipeID)
	}
	c, err = Load(File(root))
	if err != nil {
		return Source{}, nil, Agent{}, err
	}
	agent, ok := c.Agents[rec]
	if !ok {
		return Source{}, nil, Agent{}, fmt.Errorf("unknown recipe %q (try: sbx-kit recipes)", recipeID)
	}
	return Source{Name: catName, Root: root}, c, agent, nil
}
