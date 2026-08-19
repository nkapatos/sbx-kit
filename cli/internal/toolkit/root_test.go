package toolkit

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRootFromEnv(t *testing.T) {
	catalogRoot := makeCatalog(t)
	t.Setenv(CatalogEnv, catalogRoot)
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	got, err := Root()
	if err != nil {
		t.Fatal(err)
	}
	if got != catalogRoot {
		t.Fatalf("got %s want %s", got, catalogRoot)
	}
}

func TestRootFromSetupConfig(t *testing.T) {
	catalogRoot := makeCatalog(t)
	t.Setenv(CatalogEnv, "")
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	if _, err := WriteCatalog(catalogRoot); err != nil {
		t.Fatal(err)
	}
	got, err := Root()
	if err != nil {
		t.Fatal(err)
	}
	if got != catalogRoot {
		t.Fatalf("got %s want %s", got, catalogRoot)
	}
}

func TestRootWalksCwd(t *testing.T) {
	catalogRoot := makeCatalog(t)
	t.Setenv(CatalogEnv, "")
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	sub := filepath.Join(catalogRoot, "mine", "kits", "x")
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
	if got != catalogRoot {
		t.Fatalf("got %s want %s", got, catalogRoot)
	}
}

func TestRootWalksFromCatalogDir(t *testing.T) {
	catalogRoot := makeCatalog(t)
	t.Setenv(CatalogEnv, "")
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(catalogRoot); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(cwd) })
	got, err := Root()
	if err != nil {
		t.Fatal(err)
	}
	if got != catalogRoot {
		t.Fatalf("got %s want %s", got, catalogRoot)
	}
}

func TestRootRejectsSourceAsCatalog(t *testing.T) {
	catalogRoot := makeCatalog(t)
	src := filepath.Join(catalogRoot, "mine")
	t.Setenv(CatalogEnv, src)
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	_, err := Root()
	if err == nil || !strings.Contains(err.Error(), "looks like a source") {
		t.Fatalf("got %v", err)
	}
}

func TestRootEmptyCatalogOK(t *testing.T) {
	empty := t.TempDir()
	t.Setenv(CatalogEnv, empty)
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	got, err := Root()
	if err != nil {
		t.Fatal(err)
	}
	if got != empty {
		t.Fatalf("got %s want %s", got, empty)
	}
}

func TestWriteCatalogCreatesDir(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	dir := filepath.Join(t.TempDir(), "new-catalog")
	got, err := WriteCatalog(dir)
	if err != nil {
		t.Fatal(err)
	}
	st, err := os.Stat(got)
	if err != nil || !st.IsDir() {
		t.Fatalf("expected dir %s: %v", got, err)
	}
	cfg, err := ConfiguredCatalog()
	if err != nil {
		t.Fatal(err)
	}
	if cfg != got {
		t.Fatalf("config %s want %s", cfg, got)
	}
}

func TestRootMissing(t *testing.T) {
	t.Setenv(CatalogEnv, "")
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

func makeCatalog(t *testing.T) string {
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
