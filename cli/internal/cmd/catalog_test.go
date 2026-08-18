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

func TestSetupTreeAndCatalogLs(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("SBX_KIT_TREE", "")
	tree := t.TempDir()

	root := NewRoot()
	out := &bytes.Buffer{}
	root.SetOut(out)
	root.SetErr(out)
	root.SetArgs([]string{"setup", "--tree", tree})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "tree:") {
		t.Fatalf("setup out: %s", out.String())
	}

	p := filepath.Join(tree, "mine", "recipes", "agents.yaml")
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

func TestSetupRefusesCatalogDir(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("SBX_KIT_TREE", "")
	tree := t.TempDir()
	p := filepath.Join(tree, "recipes", "agents.yaml")
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte("agents: {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	root := NewRoot()
	root.SetArgs([]string{"setup", "--tree", tree})
	err := root.Execute()
	if err == nil || !strings.Contains(err.Error(), "looks like a catalog") {
		t.Fatalf("got %v", err)
	}
}
