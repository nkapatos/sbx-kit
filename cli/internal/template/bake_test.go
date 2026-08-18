package template_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nkapatos/sbx-kit/cli/internal/template"
)

func TestResolveBuildKitCoreDockerfile(t *testing.T) {
	root := fixtureTree(t)
	b, err := template.ResolveBuild(root, "kit-core", "")
	if err != nil {
		t.Fatal(err)
	}
	if b.ImageTag != "local/sbx-kit-core:latest" {
		t.Fatalf("tag: %s", b.ImageTag)
	}
	if filepath.Base(b.Context) != "kit-core" {
		t.Fatalf("context: %s", b.Context)
	}
	if filepath.Base(b.Dockerfile) != "Dockerfile" {
		t.Fatalf("dockerfile: %s", b.Dockerfile)
	}
	if len(b.BuildArgs) != 0 {
		t.Fatalf("unexpected build args: %#v", b.BuildArgs)
	}
}

func TestResolveBuildKitShellDockerfile(t *testing.T) {
	root := fixtureTree(t)
	b, err := template.ResolveBuild(root, "kit-shell", "")
	if err != nil {
		t.Fatal(err)
	}
	if b.ImageTag != "local/sbx-kit-shell:latest" {
		t.Fatalf("tag: %s", b.ImageTag)
	}
	if filepath.Base(b.Context) != "kit-shell" {
		t.Fatalf("context: %s", b.Context)
	}
}

func TestResolveBuildKitCursorDockerfile(t *testing.T) {
	root := fixtureTree(t)
	b, err := template.ResolveBuild(root, "kit-cursor", "")
	if err != nil {
		t.Fatal(err)
	}
	if b.ImageTag != "local/sbx-kit-cursor:latest" {
		t.Fatalf("tag: %s", b.ImageTag)
	}
	if filepath.Base(b.Context) != "kit-cursor" {
		t.Fatalf("context: %s", b.Context)
	}
}

func TestResolveBuildByPath(t *testing.T) {
	root := fixtureTree(t)
	rel := filepath.Join(root, "images", "kit-core")
	b, err := template.ResolveBuild(root, rel, "local/custom:dev")
	if err != nil {
		t.Fatal(err)
	}
	if b.ImageTag != "local/custom:dev" {
		t.Fatalf("tag: %s", b.ImageTag)
	}
}

func TestListLocal(t *testing.T) {
	root := fixtureTree(t)
	imgs, err := template.ListLocal(root)
	if err != nil {
		t.Fatal(err)
	}
	got := map[string]string{}
	for _, img := range imgs {
		got[img.Name] = img.ImageTag
	}
	if got["kit-core"] != "local/sbx-kit-core:latest" {
		t.Fatalf("kit-core: %#v", got)
	}
	if got["kit-shell"] != "local/sbx-kit-shell:latest" {
		t.Fatalf("kit-shell: %#v", got)
	}
	if got["kit-cursor"] != "local/sbx-kit-cursor:latest" {
		t.Fatalf("kit-cursor: %#v", got)
	}
	roles := map[string]string{}
	for _, img := range imgs {
		roles[img.Name] = img.Role
	}
	if roles["kit-core"] != "parent" {
		t.Fatalf("kit-core role: %#v", roles)
	}
	if roles["kit-shell"] != "load" || roles["kit-cursor"] != "load" {
		t.Fatalf("loadable roles: %#v", roles)
	}
}

func TestParentOnlyAndParentName(t *testing.T) {
	root := fixtureTree(t)
	core, err := template.ResolveBuild(root, "kit-core", "")
	if err != nil {
		t.Fatal(err)
	}
	if !template.IsParentOnly(core.TemplateDir) {
		t.Fatal("kit-core should be parent-only")
	}
	name, err := template.ParentTemplateName(core.Dockerfile)
	if err != nil || name != "" {
		t.Fatalf("kit-core parent name: %q %v", name, err)
	}

	shell, err := template.ResolveBuild(root, "kit-shell", "")
	if err != nil {
		t.Fatal(err)
	}
	if template.IsParentOnly(shell.TemplateDir) {
		t.Fatal("kit-shell should be loadable")
	}
	name, err = template.ParentTemplateName(shell.Dockerfile)
	if err != nil || name != "kit-core" {
		t.Fatalf("kit-shell parent: %q %v", name, err)
	}

	cur, err := template.ResolveBuild(root, "kit-cursor", "")
	if err != nil {
		t.Fatal(err)
	}
	name, err = template.ParentTemplateName(cur.Dockerfile)
	if err != nil || name != "kit-core" {
		t.Fatalf("kit-cursor parent: %q %v", name, err)
	}
}

func TestLoadRejectsParentImage(t *testing.T) {
	root := fixtureTree(t)
	err := template.Load(template.LoadOpts{
		Root:       root,
		Engine:     "docker",
		NameOrPath: "kit-core",
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "not imported") {
		t.Fatalf("got: %v", err)
	}
}

func fixtureTree(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	mustWrite := func(rel, body string) {
		t.Helper()
		p := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	mustWrite("recipes/agents.yaml", "agents: {}\n")
	mustWrite("images/kit-core/Dockerfile", "FROM scratch\n")
	mustWrite("images/kit-core/PARENT", "")
	mustWrite("images/kit-shell/Dockerfile", "ARG PARENT_IMAGE=local/sbx-kit-core:latest\nFROM scratch\n")
	mustWrite("images/kit-cursor/Dockerfile", "ARG PARENT_IMAGE=local/sbx-kit-core:latest\nFROM scratch\n")
	return root
}
