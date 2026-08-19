package catalog

import (
	"fmt"
	"strings"
)

// ValidDirName reports whether name is a safe catalog directory component.
func ValidDirName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("directory name is required")
	}
	if name == "." || name == ".." || strings.ContainsAny(name, `/\`) {
		return fmt.Errorf("invalid directory %q", name)
	}
	if strings.HasPrefix(name, ".") {
		return fmt.Errorf("invalid directory %q", name)
	}
	return nil
}

// ParseID splits <dir>/<recipe>.
func ParseID(id string) (dirName, recipeName string, err error) {
	id = strings.TrimSpace(id)
	dir, rest, ok := strings.Cut(id, "/")
	if !ok || dir == "" || rest == "" {
		return "", "", fmt.Errorf("recipe id is <dir>/<name> (got %q; try: sbx-kit recipes)", id)
	}
	if err := ValidDirName(dir); err != nil {
		return "", "", err
	}
	return dir, rest, nil
}

// JoinID builds <dir>/<recipe>.
func JoinID(dirName, recipeName string) string {
	return dirName + "/" + recipeName
}
