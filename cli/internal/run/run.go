package run

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nkapatos/sbx-kit/cli/internal/binding"
	"github.com/nkapatos/sbx-kit/cli/internal/resources"
	"github.com/nkapatos/sbx-kit/cli/internal/sbxname"
	"github.com/nkapatos/sbx-kit/cli/internal/sbxutil"
	"github.com/nkapatos/sbx-kit/cli/internal/statexfer"
	"github.com/nkapatos/sbx-kit/cli/internal/xdg"
)

type Opts struct {
	Root             string
	AgentCatalogName string // sbx-kit catalog key (e.g. cursor)
	SbxAgent         string
	ImageName        string
	TemplateFallback string
	TemplateOverride string
	KitPaths         []string
	ProjectDir       string
	Extra            []string
	Resources        *resources.Profile
	ResourcesProfile string
	RestoreState     bool
	// CreateOnly refuses to attach when the sandbox already exists.
	CreateOnly bool
	// ConfirmCreate prompts before creating a new sandbox when true.
	ConfirmCreate bool
	// ConfirmFn asks whether to create; required when ConfirmCreate is true.
	ConfirmFn func(agent, path, name string) (bool, error)
	Runner    *sbxutil.Runner
}

// Result is returned after argv is prepared (and after create/restore when applicable).
type Result struct {
	SandboxName string
	ProfileID   string
	ProjectDir  string
	Label       string
}

func Sbx(o Opts) (*Result, error) {
	r := o.Runner
	if r == nil {
		r = sbxutil.Default()
	}
	if _, err := r.LookPath(); err != nil {
		return nil, fmt.Errorf("sbx not found on PATH")
	}
	if err := xdg.Ensure(); err != nil {
		return nil, err
	}

	absProject, err := filepath.Abs(o.ProjectDir)
	if err != nil {
		return nil, err
	}
	st, err := os.Stat(absProject)
	if err != nil || !st.IsDir() {
		return nil, fmt.Errorf("not a directory: %s", o.ProjectDir)
	}

	generated := sbxname.FromProject(o.AgentCatalogName, absProject)
	passthroughName, restExtra, passthroughHasName := sbxname.ExtractFromArgs(o.Extra)

	name := generated
	profileID := generated
	if passthroughHasName {
		name = passthroughName
		profileID = passthroughName
	} else if rec, err := binding.Get(absProject, o.AgentCatalogName); err == nil && rec != nil {
		name = rec.SandboxName
		profileID = rec.ProfileID
	}

	if !sbxname.Valid(name) {
		return nil, fmt.Errorf("invalid sandbox name %q", name)
	}

	extra := append([]string{"--name", name}, restExtra...)
	extra = appendResourceFlags(extra, o.Resources)

	if err := binding.Put(binding.Record{
		ProjectDir:  absProject,
		Agent:       o.AgentCatalogName,
		SandboxName: name,
		ProfileID:   profileID,
	}); err != nil {
		return nil, fmt.Errorf("save binding: %w", err)
	}

	template := resolveTemplate(r, o.ImageName, o.TemplateFallback, o.TemplateOverride)

	env := os.Environ()
	if o.Resources.RootSize != "" {
		env = setEnv(env, "DOCKER_SANDBOXES_ROOT_SIZE", o.Resources.RootSize)
	}
	if o.Resources.DockerSize != "" {
		env = setEnv(env, "DOCKER_SANDBOXES_DOCKER_SIZE", o.Resources.DockerSize)
	}

	label := filepath.Base(absProject)
	fmt.Printf("==> resources profile=%s memory=%s cpus=%s root=%s docker=%s\n",
		o.ResourcesProfile, o.Resources.Memory, o.Resources.CPUs, o.Resources.RootSize, o.Resources.DockerSize)
	fmt.Printf("==> sandbox name=%s  label=%s  profile=%s\n", name, label, profileID)

	res := &Result{SandboxName: name, ProfileID: profileID, ProjectDir: absProject, Label: label}

	exists, err := r.Exists(name)
	if err != nil {
		fmt.Printf("==> warning: could not list sandboxes: %v\n", err)
		exists = false
	}

	if o.RestoreState {
		has, err := statexfer.HasArchive(profileID)
		if err != nil {
			return res, err
		}
		if !has {
			fmt.Printf("==> warning: --restore-state set but no archive at profile %s\n", profileID)
		}

		if exists && o.CreateOnly {
			return res, fmt.Errorf("sandbox %s already exists for recipe %s (%s)\n  attach:  sbx-kit run --name %s\n  recreate with state: sbx-kit upgrade --agent %s --path %s",
				name, o.AgentCatalogName, label, name, o.AgentCatalogName, absProject)
		}

		if !exists {
			ok, err := confirmCreate(o, absProject, name)
			if err != nil {
				return res, err
			}
			if !ok {
				return res, fmt.Errorf("create cancelled")
			}
			createArgs := buildArgs("create", o.SbxAgent, template, o.KitPaths, extra, absProject)
			if err := r.RunEnv(env, createArgs...); err != nil {
				return res, err
			}
		}
		if has {
			if err := statexfer.Import(r, name, profileID); err != nil {
				return res, err
			}
		}
		return res, r.RunEnv(env, "run", "--name", name)
	}

	// Re-attach must not pass --kit/--template (sbx rejects those on existing sandboxes).
	if exists {
		if o.CreateOnly {
			return res, fmt.Errorf("sandbox %s already exists for recipe %s (%s)\n  attach: sbx-kit run --name %s\n  or:    sbx-kit run   (from the project dir)",
				name, o.AgentCatalogName, label, name)
		}
		fmt.Printf("==> re-attaching existing sandbox %s (%s)\n", name, label)
		return res, r.RunEnv(env, "run", "--name", name)
	}

	ok, err := confirmCreate(o, absProject, name)
	if err != nil {
		return res, err
	}
	if !ok {
		return res, fmt.Errorf("create cancelled")
	}

	runArgs := buildArgs("run", o.SbxAgent, template, o.KitPaths, extra, absProject)
	return res, r.RunEnv(env, runArgs...)
}

func confirmCreate(o Opts, absProject, name string) (bool, error) {
	if !o.ConfirmCreate {
		return true, nil
	}
	if o.ConfirmFn == nil {
		return false, fmt.Errorf("create confirmation required but no prompt configured (pass --yes)")
	}
	return o.ConfirmFn(o.AgentCatalogName, absProject, name)
}

func buildArgs(verb, sbxAgent, template string, kits, extra []string, project string) []string {
	args := []string{verb, sbxAgent, "--template", template}
	for _, k := range kits {
		args = append(args, "--kit", k)
	}
	args = append(args, extra...)
	args = append(args, project)
	return args
}

func resolveTemplate(r *sbxutil.Runner, imageName, fallback, override string) string {
	if override != "" {
		return override
	}
	out, err := r.Output("template", "ls")
	if err != nil {
		return fallback
	}
	best := ""
	for _, line := range strings.Split(out, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		repo, tag := fields[0], fields[1]
		parts := strings.Split(repo, "/")
		if parts[len(parts)-1] != imageName {
			continue
		}
		ref := repo + ":" + tag
		if best == "" || strings.Count(repo, "/") > strings.Count(best, "/") {
			best = ref
		}
		if strings.HasPrefix(repo, "docker.io/") {
			best = ref
		}
	}
	if best != "" {
		return best
	}
	return fallback
}

func appendResourceFlags(extra []string, res *resources.Profile) []string {
	hasMemory, hasCPUs := false, false
	for _, a := range extra {
		if a == "--memory" || a == "-m" {
			hasMemory = true
		}
		if a == "--cpus" {
			hasCPUs = true
		}
	}
	var prepend []string
	if !hasMemory {
		prepend = append(prepend, "--memory", res.Memory)
	}
	if !hasCPUs {
		prepend = append(prepend, "--cpus", res.CPUs)
	}
	return append(prepend, extra...)
}

func setEnv(env []string, key, value string) []string {
	prefix := key + "="
	out := make([]string, 0, len(env)+1)
	found := false
	for _, e := range env {
		if strings.HasPrefix(e, prefix) {
			out = append(out, prefix+value)
			found = true
			continue
		}
		out = append(out, e)
	}
	if !found {
		out = append(out, prefix+value)
	}
	return out
}
