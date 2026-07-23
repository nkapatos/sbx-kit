package toolkit

import (
	"fmt"
	"os"
	"path/filepath"
)

// Root locates the toolkit data directory (config/, kits/, …).
//
// Order:
//  1. SBX_TREE
//  2. Homebrew share: <exe>/../share/sbx-kit (Cellar layout)
//  3. Walk up from executable dir and cwd looking for config/agents.yaml (dev checkout)
func Root() (string, error) {
	if t := os.Getenv("SBX_TREE"); t != "" {
		if isToolkitRoot(t) {
			return filepath.Clean(t), nil
		}
		return "", fmt.Errorf("SBX_TREE=%s does not contain config/agents.yaml", t)
	}

	if exe, err := os.Executable(); err == nil {
		exe, err := filepath.EvalSymlinks(exe)
		if err != nil {
			exe, _ = os.Executable()
		}
		binDir := filepath.Dir(exe)
		// Homebrew: .../Cellar/sbx-kit/<ver>/bin/sbx-kit → ../share/sbx-kit
		brewShare := filepath.Clean(filepath.Join(binDir, "..", "share", "sbx-kit"))
		if isToolkitRoot(brewShare) {
			return brewShare, nil
		}
		// Some taps install share next to opt prefix: .../opt/sbx-kit/bin + share
		optShare := filepath.Clean(filepath.Join(binDir, "..", "..", "share", "sbx-kit"))
		if isToolkitRoot(optShare) {
			return optShare, nil
		}
		if root := walkForRoot(binDir); root != "" {
			return root, nil
		}
	}

	if wd, err := os.Getwd(); err == nil {
		if root := walkForRoot(wd); root != "" {
			return root, nil
		}
	}

	return "", fmt.Errorf("cannot find toolkit root (brew share/sbx-kit or set SBX_TREE)")
}

func isToolkitRoot(dir string) bool {
	st, err := os.Stat(filepath.Join(dir, "config", "agents.yaml"))
	return err == nil && !st.IsDir()
}

func walkForRoot(start string) string {
	dir := start
	for {
		if isToolkitRoot(dir) {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}
