package toolkit

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	CatalogEnv = "SBX_KIT_CATALOG"

	sourceManifestRel = "recipes/agents.yaml"
	errNeedSetup      = "no catalog configured; run: sbx-kit setup"
)

// IsSource reports whether dir is a source (has recipes/agents.yaml).
func IsSource(dir string) bool {
	st, err := os.Stat(filepath.Join(dir, sourceManifestRel))
	return err == nil && !st.IsDir()
}

// HasSources reports whether dir has at least one source child.
func HasSources(dir string) bool {
	ents, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, e := range ents {
		if !e.IsDir() || isHidden(e.Name()) {
			continue
		}
		if IsSource(filepath.Join(dir, e.Name())) {
			return true
		}
	}
	return false
}

func isHidden(name string) bool {
	return name == "" || name[0] == '.'
}

// Root locates the catalog directory (parent of sources).
//
// Order: SBX_KIT_CATALOG, setup config, then walk cwd.
func Root() (string, error) {
	if t := os.Getenv(CatalogEnv); t != "" {
		return checkCatalog(filepath.Clean(t), CatalogEnv+"="+t)
	}

	if t, err := ConfiguredCatalog(); err == nil && t != "" {
		return checkCatalog(t, "setup catalog "+t)
	}

	if wd, err := os.Getwd(); err == nil {
		if root := walkForRoot(wd); root != "" {
			return root, nil
		}
	}

	return "", fmt.Errorf("%s", errNeedSetup)
}

func checkCatalog(dir, label string) (string, error) {
	st, err := os.Stat(dir)
	if err != nil {
		return "", fmt.Errorf("%s is not a directory", label)
	}
	if !st.IsDir() {
		return "", fmt.Errorf("%s is not a directory", label)
	}
	if IsSource(dir) {
		return "", fmt.Errorf("%s looks like a source (has %s); point setup at the catalog directory", label, sourceManifestRel)
	}
	return dir, nil
}

func walkForRoot(start string) string {
	dir := start
	for {
		if HasSources(dir) {
			return dir
		}
		if IsSource(dir) {
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
