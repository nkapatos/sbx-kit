// Package kitcreds reads credential service ids from kit spec.yaml files
// for host-side hints (sbx secret set). It does not validate kit schema —
// that belongs to sbx (see recipes verify kits).
package kitcreds

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"gopkg.in/yaml.v3"
)

// Need is a host-side secret the user should store with sbx before / at create.
type Need struct {
	Service string   // sbx secret set <Service>
	Envs    []string // in-box env vars (usually proxy-managed sentinels)
	KitName string   // kit name from spec, if present
	KitPath string
}

// ScanSpecs reads kit spec.yaml files and collects credential service ids.
// Supports schemaVersion "1" credentials.sources and a best-effort v2 list form.
func ScanSpecs(kitPaths []string) ([]Need, error) {
	seen := map[string]Need{}
	for _, p := range kitPaths {
		specPath := p
		if st, err := os.Stat(p); err == nil && st.IsDir() {
			specPath = filepath.Join(p, "spec.yaml")
		}
		b, err := os.ReadFile(specPath)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, fmt.Errorf("read %s: %w", specPath, err)
		}
		needs, err := parseSpec(b, p)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", specPath, err)
		}
		for _, n := range needs {
			if prev, ok := seen[n.Service]; ok {
				prev.Envs = mergeUnique(prev.Envs, n.Envs)
				seen[n.Service] = prev
				continue
			}
			seen[n.Service] = n
		}
	}
	out := make([]Need, 0, len(seen))
	for _, n := range seen {
		out = append(out, n)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Service < out[j].Service })
	return out, nil
}

type rawSpec struct {
	Name        string `yaml:"name"`
	Credentials any    `yaml:"credentials"`
}

func parseSpec(b []byte, kitPath string) ([]Need, error) {
	var raw rawSpec
	if err := yaml.Unmarshal(b, &raw); err != nil {
		return nil, err
	}
	if raw.Credentials == nil {
		return nil, nil
	}

	switch c := raw.Credentials.(type) {
	case map[string]any:
		// v1: credentials.sources.<service>.env: [...]
		sources, _ := c["sources"].(map[string]any)
		if sources == nil {
			return nil, nil
		}
		var out []Need
		for svc, val := range sources {
			envs := envsFromSource(val)
			out = append(out, Need{
				Service: svc,
				Envs:    envs,
				KitName: raw.Name,
				KitPath: kitPath,
			})
		}
		return out, nil
	case []any:
		// v2-ish list: [{ service, apiKey: { name } }, ...]
		var out []Need
		for _, item := range c {
			m, ok := item.(map[string]any)
			if !ok {
				continue
			}
			svc, _ := m["service"].(string)
			if svc == "" {
				continue
			}
			var envs []string
			if api, ok := m["apiKey"].(map[string]any); ok {
				if name, ok := api["name"].(string); ok && name != "" {
					envs = append(envs, name)
				}
			}
			out = append(out, Need{
				Service: svc,
				Envs:    envs,
				KitName: raw.Name,
				KitPath: kitPath,
			})
		}
		return out, nil
	default:
		return nil, nil
	}
}

func envsFromSource(val any) []string {
	m, ok := val.(map[string]any)
	if !ok {
		return nil
	}
	raw, ok := m["env"]
	if !ok {
		return nil
	}
	switch e := raw.(type) {
	case []any:
		out := make([]string, 0, len(e))
		for _, v := range e {
			if s, ok := v.(string); ok && s != "" {
				out = append(out, s)
			}
		}
		return out
	case []string:
		return append([]string(nil), e...)
	default:
		return nil
	}
}

func mergeUnique(a, b []string) []string {
	seen := map[string]struct{}{}
	var out []string
	for _, xs := range [][]string{a, b} {
		for _, x := range xs {
			if x == "" {
				continue
			}
			if _, ok := seen[x]; ok {
				continue
			}
			seen[x] = struct{}{}
			out = append(out, x)
		}
	}
	return out
}
