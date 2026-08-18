package toolkit

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	TreeEnv      = "SBX_KIT_TREE"
	catalogRel   = "recipes/agents.yaml"
	errNeedSetup = "no recipes tree; run: sbx-kit setup"
)

// IsCatalog reports whether dir is a catalog (recipes/agents.yaml).
func IsCatalog(dir string) bool {
	st, err := os.Stat(filepath.Join(dir, catalogRel))
	return err == nil && !st.IsDir()
}

// HasCatalogs reports whether dir has at least one one-level catalog child.
func HasCatalogs(dir string) bool {
	ents, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, e := range ents {
		if !e.IsDir() || isHidden(e.Name()) {
			continue
		}
		if IsCatalog(filepath.Join(dir, e.Name())) {
			return true
		}
	}
	return false
}

func isHidden(name string) bool {
	return name == "" || name[0] == '.'
}

// Root locates the recipes tree (parent of catalogs).
//
// Order: SBX_KIT_TREE, then setup config, then walk cwd.
func Root() (string, error) {
	if t := os.Getenv(TreeEnv); t != "" {
		return checkTree(filepath.Clean(t), TreeEnv+"="+t)
	}

	if t, err := ConfiguredTree(); err == nil && t != "" {
		return checkTree(t, "setup tree "+t)
	}

	if wd, err := os.Getwd(); err == nil {
		if root := walkForRoot(wd); root != "" {
			return root, nil
		}
	}

	return "", fmt.Errorf("%s", errNeedSetup)
}

func checkTree(dir, label string) (string, error) {
	st, err := os.Stat(dir)
	if err != nil {
		return "", fmt.Errorf("%s is not a directory", label)
	}
	if !st.IsDir() {
		return "", fmt.Errorf("%s is not a directory", label)
	}
	if IsCatalog(dir) {
		return "", fmt.Errorf("%s looks like a catalog (has %s); point setup at the parent directory", label, catalogRel)
	}
	return dir, nil
}

func walkForRoot(start string) string {
	dir := start
	for {
		if HasCatalogs(dir) {
			return dir
		}
		if IsCatalog(dir) {
			parent := filepath.Dir(dir)
			if parent == dir {
				return ""
			}
			return parent
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}
