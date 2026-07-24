package sbxname

import (
	"strings"
	"testing"
)

func TestFromProjectStable(t *testing.T) {
	a := FromProject("cursor", "/Users/me/proj")
	b := FromProject("cursor", "/Users/me/proj")
	if a != b {
		t.Fatalf("unstable: %s vs %s", a, b)
	}
	if !strings.HasPrefix(a, "sbxk-cursor-") {
		t.Fatalf("prefix: %s", a)
	}
	if !Valid(a) {
		t.Fatalf("invalid: %s", a)
	}
}

func TestInjectAndExtract(t *testing.T) {
	name, args := Inject([]string{"--clone", "--memory", "8g"}, "sbxk-cursor-deadbeef")
	if name != "sbxk-cursor-deadbeef" {
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
