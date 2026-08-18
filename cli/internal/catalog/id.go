package catalog

import (
	"fmt"
	"strings"
)

// ParseID splits <catalog>/<name>. Catalog is the one-level tree child.
func ParseID(id string) (catalogName, name string, err error) {
	id = strings.TrimSpace(id)
	cat, rest, ok := strings.Cut(id, "/")
	if !ok || cat == "" || rest == "" {
		return "", "", fmt.Errorf("recipe id is <catalog>/<name> (got %q; try: sbx-kit recipes)", id)
	}
	if cat == "." || cat == ".." || strings.ContainsAny(cat, `/\`) {
		return "", "", fmt.Errorf("invalid catalog %q", cat)
	}
	if strings.HasPrefix(cat, ".") {
		return "", "", fmt.Errorf("invalid catalog %q", cat)
	}
	return cat, rest, nil
}

// JoinID builds <catalog>/<name>.
func JoinID(catalogName, name string) string {
	return catalogName + "/" + name
}
