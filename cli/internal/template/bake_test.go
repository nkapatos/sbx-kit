package template_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nkapatos/sbx-kit/cli/internal/template"
)

func TestResolveBuildBake(t *testing.T) {
	root := findRepoRoot(t)
	b, err := template.ResolveBuild(root, "cursor-mise-docker", "")
	if err != nil {
		t.Fatal(err)
	}
	if b.ImageTag != "local/sbx-cursor-mise-docker:latest" {
		t.Fatalf("tag: %s", b.ImageTag)
	}
	if filepath.Base(b.Context) != "_bake" {
		t.Fatalf("context: %s", b.Context)
	}
	if len(b.BuildArgs) != 1 || b.BuildArgs[0] == "" {
		t.Fatalf("build args: %#v", b.BuildArgs)
	}
}

func TestResolveBuildByPath(t *testing.T) {
	root := findRepoRoot(t)
	rel := filepath.Join(root, "templates", "cursor-mise-docker")
	b, err := template.ResolveBuild(root, rel, "local/custom:dev")
	if err != nil {
		t.Fatal(err)
	}
	if b.ImageTag != "local/custom:dev" {
		t.Fatalf("tag: %s", b.ImageTag)
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
