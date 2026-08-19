package overlay

import (
	"strings"
	"testing"
)

type fakeRunner struct {
	cps   [][2]string
	execs [][]string
	fail  string
}

func (f *fakeRunner) Cp(src, dst string) error {
	f.cps = append(f.cps, [2]string{src, dst})
	return nil
}

func (f *fakeRunner) Exec(sandbox string, command ...string) error {
	f.execs = append(f.execs, append([]string{sandbox}, command...))
	if f.fail != "" && strings.Contains(strings.Join(command, " "), f.fail) {
		return errFail{}
	}
	return nil
}

type errFail struct{}

func (errFail) Error() string { return "fail" }

func TestInstallCopiesOverlay(t *testing.T) {
	fk := &fakeRunner{}
	if err := Install(fk, "box", nil); err != nil {
		t.Fatal(err)
	}
	if len(fk.execs) < 2 {
		t.Fatalf("execs=%d", len(fk.execs))
	}
	var sawHelper, sawManifest, sawCLI, sawContext bool
	for _, c := range fk.cps {
		if strings.Contains(c[1], "/bin/"+Helper) {
			sawHelper = true
		}
		if strings.HasSuffix(c[1], "state.manifest") {
			sawManifest = true
		}
		if strings.HasSuffix(c[1], "/cli.md") {
			sawCLI = true
		}
		if strings.HasSuffix(c[1], "/context.md") {
			sawContext = true
		}
	}
	if !sawHelper || !sawManifest || !sawCLI || !sawContext {
		t.Fatalf("cps=%v", fk.cps)
	}
}

func TestRenderCLIHasCurrentVerbs(t *testing.T) {
	s, err := renderCLI()
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"sbx-kit box run", "box state", "floor.md", "/etc/sbx-kit/context.md"} {
		if !strings.Contains(s, want) {
			t.Fatalf("missing %q in:\n%s", want, s)
		}
	}
	if strings.Contains(s, "sbx-kit run ") || strings.Contains(s, "agent-workspace") {
		t.Fatalf("stale kit wording in cli.md:\n%s", s)
	}
}

func TestInstallMkdirFailure(t *testing.T) {
	fk := &fakeRunner{fail: "mkdir"}
	if err := Install(fk, "box", nil); err == nil {
		t.Fatal("expected error")
	}
}

func TestHelperIsAgentAgnostic(t *testing.T) {
	b, err := files.ReadFile("files/sbx-kit-state")
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)
	if strings.Contains(s, "agent-workspace") || strings.Contains(s, ".cursor") {
		t.Fatalf("helper still mentions kit/agent paths:\n%s", s)
	}
	if !strings.Contains(s, "overlay not installed") {
		t.Fatal("expected overlay wording in helper")
	}
}

func TestManifestIsPortableOnly(t *testing.T) {
	b, err := files.ReadFile("files/state.manifest")
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)
	if strings.Contains(s, ".cursor") || strings.Contains(s, ".pi") {
		t.Fatalf("manifest should not include agent homes:\n%s", s)
	}
	if !strings.Contains(s, ShareDir+"/portable") {
		t.Fatalf("missing portable include:\n%s", s)
	}
}

func TestContextIndexPointsAtDocs(t *testing.T) {
	b, err := files.ReadFile("files/context.md")
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)
	for _, want := range []string{"/etc/sbx-kit/floor.md", "/etc/sbx-kit/cli.md", "sbx-kit land", "agentContext"} {
		if !strings.Contains(s, want) {
			t.Fatalf("missing %q in context.md:\n%s", want, s)
		}
	}
	if strings.Contains(s, "agent-workspace") || strings.Contains(s, ".cursor") {
		t.Fatalf("context.md still kit/agent-branded:\n%s", s)
	}
}
