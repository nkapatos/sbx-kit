package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

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

func TestExperimentalPromptsParked(t *testing.T) {
	root := NewRoot()
	out := &bytes.Buffer{}
	root.SetOut(out)
	root.SetErr(out)
	root.SetArgs([]string{"experimental", "prompts", "state"})
	err := root.Execute()
	if err == nil || !strings.Contains(err.Error(), "prompts") {
		t.Fatalf("expected prompts stub error, got %v", err)
	}
	if !strings.Contains(out.String(), "sbx-kit-state") {
		t.Fatalf("expected state prompt, got %s", out.String())
	}
}

func TestRecipesVerifySkipKits(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	catalogRoot := t.TempDir()
	t.Setenv("SBX_KIT_CATALOG", catalogRoot)

	dir := filepath.Join(catalogRoot, "mine")
	agents := filepath.Join(dir, "recipes", "agents.yaml")
	if err := os.MkdirAll(filepath.Dir(agents), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(agents, []byte("agents:\n  cursor:\n    sbx_agent: cursor\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	root := NewRoot()
	out := &bytes.Buffer{}
	root.SetOut(out)
	root.SetErr(out)
	root.SetArgs([]string{"recipes", "verify", "--skip-kits", "mine/cursor"})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "recipe mine/cursor: ok") {
		t.Fatalf("out: %s", out.String())
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
