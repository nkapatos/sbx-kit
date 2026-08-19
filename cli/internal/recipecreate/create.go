package recipecreate

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/nkapatos/sbx-kit/cli/internal/sbxcompat"
)

//go:embed templates/*
var templateFS embed.FS

// CreateOpts scaffolds a new catalog directory bundle.
type CreateOpts struct {
	CatalogRoot string
	DirName     string
	RecipeName  string
	SbxAgent    string
	DefaultKits []string
	Resources   string
	WriteAgents bool
	Force       bool
}

// SkillOpts renders agent skill markdown.
type SkillOpts struct {
	CatalogRoot string // optional; adds catalog path to skill body
	DirName     string // optional; adds directory-specific section
}

type createData struct {
	DirName     string
	RecipeName  string
	SbxAgent    string
	DefaultKits []string
	Resources   string
	RecipeID    string
}

type skillData struct {
	MinSbx         string
	MinKitVerify   string
	CatalogRoot    string
	DirName        string
	HasCatalogHint bool
	HasDirHint     bool
}

// Create writes recipes/, optional kits/ and images/, and AGENTS.md.
func Create(o CreateOpts) error {
	if err := validateDirName(o.DirName); err != nil {
		return err
	}
	if o.RecipeName == "" {
		o.RecipeName = "shell"
	}
	if o.SbxAgent == "" {
		o.SbxAgent = "shell"
	}
	if o.Resources == "" {
		o.Resources = "remote-llm"
	}
	if len(o.DefaultKits) == 0 {
		o.DefaultKits = []string{"agent-workspace"}
	}

	root := filepath.Join(o.CatalogRoot, o.DirName)
	if st, err := os.Stat(root); err == nil {
		if !o.Force {
			return fmt.Errorf("directory %q already exists (pass --force)", o.DirName)
		}
		if !st.IsDir() {
			return fmt.Errorf("%q is not a directory", root)
		}
	} else if !os.IsNotExist(err) {
		return err
	}

	data := createData{
		DirName:     o.DirName,
		RecipeName:  o.RecipeName,
		SbxAgent:    o.SbxAgent,
		DefaultKits: o.DefaultKits,
		Resources:   o.Resources,
		RecipeID:    o.DirName + "/" + o.RecipeName,
	}

	if err := os.MkdirAll(filepath.Join(root, "recipes"), 0o755); err != nil {
		return err
	}
	if err := writeTemplate(filepath.Join(root, "recipes", "agents.yaml"), "agents.yaml.tmpl", data, o.Force); err != nil {
		return err
	}

	for _, sub := range []string{"kits", "images"} {
		if err := os.MkdirAll(filepath.Join(root, sub), 0o755); err != nil {
			return err
		}
		keep := filepath.Join(root, sub, ".gitkeep")
		if err := writeFileIfMissing(keep, o.Force); err != nil {
			return err
		}
	}

	if o.WriteAgents {
		if err := writeTemplate(filepath.Join(root, "AGENTS.md"), "agents.md.tmpl", data, o.Force); err != nil {
			return err
		}
	}

	fmt.Printf("created catalog directory %s\n", root)
	fmt.Printf("  recipe id:  %s\n", data.RecipeID)
	fmt.Printf("  edit:       %s/recipes/agents.yaml\n", root)
	fmt.Printf("  verify:     sbx-kit recipes verify %s\n", data.RecipeID)
	fmt.Printf("  run:        sbx-kit box run %s --yes\n", data.RecipeID)
	return nil
}

// RenderSkill returns Cursor-style SKILL.md content for sbx-kit.
func RenderSkill(o SkillOpts) (string, error) {
	data := skillData{
		MinSbx:       sbxcompat.MinVersion,
		MinKitVerify: sbxcompat.MinKitVerify,
	}
	if o.CatalogRoot != "" {
		data.CatalogRoot = o.CatalogRoot
		data.HasCatalogHint = true
	}
	if o.DirName != "" {
		data.DirName = o.DirName
		data.HasDirHint = true
	}
	return execTemplate("skill.md.tmpl", data)
}

func validateDirName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("directory name is required")
	}
	if name == "." || name == ".." || strings.ContainsAny(name, `/\`) {
		return fmt.Errorf("invalid directory name %q", name)
	}
	if strings.HasPrefix(name, ".") {
		return fmt.Errorf("invalid directory name %q", name)
	}
	return nil
}

func writeTemplate(path, tmplName string, data any, force bool) error {
	body, err := execTemplate(tmplName, data)
	if err != nil {
		return err
	}
	if _, err := os.Stat(path); err == nil && !force {
		return nil
	} else if err != nil && !os.IsNotExist(err) {
		return err
	}
	return os.WriteFile(path, []byte(body), 0o644)
}

func writeFileIfMissing(path string, force bool) error {
	if _, err := os.Stat(path); err == nil {
		if force {
			return nil
		}
		return nil
	}
	return os.WriteFile(path, []byte(""), 0o644)
}

func execTemplate(name string, data any) (string, error) {
	b, err := templateFS.ReadFile(filepath.Join("templates", name))
	if err != nil {
		return "", fmt.Errorf("read template %s: %w", name, err)
	}
	t, err := template.New(name).Parse(string(b))
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
