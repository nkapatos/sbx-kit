package catalog

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseID(t *testing.T) {
	dir, name, err := ParseID("mine/cursor")
	if err != nil || dir != "mine" || name != "cursor" {
		t.Fatalf("got %q %q %v", dir, name, err)
	}
	if _, _, err := ParseID("cursor"); err == nil {
		t.Fatal("expected error")
	}
	if _, _, err := ParseID("../x/y"); err == nil {
		t.Fatal("expected invalid directory")
	}
	dir, name, err = ParseID("mine/nested/name")
	if err != nil || dir != "mine" || name != "nested/name" {
		t.Fatalf("got %q %q %v", dir, name, err)
	}
}

func TestListAndLookup(t *testing.T) {
	catalogRoot := t.TempDir()
	writeDir(t, catalogRoot, "aa", "cursor: {}\n")
	writeDir(t, catalogRoot, "bb", "shell: {}\n")
	if err := os.MkdirAll(filepath.Join(catalogRoot, "not-a-dir"), 0o755); err != nil {
		t.Fatal(err)
	}

	dirs, err := List(catalogRoot)
	if err != nil {
		t.Fatal(err)
	}
	if len(dirs) != 2 || dirs[0].Name != "aa" || dirs[1].Name != "bb" {
		t.Fatalf("got %+v", dirs)
	}

	d, manifest, ag, err := Lookup(catalogRoot, "aa/cursor")
	if err != nil {
		t.Fatal(err)
	}
	if d.Name != "aa" || manifest.Agents["cursor"].SbxAgent != "" {
		t.Fatalf("dir=%+v manifest=%+v agent=%+v", d, manifest, ag)
	}
	if _, _, _, err := Lookup(catalogRoot, "aa/missing"); err == nil {
		t.Fatal("expected missing recipe")
	}
	if _, _, _, err := Lookup(catalogRoot, "nope/cursor"); err == nil {
		t.Fatal("expected missing directory")
	}
}

func TestFilterDirs(t *testing.T) {
	dirs := []Dir{{Name: "aa"}, {Name: "bb"}}
	got, err := FilterDirs(dirs, "aa")
	if err != nil || len(got) != 1 || got[0].Name != "aa" {
		t.Fatalf("got %+v err=%v", got, err)
	}
	if _, err := FilterDirs(dirs, "nope"); err == nil {
		t.Fatal("expected error")
	}
}

func writeDir(t *testing.T, catalogRoot, name, agentsBody string) {
	t.Helper()
	p := filepath.Join(catalogRoot, name, "recipes", "agents.yaml")
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	body := "agents:\n"
	for _, line := range splitLines(agentsBody) {
		if line == "" {
			continue
		}
		body += "  " + line + "\n"
	}
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func splitLines(s string) []string {
	var out []string
	cur := ""
	for _, r := range s {
		if r == '\n' {
			out = append(out, cur)
			cur = ""
			continue
		}
		cur += string(r)
	}
	if cur != "" {
		out = append(out, cur)
	}
	return out
}
