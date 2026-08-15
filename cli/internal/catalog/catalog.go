package catalog

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Catalog struct {
	Defaults Defaults         `yaml:"defaults"`
	Agents   map[string]Agent `yaml:"agents"`
}

type Defaults struct {
	Resources string   `yaml:"resources"`
	Kits      []string `yaml:"kits"`
}

type Agent struct {
	SbxAgent         string   `yaml:"sbx_agent"`
	ImageName        string   `yaml:"image_name"`
	TemplateFallback string   `yaml:"template_fallback"`
	Kits             []string `yaml:"kits"`
	Stub             bool     `yaml:"stub"`
}

func Load(path string) (*Catalog, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var c Catalog
	if err := yaml.Unmarshal(b, &c); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	if c.Agents == nil {
		c.Agents = map[string]Agent{}
	}
	return &c, nil
}

// ResolveKits is recipe kits first, then catalog defaults not already listed.
// An empty recipe list means "defaults only" (typically agent-workspace).
func ResolveKits(recipeKits, defaults []string) []string {
	if len(recipeKits) == 0 {
		return append([]string(nil), defaults...)
	}
	seen := make(map[string]struct{}, len(recipeKits)+len(defaults))
	out := make([]string, 0, len(recipeKits)+len(defaults))
	for _, xs := range [][]string{recipeKits, defaults} {
		for _, k := range xs {
			if k == "" {
				continue
			}
			if _, ok := seen[k]; ok {
				continue
			}
			seen[k] = struct{}{}
			out = append(out, k)
		}
	}
	return out
}

// KitPaths maps catalog kit names to directories under the toolkit root.
func KitPaths(root string, names []string) []string {
	out := make([]string, 0, len(names))
	for _, n := range names {
		out = append(out, filepath.Join(root, "kits", n))
	}
	return out
}
