package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTemplateOverrideEnv(t *testing.T) {
	got := templateOverrideEnv("mine/cursor")
	if got != "SBX_MINE_CURSOR_TEMPLATE" {
		t.Fatalf("got %s", got)
	}
}

func TestSetupCatalogAndCatalogLs(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("SBX_KIT_CATALOG", "")
	catalogRoot := t.TempDir()

	root := NewRoot()
	out := &bytes.Buffer{}
	root.SetOut(out)
	root.SetErr(out)
	root.SetArgs([]string{"setup", "--catalog", catalogRoot})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "catalog:") {
		t.Fatalf("setup out: %s", out.String())
	}

	p := filepath.Join(catalogRoot, "mine", "recipes", "agents.yaml")
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte("agents:\n  cursor:\n    sbx_agent: cursor\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	root = NewRoot()
	out = &bytes.Buffer{}
	root.SetOut(out)
	root.SetErr(out)
	root.SetArgs([]string{"catalog", "ls"})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "mine") {
		t.Fatalf("catalog ls: %s", out.String())
	}

	root = NewRoot()
	out = &bytes.Buffer{}
	root.SetOut(out)
	root.SetErr(out)
	root.SetArgs([]string{"recipes"})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "mine/cursor") {
		t.Fatalf("recipes: %s", out.String())
	}
}

func TestSetupRefusesRecipeDir(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("SBX_KIT_CATALOG", "")
	catalogRoot := t.TempDir()
	p := filepath.Join(catalogRoot, "recipes", "agents.yaml")
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte("agents: {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	root := NewRoot()
	root.SetArgs([]string{"setup", "--catalog", catalogRoot})
	err := root.Execute()
	if err == nil || !strings.Contains(err.Error(), "recipe directory") {
		t.Fatalf("got %v", err)
	}
}
