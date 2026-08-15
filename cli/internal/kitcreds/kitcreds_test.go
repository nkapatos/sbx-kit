package kitcreds

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanSpecsV1Sources(t *testing.T) {
	dir := t.TempDir()
	kit := filepath.Join(dir, "example")
	if err := os.MkdirAll(kit, 0o755); err != nil {
		t.Fatal(err)
	}
	spec := []byte(`
schemaVersion: "1"
kind: mixin
name: example
credentials:
  sources:
    deepseek:
      env:
        - DEEPSEEK_API_KEY
    openai:
      env:
        - OPENAI_API_KEY
`)
	if err := os.WriteFile(filepath.Join(kit, "spec.yaml"), spec, 0o644); err != nil {
		t.Fatal(err)
	}

	needs, err := ScanSpecs([]string{kit})
	if err != nil {
		t.Fatal(err)
	}
	if len(needs) != 2 {
		t.Fatalf("got %#v", needs)
	}
	by := map[string]Need{}
	for _, n := range needs {
		by[n.Service] = n
	}
	if by["deepseek"].Envs[0] != "DEEPSEEK_API_KEY" || by["openai"].Envs[0] != "OPENAI_API_KEY" {
		t.Fatalf("envs %#v", needs)
	}
}

func TestScanSpecsV2List(t *testing.T) {
	dir := t.TempDir()
	kit := filepath.Join(dir, "x")
	if err := os.MkdirAll(kit, 0o755); err != nil {
		t.Fatal(err)
	}
	spec := []byte(`
name: x
credentials:
  - service: my-service
    apiKey:
      name: MY_SERVICE_API_KEY
`)
	if err := os.WriteFile(filepath.Join(kit, "spec.yaml"), spec, 0o644); err != nil {
		t.Fatal(err)
	}
	needs, err := ScanSpecs([]string{kit})
	if err != nil {
		t.Fatal(err)
	}
	if len(needs) != 1 || needs[0].Service != "my-service" {
		t.Fatalf("got %#v", needs)
	}
}
