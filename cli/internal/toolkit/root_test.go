package toolkit

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRootFromEnv(t *testing.T) {
	tree := makeTree(t)
	t.Setenv(TreeEnv, tree)
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	got, err := Root()
	if err != nil {
		t.Fatal(err)
	}
	if got != tree {
		t.Fatalf("got %s want %s", got, tree)
	}
}

func TestRootFromSetupConfig(t *testing.T) {
	tree := makeTree(t)
	t.Setenv(TreeEnv, "")
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	if _, err := WriteTree(tree); err != nil {
		t.Fatal(err)
	}
	got, err := Root()
	if err != nil {
		t.Fatal(err)
	}
	if got != tree {
		t.Fatalf("got %s want %s", got, tree)
	}
}

func TestRootWalksCwd(t *testing.T) {
	tree := makeTree(t)
	t.Setenv(TreeEnv, "")
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	sub := filepath.Join(tree, "mine", "kits", "x")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(sub); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(cwd) })
	got, err := Root()
	if err != nil {
		t.Fatal(err)
	}
	if got != tree {
		t.Fatalf("got %s want %s", got, tree)
	}
}

func TestRootWalksFromTreeDir(t *testing.T) {
	tree := makeTree(t)
	t.Setenv(TreeEnv, "")
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(tree); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(cwd) })
	got, err := Root()
	if err != nil {
		t.Fatal(err)
	}
	if got != tree {
		t.Fatalf("got %s want %s", got, tree)
	}
}

func TestRootRejectsCatalogAsTree(t *testing.T) {
	tree := makeTree(t)
	cat := filepath.Join(tree, "mine")
	t.Setenv(TreeEnv, cat)
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	_, err := Root()
	if err == nil || !strings.Contains(err.Error(), "looks like a catalog") {
		t.Fatalf("got %v", err)
	}
}

func TestRootEmptyTreeOK(t *testing.T) {
	empty := t.TempDir()
	t.Setenv(TreeEnv, empty)
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	got, err := Root()
	if err != nil {
		t.Fatal(err)
	}
	if got != empty {
		t.Fatalf("got %s want %s", got, empty)
	}
}

func TestWriteTreeCreatesDir(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	dir := filepath.Join(t.TempDir(), "new-tree")
	got, err := WriteTree(dir)
	if err != nil {
		t.Fatal(err)
	}
	st, err := os.Stat(got)
	if err != nil || !st.IsDir() {
		t.Fatalf("expected dir %s: %v", got, err)
	}
	cfg, err := ConfiguredTree()
	if err != nil {
		t.Fatal(err)
	}
	if cfg != got {
		t.Fatalf("config %s want %s", cfg, got)
	}
}

func TestRootMissing(t *testing.T) {
	t.Setenv(TreeEnv, "")
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	empty := t.TempDir()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(empty); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(cwd) })
	_, err = Root()
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != errNeedSetup {
		t.Fatalf("got %v", err)
	}
}

func makeTree(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	p := filepath.Join(root, "mine", "recipes", "agents.yaml")
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte("agents: {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}
