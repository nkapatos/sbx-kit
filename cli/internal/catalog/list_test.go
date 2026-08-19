package catalog

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseID(t *testing.T) {
	cat, name, err := ParseID("mine/cursor")
	if err != nil || cat != "mine" || name != "cursor" {
		t.Fatalf("got %q %q %v", cat, name, err)
	}
	if _, _, err := ParseID("cursor"); err == nil {
		t.Fatal("expected error")
	}
	if _, _, err := ParseID("../x/y"); err == nil {
		t.Fatal("expected invalid source")
	}
	cat, name, err = ParseID("mine/nested/name")
	if err != nil || cat != "mine" || name != "nested/name" {
		t.Fatalf("got %q %q %v", cat, name, err)
	}
}

func TestListAndLookup(t *testing.T) {
	tree := t.TempDir()
	writeCat(t, tree, "aa", "cursor: {}\n")
	writeCat(t, tree, "bb", "shell: {}\n")
	if err := os.MkdirAll(filepath.Join(tree, "not-a-catalog"), 0o755); err != nil {
		t.Fatal(err)
	}

	srcs, err := List(tree)
	if err != nil {
		t.Fatal(err)
	}
	if len(srcs) != 2 || srcs[0].Name != "aa" || srcs[1].Name != "bb" {
		t.Fatalf("got %+v", srcs)
	}

	src, c, ag, err := Lookup(tree, "aa/cursor")
	if err != nil {
		t.Fatal(err)
	}
	if src.Name != "aa" || c.Agents["cursor"].SbxAgent != "" {
		t.Fatalf("src=%+v cat=%+v agent=%+v", src, c, ag)
	}
	if _, _, _, err := Lookup(tree, "aa/missing"); err == nil {
		t.Fatal("expected missing recipe")
	}
	if _, _, _, err := Lookup(tree, "nope/cursor"); err == nil {
		t.Fatal("expected missing source")
	}
}

func writeCat(t *testing.T, tree, name, agentsBody string) {
	t.Helper()
	p := filepath.Join(tree, name, "recipes", "agents.yaml")
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
