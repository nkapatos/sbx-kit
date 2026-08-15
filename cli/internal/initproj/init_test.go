package initproj

import (
	"strings"
	"testing"

	"github.com/nkapatos/sbx-kit/cli/internal/catalog"
)

func TestBuildSectionHubOmitsMise(t *testing.T) {
	got := buildSection("cursor", catalog.Agent{SbxAgent: "cursor"}, []string{"agent-workspace"})
	if strings.Contains(got, "mise install") {
		t.Fatalf("hub recipe should not tell humans to mise install:\n%s", got)
	}
	if !strings.Contains(got, "floor.md") {
		t.Fatalf("expected floor.md pointer:\n%s", got)
	}
}

func TestBuildSectionCustomIncludesMise(t *testing.T) {
	got := buildSection("kit-cursor", catalog.Agent{
		SbxAgent:  "cursor",
		ImageName: "sbx-kit-cursor",
		Kits:      []string{"mise-workspace"},
	}, []string{"agent-workspace"})
	if !strings.Contains(got, "mise install") {
		t.Fatalf("custom+mise recipe should mention mise install:\n%s", got)
	}
}
