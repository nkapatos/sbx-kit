package catalog

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Catalog struct {
	Defaults Defaults          `yaml:"defaults"`
	Agents   map[string]Agent  `yaml:"agents"`
}

type Defaults struct {
	Resources string   `yaml:"resources"`
	Kits      []string `yaml:"kits"`
}

type Agent struct {
	SbxAgent          string   `yaml:"sbx_agent"`
	ImageName         string   `yaml:"image_name"`
	TemplateFallback  string   `yaml:"template_fallback"`
	Kits              []string `yaml:"kits"`
	Stub              bool     `yaml:"stub"`
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
