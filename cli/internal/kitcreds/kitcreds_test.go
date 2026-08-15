package kitcreds

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanSpecsV1Sources(t *testing.T) {
	dir := t.TempDir()
	kit := filepath.Join(dir, "deepseek-creds")
	if err := os.MkdirAll(kit, 0o755); err != nil {
		t.Fatal(err)
	}
	spec := []byte(`
schemaVersion: "1"
kind: mixin
name: deepseek-creds
credentials:
  sources:
    deepseek:
      env:
        - DEEPSEEK_API_KEY
`)
	if err := os.WriteFile(filepath.Join(kit, "spec.yaml"), spec, 0o644); err != nil {
		t.Fatal(err)
	}

	needs, err := ScanSpecs([]string{kit})
	if err != nil {
		t.Fatal(err)
	}
	if len(needs) != 1 || needs[0].Service != "deepseek" {
		t.Fatalf("got %#v", needs)
	}
	if len(needs[0].Envs) != 1 || needs[0].Envs[0] != "DEEPSEEK_API_KEY" {
		t.Fatalf("envs %#v", needs[0].Envs)
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
