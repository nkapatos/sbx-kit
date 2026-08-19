package sbxname

import (
	"strings"
	"testing"
)

func TestNewProfileIDStable(t *testing.T) {
	a := NewProfileID("cursor", "/Users/me/proj")
	b := NewProfileID("cursor", "/Users/me/proj")
	if a != b {
		t.Fatalf("unstable: %s vs %s", a, b)
	}
	if !strings.HasPrefix(a, "sbxk-cursor-") {
		t.Fatalf("prefix: %s", a)
	}
	slash := NewProfileID("mine/cursor", "/Users/me/proj")
	if !strings.Contains(slash, "mine-cursor") {
		t.Fatalf("namespaced: %s", slash)
	}
	if !Valid(a) {
		t.Fatalf("invalid: %s", a)
	}
}

func TestFromDir(t *testing.T) {
	got, err := FromDir("/Users/me/posselt.studio")
	if err != nil {
		t.Fatal(err)
	}
	if got != "posselt.studio" {
		t.Fatalf("got %q", got)
	}
	got, err = FromDir("/tmp/My App!")
	if err != nil {
		t.Fatal(err)
	}
	if got != "My-App" {
		t.Fatalf("got %q", got)
	}
}

func TestInjectAndExtract(t *testing.T) {
	name, args := Inject([]string{"--clone", "--memory", "8g"}, "my-project")
	if name != "my-project" {
		t.Fatalf("name=%s", name)
	}
	if args[0] != "--name" || args[1] != name {
		t.Fatalf("args=%v", args)
	}

	got, rest, ok := ExtractFromArgs([]string{"--clone", "--name", "mine", "--cpus", "4"})
	if !ok || got != "mine" {
		t.Fatalf("extract got=%s ok=%v", got, ok)
	}
	if len(rest) != 3 || rest[0] != "--clone" {
		t.Fatalf("rest=%v", rest)
	}

	name, args = Inject([]string{"--name=custom", "--clone"}, "ignored")
	if name != "custom" || args[1] != "custom" {
		t.Fatalf("name=%s args=%v", name, args)
	}
}
