package binding

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/nkapatos/sbx-kit/cli/internal/xdg"
)

const fileName = "bindings.json"

// Record ties a project+agent to an sbx sandbox name and host profile id.
type Record struct {
	ProjectDir  string    `json:"projectDir"`
	Agent       string    `json:"agent"`
	SandboxName string    `json:"sandboxName"`
	ProfileID   string    `json:"profileID"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type store struct {
	Bindings []Record `json:"bindings"`
}

func path() (string, error) {
	dir, err := xdg.StateDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, fileName), nil
}

func load() (*store, error) {
	p, err := path()
	if err != nil {
		return nil, err
	}
	b, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return &store{}, nil
		}
		return nil, err
	}
	var s store
	if err := json.Unmarshal(b, &s); err != nil {
		return nil, fmt.Errorf("parse %s: %w", p, err)
	}
	return &s, nil
}

func save(s *store) error {
	if err := xdg.Ensure(); err != nil {
		return err
	}
	p, err := path()
	if err != nil {
		return err
	}
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, append(b, '\n'), 0o644)
}

// Put upserts a binding for (projectDir, agent).
func Put(rec Record) error {
	abs, err := filepath.Abs(rec.ProjectDir)
	if err != nil {
		return err
	}
	rec.ProjectDir = abs
	rec.UpdatedAt = time.Now().UTC()
	if rec.ProfileID == "" {
		rec.ProfileID = rec.SandboxName
	}

	s, err := load()
	if err != nil {
		return err
	}
	found := false
	for i := range s.Bindings {
		if s.Bindings[i].ProjectDir == rec.ProjectDir && s.Bindings[i].Agent == rec.Agent {
			s.Bindings[i] = rec
			found = true
			break
		}
	}
	if !found {
		s.Bindings = append(s.Bindings, rec)
	}
	return save(s)
}

// Get returns the binding for (projectDir, agent).
func Get(projectDir, agent string) (*Record, error) {
	abs, err := filepath.Abs(projectDir)
	if err != nil {
		return nil, err
	}
	s, err := load()
	if err != nil {
		return nil, err
	}
	for i := range s.Bindings {
		if s.Bindings[i].ProjectDir == abs && s.Bindings[i].Agent == agent {
			rec := s.Bindings[i]
			return &rec, nil
		}
	}
	return nil, nil
}

// GetBySandbox finds a binding by sandbox name.
func GetBySandbox(name string) (*Record, error) {
	s, err := load()
	if err != nil {
		return nil, err
	}
	for i := range s.Bindings {
		if s.Bindings[i].SandboxName == name {
			rec := s.Bindings[i]
			return &rec, nil
		}
	}
	return nil, nil
}

// List returns all bindings.
func List() ([]Record, error) {
	s, err := load()
	if err != nil {
		return nil, err
	}
	out := make([]Record, len(s.Bindings))
	copy(out, s.Bindings)
	return out, nil
}

// Delete removes a binding for (projectDir, agent).
func Delete(projectDir, agent string) error {
	abs, err := filepath.Abs(projectDir)
	if err != nil {
		return err
	}
	s, err := load()
	if err != nil {
		return err
	}
	next := s.Bindings[:0]
	for _, r := range s.Bindings {
		if r.ProjectDir == abs && r.Agent == agent {
			continue
		}
		next = append(next, r)
	}
	s.Bindings = next
	return save(s)
}
