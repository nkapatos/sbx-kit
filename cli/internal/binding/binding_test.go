package binding

import (
	"path/filepath"
	"testing"
)

func TestPutGetRoundTrip(t *testing.T) {
	t.Setenv("XDG_DATA_HOME", t.TempDir())
	t.Setenv("XDG_STATE_HOME", t.TempDir())

	project := filepath.Join(t.TempDir(), "proj")
	if err := Put(Record{
		ProjectDir:  project,
		Agent:       "cursor",
		SandboxName: "sbxk-cursor-abc12345",
		ProfileID:   "sbxk-cursor-abc12345",
	}); err != nil {
		t.Fatal(err)
	}

	got, err := Get(project, "cursor")
	if err != nil {
		t.Fatal(err)
	}
	if got == nil || got.SandboxName != "sbxk-cursor-abc12345" {
		t.Fatalf("got %+v", got)
	}

	byName, err := GetBySandbox("sbxk-cursor-abc12345")
	if err != nil || byName == nil {
		t.Fatalf("byName=%v err=%v", byName, err)
	}

	forProj, err := ListForProject(project)
	if err != nil || len(forProj) != 1 {
		t.Fatalf("ListForProject=%v err=%v", forProj, err)
	}
	if Label(&forProj[0]) != "proj" {
		t.Fatalf("Label=%q", Label(&forProj[0]))
	}

	if err := Delete(project, "cursor"); err != nil {
		t.Fatal(err)
	}
	got, err = Get(project, "cursor")
	if err != nil || got != nil {
		t.Fatalf("expected deleted, got=%v err=%v", got, err)
	}
}
