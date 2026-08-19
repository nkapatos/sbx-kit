package recipeverify

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type fakeKitRunner struct {
	verifyCalls []string
	ensureErr   error
	verifyErr   error
}

func (f *fakeKitRunner) EnsureKitVerify() error { return f.ensureErr }

func (f *fakeKitRunner) KitVerify(path string) error {
	f.verifyCalls = append(f.verifyCalls, path)
	return f.verifyErr
}

func TestVerifyRecipeManifest(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "mine")
	agents := filepath.Join(dir, "recipes", "agents.yaml")
	kits := filepath.Join(dir, "kits", "agent-workspace")
	if err := os.MkdirAll(kits, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(kits, "spec.yaml"), []byte("name: agent-workspace\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(agents), 0o755); err != nil {
		t.Fatal(err)
	}
	yaml := `defaults:
  kits: [agent-workspace]
agents:
  cursor:
    sbx_agent: cursor
`
	if err := os.WriteFile(agents, []byte(yaml), 0o644); err != nil {
		t.Fatal(err)
	}

	out := &bytes.Buffer{}
	fk := &fakeKitRunner{}
	err := VerifyRecipe(root, "mine/cursor", Options{Out: out, Runner: fk})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "recipe mine/cursor: ok") {
		t.Fatalf("out: %s", out.String())
	}
	if len(fk.verifyCalls) != 1 {
		t.Fatalf("verify calls: %v", fk.verifyCalls)
	}
}

func TestVerifyRecipeMissingKit(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "mine")
	agents := filepath.Join(dir, "recipes", "agents.yaml")
	if err := os.MkdirAll(filepath.Dir(agents), 0o755); err != nil {
		t.Fatal(err)
	}
	yaml := `agents:
  cursor:
    sbx_agent: cursor
    kits: [missing-kit]
`
	if err := os.WriteFile(agents, []byte(yaml), 0o644); err != nil {
		t.Fatal(err)
	}

	err := VerifyRecipe(root, "mine/cursor", Options{Out: os.Stderr, SkipKits: true})
	if err == nil || !strings.Contains(err.Error(), "missing kit") {
		t.Fatalf("got %v", err)
	}
}

func TestVerifyRecipeSkipKits(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "mine")
	agents := filepath.Join(dir, "recipes", "agents.yaml")
	if err := os.MkdirAll(filepath.Dir(agents), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(agents, []byte("agents:\n  cursor:\n    sbx_agent: cursor\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	fk := &fakeKitRunner{}
	err := VerifyRecipe(root, "mine/cursor", Options{Out: os.Stderr, SkipKits: true, Runner: fk})
	if err != nil {
		t.Fatal(err)
	}
	if len(fk.verifyCalls) != 0 {
		t.Fatalf("expected no kit verify, got %v", fk.verifyCalls)
	}
}
