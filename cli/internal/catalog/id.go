package catalog

import (
	"fmt"
	"strings"
)

// ParseID splits <source>/<name>.
func ParseID(id string) (sourceName, name string, err error) {
	id = strings.TrimSpace(id)
	src, rest, ok := strings.Cut(id, "/")
	if !ok || src == "" || rest == "" {
		return "", "", fmt.Errorf("recipe id is <source>/<name> (got %q; try: sbx-kit recipes)", id)
	}
	if src == "." || src == ".." || strings.ContainsAny(src, `/\`) {
		return "", "", fmt.Errorf("invalid source %q", src)
	}
	if strings.HasPrefix(src, ".") {
		return "", "", fmt.Errorf("invalid source %q", src)
	}
	return src, rest, nil
}

// JoinID builds <source>/<name>.
func JoinID(sourceName, name string) string {
	return sourceName + "/" + name
}
