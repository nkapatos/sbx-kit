package cmd

import (
	"reflect"
	"strings"
	"testing"
)

func TestExtractPassthrough(t *testing.T) {
	got := extractPassthrough([]string{"sbx-kit", "run", "cursor", "--", "--memory", "8g"})
	if len(got) != 2 || got[0] != "--memory" || got[1] != "8g" {
		t.Fatalf("got %v", got)
	}
	if extractPassthrough([]string{"sbx-kit", "run", "cursor"}) != nil {
		t.Fatal("expected nil")
	}
}

func TestPositionalArgsStripsPassthrough(t *testing.T) {
	extra := []string{"--memory", "8g"}
	got := positionalArgs([]string{"cursor", "--memory", "8g"}, extra)
	if !reflect.DeepEqual(got, []string{"cursor"}) {
		t.Fatalf("got %v", got)
	}
	if got := positionalArgs([]string{"cursor"}, nil); !reflect.DeepEqual(got, []string{"cursor"}) {
		t.Fatalf("got %v", got)
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
