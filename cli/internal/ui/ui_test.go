package ui

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewHonorsNoColor(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	t.Setenv("TERM", "xterm")
	w := New(&bytes.Buffer{}, &bytes.Buffer{})
	if w.Color() {
		t.Fatal("expected color off when NO_COLOR is set")
	}
}

func TestNewHonorsDumbTerm(t *testing.T) {
	t.Setenv("NO_COLOR", "")
	t.Setenv("TERM", "dumb")
	w := New(&bytes.Buffer{}, &bytes.Buffer{})
	if w.Color() {
		t.Fatal("expected color off when TERM=dumb")
	}
}

func TestHeaderAndDetail(t *testing.T) {
	out := &bytes.Buffer{}
	w := New(out, &bytes.Buffer{})
	w.NoColor = true
	w.Header("check")
	w.Detail("sandbox", "demo")
	w.Detail("recipe", "")
	got := out.String()
	if !strings.Contains(got, "==> check\n") {
		t.Fatalf("header: %q", got)
	}
	if !strings.Contains(got, "sandbox:") || !strings.Contains(got, "demo") {
		t.Fatalf("detail: %q", got)
	}
	if strings.Contains(got, "recipe:") {
		t.Fatalf("empty detail should be skipped: %q", got)
	}
}

func TestWarnGoesToErr(t *testing.T) {
	out := &bytes.Buffer{}
	err := &bytes.Buffer{}
	w := New(out, err)
	w.Warn("sbx ls failed")
	if out.Len() != 0 {
		t.Fatalf("warn leaked to Out: %q", out.String())
	}
	if !strings.Contains(err.String(), "warning: sbx ls failed") {
		t.Fatalf("err: %q", err.String())
	}
}

func TestEmpty(t *testing.T) {
	out := &bytes.Buffer{}
	w := New(out, &bytes.Buffer{})
	w.Empty("directories", "add one:  sbx-kit catalog add <url>")
	got := out.String()
	if !strings.Contains(got, "(no directories)\n") {
		t.Fatalf("empty: %q", got)
	}
	if !strings.Contains(got, "add one:") {
		t.Fatalf("hint: %q", got)
	}
}

func TestTable(t *testing.T) {
	out := &bytes.Buffer{}
	w := New(out, &bytes.Buffer{})
	if err := w.Table([]string{"DIR", "ORIGIN"}, [][]string{{"mine", "local"}}); err != nil {
		t.Fatal(err)
	}
	got := out.String()
	if !strings.Contains(got, "DIR") || !strings.Contains(got, "mine") {
		t.Fatalf("table: %q", got)
	}
}
