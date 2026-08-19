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
	if !strings.Contains(text, "kits: []") {
		t.Fatalf("expected empty kits list: %s", text)
	}
	if strings.Contains(text, "agent-workspace") {
		t.Fatalf("overlay must not be injected as a kit: %s", text)
	}
	stub := filepath.Join(root, "mine", "kits", "agent-workspace", "README.md")
	if _, err := os.Stat(stub); err == nil {
		t.Fatal("did not expect agent-workspace pull stub")
	}
	res := filepath.Join(root, "mine", "recipes", "resources-remote-llm.env")
	if _, err := os.Stat(res); err != nil {
		t.Fatal(err)
	}
	agents := filepath.Join(root, "mine", "AGENTS.md")
	ab, err := os.ReadFile(agents)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(ab), "agent-workspace") {
		t.Fatalf("AGENTS.md still requires a core kit: %s", ab)
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
		"overlay",
		"/etc/sbx-kit/context.md",
		"Two lands",
		"agentContext",
		"/tmp/cat",
		"mine",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("missing %q in skill output", want)
		}
	}
	if strings.Contains(out, "kit tree") || strings.Contains(out, "pull stub") {
		t.Fatalf("skill still tells agents to copy a core kit:\n%s", out)
	}
}

func TestCreateOptionalKits(t *testing.T) {
	root := t.TempDir()
	err := Create(CreateOpts{
		CatalogRoot: root,
		DirName:     "mine",
		DefaultKits: []string{"mise-workspace"},
	})
	if err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(filepath.Join(root, "mine", "recipes", "agents.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	text := string(b)
	if !strings.Contains(text, "mise-workspace") {
		t.Fatalf("expected optional kit: %s", text)
	}
	if strings.Contains(text, "agent-workspace") {
		t.Fatalf("did not expect core kit injection: %s", text)
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
