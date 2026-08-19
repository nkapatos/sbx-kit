package catalog

import (
	"os"
	"path/filepath"
	"testing"
)

func TestListKitPaths(t *testing.T) {
	root := t.TempDir()
	kits := filepath.Join(root, "kits")
	if err := os.MkdirAll(filepath.Join(kits, "a"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(kits, "a", "spec.yaml"), []byte("name: a\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(kits, "empty"), 0o755); err != nil {
		t.Fatal(err)
	}

	got, err := ListKitPaths(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0] != filepath.Join(kits, "a") {
		t.Fatalf("got %v", got)
	}
}

func TestRequireKitPath(t *testing.T) {
	root := t.TempDir()
	p := filepath.Join(root, "kits", "x")
	if err := os.MkdirAll(p, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := RequireKitPath(p); err == nil {
		t.Fatal("expected missing spec error")
	}
}
