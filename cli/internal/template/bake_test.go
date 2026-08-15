package template_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nkapatos/sbx-kit/cli/internal/template"
)

func TestResolveBuildKitShellDockerfile(t *testing.T) {
	root := findRepoRoot(t)
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
	if filepath.Base(b.Dockerfile) != "Dockerfile" {
		t.Fatalf("dockerfile: %s", b.Dockerfile)
	}
	if len(b.BuildArgs) != 0 {
		t.Fatalf("unexpected build args: %#v", b.BuildArgs)
	}
}

func TestResolveBuildKitCursorDockerfile(t *testing.T) {
	root := findRepoRoot(t)
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
	root := findRepoRoot(t)
	rel := filepath.Join(root, "templates", "kit-shell")
	b, err := template.ResolveBuild(root, rel, "local/custom:dev")
	if err != nil {
		t.Fatal(err)
	}
	if b.ImageTag != "local/custom:dev" {
		t.Fatalf("tag: %s", b.ImageTag)
	}
}

func TestListLocal(t *testing.T) {
	root := findRepoRoot(t)
	imgs, err := template.ListLocal(root)
	if err != nil {
		t.Fatal(err)
	}
	got := map[string]string{}
	for _, img := range imgs {
		got[img.Name] = img.ImageTag
	}
	if got["kit-shell"] != "local/sbx-kit-shell:latest" {
		t.Fatalf("kit-shell: %#v", got)
	}
	if got["kit-cursor"] != "local/sbx-kit-cursor:latest" {
		t.Fatalf("kit-cursor: %#v", got)
	}
}

func findRepoRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	// cli/internal/template → repo root
	root := filepath.Clean(filepath.Join(wd, "../../.."))
	if _, err := os.Stat(filepath.Join(root, "config", "agents.yaml")); err != nil {
		t.Fatalf("repo root not found from %s: %v", wd, err)
	}
	return root
}
