package cmd

import (
	"strings"
	"testing"
)

func TestExtractPassthrough(t *testing.T) {
	got := extractPassthrough([]string{"sbx-kit", "run", "--agent", "cursor", "--", "--memory", "8g"})
	if len(got) != 2 || got[0] != "--memory" || got[1] != "8g" {
		t.Fatalf("got %v", got)
	}
	if extractPassthrough([]string{"sbx-kit", "run", "--agent", "cursor"}) != nil {
		t.Fatal("expected nil")
	}
}

func TestSoleBindingRecordErrors(t *testing.T) {
	t.Setenv("XDG_DATA_HOME", t.TempDir())
	t.Setenv("XDG_STATE_HOME", t.TempDir())

	dir := t.TempDir()
	_, err := soleBindingRecord(dir)
	if err == nil || !strings.Contains(err.Error(), "no sandbox bound") {
		t.Fatalf("expected no binding error, got %v", err)
	}
}
