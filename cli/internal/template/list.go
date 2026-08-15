package template

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// LocalImage is a custom image directory under templates/.
type LocalImage struct {
	Name     string
	ImageTag string
	Role     string // "parent" (FROM only) or "load" (import into sbx)
}

// ListLocal returns templates/ dirs that have a Dockerfile or bake.env.
func ListLocal(root string) ([]LocalImage, error) {
	dir := filepath.Join(root, "templates")
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var out []LocalImage
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		if name == "" || strings.HasPrefix(name, "_") || strings.HasPrefix(name, ".") {
			continue
		}
		base := filepath.Join(dir, name)
		if !fileExists(filepath.Join(base, "Dockerfile")) && !fileExists(filepath.Join(base, "bake.env")) {
			continue
		}
		role := "load"
		if IsParentOnly(base) {
			role = "parent"
		}
		out = append(out, LocalImage{
			Name:     name,
			ImageTag: "local/sbx-" + name + ":latest",
			Role:     role,
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}

func fileExists(path string) bool {
	st, err := os.Stat(path)
	return err == nil && !st.IsDir()
}
