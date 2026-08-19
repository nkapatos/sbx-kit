package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRecipesCreateAndSkill(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	catalogRoot := t.TempDir()
	t.Setenv("SBX_KIT_CATALOG", catalogRoot)

	root := NewRoot()
	out := &bytes.Buffer{}
	root.SetOut(out)
	root.SetErr(out)
	root.SetArgs([]string{"recipes", "create", "mine", "--recipe", "shell"})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}

	manifest := filepath.Join(catalogRoot, "mine", "recipes", "agents.yaml")
	if _, err := os.Stat(manifest); err != nil {
		t.Fatal(err)
	}

	out = &bytes.Buffer{}
	root = NewRoot()
	root.SetOut(out)
	root.SetErr(out)
	root.SetArgs([]string{"recipes", "skill", "--dir", "mine"})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "name: sbx-kit") {
		t.Fatalf("skill out: %s", out.String())
	}
	if !strings.Contains(out.String(), "/etc/sbx-kit/context.md") {
		t.Fatalf("skill should point box agents at overlay context.md: %s", out.String())
	}
	if strings.Contains(out.String(), "agent-workspace") {
		t.Fatalf("skill must not require agent-workspace: %s", out.String())
	}
	stub := filepath.Join(catalogRoot, "mine", "kits", "agent-workspace", "README.md")
	if _, err := os.Stat(stub); err == nil {
		t.Fatal("did not expect agent-workspace pull stub")
	}
}
