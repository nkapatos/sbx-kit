package recipecreate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCreateBundle(t *testing.T) {
	root := t.TempDir()
	err := Create(CreateOpts{
		CatalogRoot: root,
		DirName:     "mine",
		RecipeName:  "cursor",
		SbxAgent:    "cursor",
		DefaultKits: []string{"agent-workspace"},
		WriteAgents: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	manifest := filepath.Join(root, "mine", "recipes", "agents.yaml")
	b, err := os.ReadFile(manifest)
	if err != nil {
		t.Fatal(err)
	}
	text := string(b)
	if !strings.Contains(text, "sbx_agent: cursor") {
		t.Fatalf("manifest: %s", text)
	}
	agents := filepath.Join(root, "mine", "AGENTS.md")
	if _, err := os.Stat(agents); err != nil {
		t.Fatal(err)
	}
}

func TestRenderSkill(t *testing.T) {
	out, err := RenderSkill(SkillOpts{CatalogRoot: "/tmp/cat", DirName: "mine"})
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"name: sbx-kit",
		"docs.docker.com/ai/sandboxes/",
		"SPEC-v2.md",
		"sbx-kit recipes verify",
		"/tmp/cat",
		"mine",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("missing %q in skill output", want)
		}
	}
}

func TestCreateRejectsExisting(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "mine"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := Create(CreateOpts{CatalogRoot: root, DirName: "mine"}); err == nil {
		t.Fatal("expected error for existing dir")
	}
}
