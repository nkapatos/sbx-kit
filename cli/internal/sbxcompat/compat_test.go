package sbxcompat

import "testing"

func TestParseVersion(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"0.34.0", "0.34.0"},
		{"v0.34.0", "0.34.0"},
		{"sbx version 0.34.0\n", "0.34.0"},
		{"Version: 0.34.0 (abc)", "0.34.0"},
		{"0.34.0-rc3", "0.34.0"},
	}
	for _, tc := range cases {
		got, err := ParseVersion(tc.in)
		if err != nil {
			t.Fatalf("%q: %v", tc.in, err)
		}
		if got != tc.want {
			t.Fatalf("%q: got %q want %q", tc.in, got, tc.want)
		}
	}
}

func TestCheckMin(t *testing.T) {
	Reset()
	t.Setenv(skipEnv, "")
	if err := Check("0.33.0"); err == nil {
		t.Fatal("expected error for 0.33.0")
	}
	if err := Check("0.34.0"); err != nil {
		t.Fatalf("0.34.0: %v", err)
	}
	if err := Check("1.0.0"); err != nil {
		t.Fatalf("1.0.0: %v", err)
	}
}

func TestSkipEnv(t *testing.T) {
	Reset()
	t.Setenv(skipEnv, "1")
	if err := Check("0.1.0"); err != nil {
		t.Fatalf("skip should allow old sbx: %v", err)
	}
}
