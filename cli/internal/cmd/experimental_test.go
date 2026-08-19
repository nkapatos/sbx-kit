package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExperimentalVerifyRecipeStub(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	catalogRoot := t.TempDir()
	t.Setenv("SBX_KIT_CATALOG", catalogRoot)

	root := NewRoot()
	out := &bytes.Buffer{}
	root.SetOut(out)
	root.SetErr(out)
	root.SetArgs([]string{"experimental", "verify", "recipe"})
	err := root.Execute()
	if err == nil || !strings.Contains(err.Error(), "not implemented") {
		t.Fatalf("expected stub error, got %v out=%s", err, out.String())
	}
}

func TestExperimentalSpecPrintsStatus(t *testing.T) {
	root := NewRoot()
	out := &bytes.Buffer{}
	root.SetOut(out)
	root.SetErr(out)
	root.SetArgs([]string{"experimental", "spec"})
	err := root.Execute()
	if err == nil || !strings.Contains(err.Error(), "recipe spec") {
		t.Fatalf("expected spec stub error, got %v", err)
	}
	if !strings.Contains(out.String(), "agents.yaml") {
		t.Fatalf("expected spec text, got %s", out.String())
	}
}

func TestBoxRunHelp(t *testing.T) {
	root := NewRoot()
	out := &bytes.Buffer{}
	root.SetOut(out)
	root.SetErr(out)
	root.SetArgs([]string{"box", "run", "--help"})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "run") {
		t.Fatalf("box run help: %s", out.String())
	}
}

func TestRecipesImageNested(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	catalogRoot := t.TempDir()
	t.Setenv("SBX_KIT_CATALOG", catalogRoot)

	p := filepath.Join(catalogRoot, "mine", "recipes", "agents.yaml")
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte("agents: {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	root := NewRoot()
	out := &bytes.Buffer{}
	root.SetOut(out)
	root.SetErr(out)
	root.SetArgs([]string{"recipes", "image", "ls"})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
}
