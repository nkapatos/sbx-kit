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
	if got := ResolveKits([]string{"pi", "deepseek-creds"}, defaults); !reflect.DeepEqual(got, []string{"pi", "deepseek-creds", "agent-workspace"}) {
		t.Fatalf("merge: %v", got)
	}
	if got := ResolveKits([]string{"mise-workspace", "agent-workspace"}, defaults); !reflect.DeepEqual(got, []string{"mise-workspace", "agent-workspace"}) {
		t.Fatalf("already present: %v", got)
	}
}
