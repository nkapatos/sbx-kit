package gitdir

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestDirName(t *testing.T) {
	cases := map[string]string{
		"https://github.com/foo/bar.git": "bar",
		"https://github.com/foo/bar":     "bar",
		"http://example.com/x/y.git/":    "y",
		"git@github.com:foo/bar.git":     "bar",
		"ssh://git@host/org/repo.git":    "repo",
	}
	for in, want := range cases {
		if got := DirName(in); got != want {
			t.Errorf("%s: got %q want %q", in, got, want)
		}
	}
}

func TestCloneStatus(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	src := t.TempDir()
	run := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = src
		cmd.Env = append(os.Environ(),
			"GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t",
			"GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	run("init", "-q")
	run("config", "user.email", "t@t")
	run("config", "user.name", "t")
	if err := os.WriteFile(filepath.Join(src, "README"), []byte("x\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	run("add", "README")
	run("commit", "-qm", "init")

	dest := filepath.Join(t.TempDir(), "bar")
	if err := Clone(src, dest); err != nil {
		t.Fatal(err)
	}
	if !IsRepo(dest) {
		t.Fatal("expected clone to be a repo")
	}
	st, err := Status(dest)
	if err != nil {
		t.Fatal(err)
	}
	if st == "" {
		t.Fatal("empty status")
	}
}
