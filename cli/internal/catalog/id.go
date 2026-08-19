package catalog

import (
	"fmt"
	"strings"
)

// ParseID splits <dir>/<recipe>.
func ParseID(id string) (dirName, recipeName string, err error) {
	id = strings.TrimSpace(id)
	dir, rest, ok := strings.Cut(id, "/")
	if !ok || dir == "" || rest == "" {
		return "", "", fmt.Errorf("recipe id is <dir>/<name> (got %q; try: sbx-kit recipes)", id)
	}
	if dir == "." || dir == ".." || strings.ContainsAny(dir, `/\`) {
		return "", "", fmt.Errorf("invalid directory %q", dir)
	}
	if strings.HasPrefix(dir, ".") {
		return "", "", fmt.Errorf("invalid directory %q", dir)
	}
	return dir, rest, nil
}

// JoinID builds <dir>/<recipe>.
func JoinID(dirName, recipeName string) string {
	return dirName + "/" + recipeName
}
