package cmd

import (
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
