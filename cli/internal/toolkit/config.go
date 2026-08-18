package toolkit

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/nkapatos/sbx-kit/cli/internal/xdg"
)

type userConfig struct {
	Tree string `yaml:"tree"`
}

func configFile() (string, error) {
	dir, err := xdg.ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.yaml"), nil
}

func configuredTree() (string, error) {
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
	if c.Tree == "" {
		return "", nil
	}
	return filepath.Clean(c.Tree), nil
}

// WriteTree records the recipes tree path in ~/.config/sbx-kit/config.yaml.
func WriteTree(dir string) (string, error) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}
	if !IsTree(abs) {
		return "", fmt.Errorf("%s is not a recipes tree (need %s)", abs, catalogRel)
	}
	path, err := configFile()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", err
	}
	c := userConfig{}
	if b, err := os.ReadFile(path); err == nil {
		_ = yaml.Unmarshal(b, &c)
	}
	c.Tree = abs
	out, err := yaml.Marshal(&c)
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(path, out, 0o644); err != nil {
		return "", err
	}
	return abs, nil
}
