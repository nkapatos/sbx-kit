package toolkit

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/nkapatos/sbx-kit/cli/internal/xdg"
)

type userConfig struct {
	Catalog string `yaml:"catalog"`
}

func configFile() (string, error) {
	dir, err := xdg.ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.yaml"), nil
}

// ConfiguredCatalog returns the catalog path from setup config, or empty.
func ConfiguredCatalog() (string, error) {
	path, err := configFile()
	if err != nil {
		return "", err
	}
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	var c userConfig
	if err := yaml.Unmarshal(b, &c); err != nil {
		return "", fmt.Errorf("parse %s: %w", path, err)
	}
	if c.Catalog == "" {
		return "", nil
	}
	return filepath.Clean(c.Catalog), nil
}

// DefaultCatalog is ~/sbx-kit-catalog.
func DefaultCatalog() string {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return "sbx-kit-catalog"
	}
	return filepath.Join(home, "sbx-kit-catalog")
}

// WriteCatalog records the catalog path in ~/.config/sbx-kit/config.yaml and
// creates the directory if needed.
func WriteCatalog(dir string) (string, error) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}
	if IsRecipeDir(abs) {
		return "", fmt.Errorf("%s is a recipe directory; pass the catalog path that contains it", abs)
	}
	if err := os.MkdirAll(abs, 0o755); err != nil {
		return "", err
	}
	path, err := configFile()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", err
	}
	out, err := yaml.Marshal(userConfig{Catalog: abs})
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(path, out, 0o644); err != nil {
		return "", err
	}
	return abs, nil
}
