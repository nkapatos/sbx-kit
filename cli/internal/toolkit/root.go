package toolkit

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	CatalogEnv = "SBX_KIT_CATALOG"

	recipeManifestRel = "recipes/agents.yaml"
	errNeedSetup      = "no catalog configured; run: sbx-kit setup"
)

// IsRecipeDir reports whether dir holds recipes (has recipes/agents.yaml).
func IsRecipeDir(dir string) bool {
	st, err := os.Stat(filepath.Join(dir, recipeManifestRel))
	return err == nil && !st.IsDir()
}

// HasRecipeDirs reports whether dir has at least one recipe directory child.
func HasRecipeDirs(dir string) bool {
	ents, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, e := range ents {
		if !e.IsDir() || isHidden(e.Name()) {
			continue
		}
		if IsRecipeDir(filepath.Join(dir, e.Name())) {
			return true
		}
	}
	return false
}

func isHidden(name string) bool {
	return name == "" || name[0] == '.'
}

// Root locates the catalog path configured by setup.
//
// Order: SBX_KIT_CATALOG, setup config, then walk cwd.
func Root() (string, error) {
	if t := os.Getenv(CatalogEnv); t != "" {
		return checkCatalogRoot(filepath.Clean(t), CatalogEnv+"="+t)
	}

	if t, err := ConfiguredCatalog(); err == nil && t != "" {
		return checkCatalogRoot(t, "setup catalog "+t)
	}

	if wd, err := os.Getwd(); err == nil {
		if root := walkForRoot(wd); root != "" {
			return root, nil
		}
	}

	return "", fmt.Errorf("%s", errNeedSetup)
}

func checkCatalogRoot(dir, label string) (string, error) {
	st, err := os.Stat(dir)
	if err != nil {
		return "", fmt.Errorf("%s is not a directory", label)
	}
	if !st.IsDir() {
		return "", fmt.Errorf("%s is not a directory", label)
	}
	if IsRecipeDir(dir) {
		return "", fmt.Errorf("%s is a recipe directory; run setup on the catalog path above it", label)
	}
	return dir, nil
}

func walkForRoot(start string) string {
	dir := start
	for {
		if HasRecipeDirs(dir) {
			return dir
		}
		if IsRecipeDir(dir) {
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
