package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nkapatos/sbx-kit/cli/internal/binding"
	"github.com/nkapatos/sbx-kit/cli/internal/kitcreds"
	"github.com/nkapatos/sbx-kit/cli/internal/ui"
)

func TestRenderCheck(t *testing.T) {
	buf := &bytes.Buffer{}
	w := ui.New(buf, buf)
	renderCheck(w, &CheckResult{
		SandboxName: "proj",
		RecipeID:    "mine/shell",
		Profile:     "sbxk-mine-shell-abcd1234",
		Project:     "/tmp/proj",
		SbxAgent:    "shell",
		CredNeeds: []kitcreds.Need{{
			Service: "openai",
			Envs:    []string{"OPENAI_API_KEY"},
			KitName: "agent-workspace",
		}},
		SbxOnPath: true,
	})
	out := buf.String()
	for _, want := range []string{"check", "proj", "mine/shell", "openai", "OPENAI_API_KEY"} {
		if !strings.Contains(out, want) {
			t.Fatalf("missing %q in:\n%s", want, out)
		}
	}
}

func TestRenderBindingsEmpty(t *testing.T) {
	buf := &bytes.Buffer{}
	w := ui.New(buf, buf)
	renderBindings(w, &BindingsResult{Share: "/share", State: "/state"})
	if !strings.Contains(buf.String(), "(no bindings)") {
		t.Fatalf("got %s", buf.String())
	}
}

func TestComputeBindingsFiltersProject(t *testing.T) {
	t.Setenv("XDG_DATA_HOME", t.TempDir())
	t.Setenv("XDG_STATE_HOME", t.TempDir())

	a := filepath.Join(t.TempDir(), "a")
	b := filepath.Join(t.TempDir(), "b")
	if err := os.MkdirAll(a, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(b, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := binding.Put(binding.Record{ProjectDir: a, Agent: "mine/shell", SandboxName: "a", ProfileID: "p-a"}); err != nil {
		t.Fatal(err)
	}
	if err := binding.Put(binding.Record{ProjectDir: b, Agent: "mine/cursor", SandboxName: "b", ProfileID: "p-b"}); err != nil {
		t.Fatal(err)
	}

	res, err := computeBindings(a)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Rows) != 1 || res.Rows[0].Name != "a" {
		t.Fatalf("rows=%+v", res.Rows)
	}
	if res.Rows[0].Status == "" {
		t.Fatal("expected a status even when sbx ls fails")
	}
}
