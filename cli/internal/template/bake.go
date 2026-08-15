package template

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Build describes how to docker/container-build a thin or legacy template.
type Build struct {
	Name        string
	TemplateDir string
	Context     string
	Dockerfile  string
	BuildArgs   []string // e.g. BASE_IMAGE=...
	ImageTag    string
}

// ResolveBuild finds templates/<name> (or an absolute/relative dir) and
// resolves bake.env → _bake + BASE_IMAGE, or a local Dockerfile.
func ResolveBuild(root, nameOrPath, imageTag string) (*Build, error) {
	dir, err := resolveTemplateDir(root, nameOrPath)
	if err != nil {
		return nil, err
	}
	name := filepath.Base(dir)
	if imageTag == "" {
		imageTag = "local/sbx-" + name + ":latest"
	}

	bakeEnv := filepath.Join(dir, "bake.env")
	if st, err := os.Stat(bakeEnv); err == nil && !st.IsDir() {
		baseImage, err := readBakeBaseImage(bakeEnv)
		if err != nil {
			return nil, err
		}
		bakeDir := filepath.Join(filepath.Dir(dir), "_bake")
		df := filepath.Join(bakeDir, "Dockerfile")
		if _, err := os.Stat(df); err != nil {
			return nil, fmt.Errorf("shared bake missing: %s", df)
		}
		fmt.Printf("==> shared bake: BASE_IMAGE=%s\n", baseImage)
		return &Build{
			Name:        name,
			TemplateDir: dir,
			Context:     bakeDir,
			Dockerfile:  df,
			BuildArgs:   []string{"BASE_IMAGE=" + baseImage},
			ImageTag:    imageTag,
		}, nil
	}

	df := filepath.Join(dir, "Dockerfile")
	if _, err := os.Stat(df); err != nil {
		return nil, fmt.Errorf("need %s (shared bake) or %s", bakeEnv, df)
	}
	return &Build{
		Name:        name,
		TemplateDir: dir,
		Context:     dir,
		Dockerfile:  df,
		ImageTag:    imageTag,
	}, nil
}

func resolveTemplateDir(root, nameOrPath string) (string, error) {
	if nameOrPath == "" {
		return "", fmt.Errorf("template name or path required")
	}
	// Absolute or existing relative path to a directory.
	if st, err := os.Stat(nameOrPath); err == nil && st.IsDir() {
		abs, err := filepath.Abs(nameOrPath)
		if err != nil {
			return "", err
		}
		return abs, nil
	}
	// Strip optional templates/ prefix.
	name := strings.TrimPrefix(nameOrPath, "templates/")
	name = strings.TrimPrefix(name, "templates"+string(filepath.Separator))
	cand := filepath.Join(root, "templates", name)
	st, err := os.Stat(cand)
	if err != nil || !st.IsDir() {
		return "", fmt.Errorf("template not found: %s (tried %s)", nameOrPath, cand)
	}
	return cand, nil
}

func readBakeBaseImage(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var base string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, val, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		if strings.TrimSpace(key) == "BASE_IMAGE" {
			base = strings.TrimSpace(val)
			base = strings.Trim(base, `"'`)
		}
	}
	if err := sc.Err(); err != nil {
		return "", err
	}
	if base == "" {
		return "", fmt.Errorf("%s must set BASE_IMAGE=", path)
	}
	return base, nil
}
