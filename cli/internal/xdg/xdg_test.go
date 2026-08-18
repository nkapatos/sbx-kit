package xdg

import (
	"os"
	"path/filepath"
	"testing"
)

func TestShareAndStateHonorEnv(t *testing.T) {
	t.Setenv("XDG_DATA_HOME", filepath.Join(t.TempDir(), "data"))
	t.Setenv("XDG_STATE_HOME", filepath.Join(t.TempDir(), "state"))
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(t.TempDir(), "config"))

	share, err := ShareDir()
	if err != nil {
		t.Fatal(err)
	}
	state, err := StateDir()
	if err != nil {
		t.Fatal(err)
	}
	cfg, err := ConfigDir()
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(share) != "sbx-kit" || filepath.Base(state) != "sbx-kit" || filepath.Base(cfg) != "sbx-kit" {
		t.Fatalf("unexpected dirs share=%s state=%s config=%s", share, state, cfg)
	}
}

func TestEnsureCreatesDirs(t *testing.T) {
	data := t.TempDir()
	state := t.TempDir()
	t.Setenv("XDG_DATA_HOME", data)
	t.Setenv("XDG_STATE_HOME", state)

	if err := Ensure(); err != nil {
		t.Fatal(err)
	}
	profiles, err := ProfilesDir()
	if err != nil {
		t.Fatal(err)
	}
	for _, d := range []string{
		filepath.Join(data, "sbx-kit"),
		filepath.Join(state, "sbx-kit"),
		profiles,
	} {
		st, err := os.Stat(d)
		if err != nil || !st.IsDir() {
			t.Fatalf("missing dir %s: %v", d, err)
		}
	}

	arch, err := ProfileArchive("demo")
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(profiles, "demo", "state.tgz")
	if arch != want {
		t.Fatalf("archive=%s want=%s", arch, want)
	}
}
