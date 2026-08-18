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

// ResolveBuild finds images/<name> (or an absolute/relative dir) and
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
	// Strip optional images/ or templates/ prefix.
	name := strings.TrimPrefix(nameOrPath, "images/")
	name = strings.TrimPrefix(name, "images"+string(filepath.Separator))
	name = strings.TrimPrefix(name, "templates/")
	name = strings.TrimPrefix(name, "templates"+string(filepath.Separator))
	cand := filepath.Join(root, "images", name)
	st, err := os.Stat(cand)
	if err != nil || !st.IsDir() {
		return "", fmt.Errorf("template not found: %s (tried %s)", nameOrPath, cand)
	}
	return cand, nil
}

const parentMarker = "PARENT"

// IsParentOnly is true when the template is a Docker FROM base and must not be
// imported into sbx (no sandbox recipe).
func IsParentOnly(templateDir string) bool {
	st, err := os.Stat(filepath.Join(templateDir, parentMarker))
	return err == nil && !st.IsDir()
}

// ParentTemplateName reads ARG PARENT_IMAGE=local/sbx-<name>:… from a
// Dockerfile. Empty means no in-tree parent to docker-build first.
func ParentTemplateName(dockerfile string) (string, error) {
	f, err := os.Open(dockerfile)
	if err != nil {
		return "", err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		rest, ok := strings.CutPrefix(line, "ARG PARENT_IMAGE=")
		if !ok {
			continue
		}
		val := strings.Trim(strings.TrimSpace(rest), `"'`)
		return localSbxTemplateName(val), nil
	}
	return "", sc.Err()
}

func localSbxTemplateName(imageTag string) string {
	const p = "local/sbx-"
	if !strings.HasPrefix(imageTag, p) {
		return ""
	}
	name, _, _ := strings.Cut(strings.TrimPrefix(imageTag, p), ":")
	return name
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
