package catalog

import (
	"reflect"
	"testing"
)

func TestResolveKitsMergesDefaults(t *testing.T) {
	defaults := []string{"agent-workspace"}

	if got := ResolveKits(nil, defaults); !reflect.DeepEqual(got, defaults) {
		t.Fatalf("empty recipe: %v", got)
	}
	if got := ResolveKits([]string{"pi", "mise-workspace"}, defaults); !reflect.DeepEqual(got, []string{"pi", "mise-workspace", "agent-workspace"}) {
		t.Fatalf("merge: %v", got)
	}
	if got := ResolveKits([]string{"mise-workspace", "agent-workspace"}, defaults); !reflect.DeepEqual(got, []string{"mise-workspace", "agent-workspace"}) {
		t.Fatalf("already present: %v", got)
	}
}

func TestResolveKitsEmptyDefaults(t *testing.T) {
	if got := ResolveKits(nil, nil); len(got) != 0 {
		t.Fatalf("empty: %v", got)
	}
	if got := ResolveKits([]string{"mise-workspace"}, nil); !reflect.DeepEqual(got, []string{"mise-workspace"}) {
		t.Fatalf("recipe only: %v", got)
	}
}
