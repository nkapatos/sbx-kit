package cmd

import (
	"bytes"
	"testing"

	"github.com/nkapatos/sbx-kit/cli/internal/ui"
)

func TestUIFollowsCobraOut(t *testing.T) {
	root := NewRoot()
	out := &bytes.Buffer{}
	errb := &bytes.Buffer{}
	root.SetOut(out)
	root.SetErr(errb)
	root.SetArgs([]string{"version"})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	if UI().Out != out {
		t.Fatal("UI.Out should follow cobra SetOut after Execute")
	}
}

func TestSetUI(t *testing.T) {
	prev := UI()
	t.Cleanup(func() { SetUI(prev) })
	buf := &bytes.Buffer{}
	SetUI(ui.New(buf, buf))
	UI().Header("hi")
	if buf.String() == "" {
		t.Fatal("expected header on injected writer")
	}
}
