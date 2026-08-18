package toolkit

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	TreeEnv      = "SBX_KIT_TREE"
	catalogRel   = "recipes/agents.yaml"
	errNeedSetup = "no recipes tree; run: sbx-kit setup --tree <dir>"
)

// IsTree reports whether dir looks like a recipes/kits/images tree.
func IsTree(dir string) bool {
	st, err := os.Stat(filepath.Join(dir, catalogRel))
	return err == nil && !st.IsDir()
}

// Root locates the recipes/kits/images tree.
//
// Order: SBX_KIT_TREE, then setup config, then walk cwd for recipes/agents.yaml.
func Root() (string, error) {
	if t := os.Getenv(TreeEnv); t != "" {
		if IsTree(t) {
			return filepath.Clean(t), nil
		}
		return "", fmt.Errorf("%s=%s is not a recipes tree (need %s)", TreeEnv, t, catalogRel)
	}

	if t, err := configuredTree(); err == nil && t != "" {
		if IsTree(t) {
			return t, nil
		}
		return "", fmt.Errorf("setup tree %s is not a recipes tree (need %s)", t, catalogRel)
	}

	if wd, err := os.Getwd(); err == nil {
		if root := walkForRoot(wd); root != "" {
			return root, nil
		}
	}

	return "", fmt.Errorf("%s", errNeedSetup)
}

func walkForRoot(start string) string {
	dir := start
	for {
		if IsTree(dir) {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}
