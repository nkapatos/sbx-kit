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
	AgentCatalogName string // catalog recipe id
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
	// ConfirmCreate prompts before creating when true.
	ConfirmCreate bool
	ConfirmFn     func(agent, path, name string) (bool, error)
	// SandboxName is the friendly sbx --name (optional override).
	SandboxName string
	// ProfileID overrides the opaque vault id (e.g. upgrade reuses binding).
	ProfileID string
	// NamePromptFn asks for a friendly name; defaultName is the cwd basename.
	// Called only when ConfirmCreate and SandboxName is empty.
	NamePromptFn func(defaultName string) (string, error)
	// StaleArchiveFn is called when a vault archive exists but the box does not.
	// restore=true → import after create; discard=true → delete archive; both false → abort.
	StaleArchiveFn func(profileID string) (restore, discard bool, err error)
	Runner         *sbxutil.Runner
}

// Result is returned after create/attach preparation.
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

	passthroughName, restExtra, passthroughHasName := sbxname.ExtractFromArgs(o.Extra)

	// Opaque vault id: stable per recipe+path (or reused from binding / caller).
	profileID := o.ProfileID
	var existing *binding.Record
	if rec, err := binding.Get(absProject, o.AgentCatalogName); err == nil && rec != nil {
		existing = rec
		if profileID == "" {
			profileID = rec.ProfileID
		}
	}
	if profileID == "" {
		profileID = sbxname.NewProfileID(o.AgentCatalogName, absProject)
	}

	// Friendly sbx name (what sbx ls shows).
	name, err := resolveFriendlyName(o, absProject, existing, passthroughHasName, passthroughName)
	if err != nil {
		return nil, err
	}
	if !sbxname.Valid(name) {
		return nil, fmt.Errorf("invalid sandbox name %q", name)
	}

	extra := append([]string{"--name", name}, restExtra...)
	extra = appendResourceFlags(extra, o.Resources)

	template := resolveTemplate(r, o.ImageName, o.TemplateFallback, o.TemplateOverride)

	env := os.Environ()
	if o.Resources.RootSize != "" {
		env = setEnv(env, "DOCKER_SANDBOXES_ROOT_SIZE", o.Resources.RootSize)
	}
	if o.Resources.DockerSize != "" {
		env = setEnv(env, "DOCKER_SANDBOXES_DOCKER_SIZE", o.Resources.DockerSize)
	}

	fmt.Printf("==> resources profile=%s memory=%s cpus=%s root=%s docker=%s\n",
		o.ResourcesProfile, o.Resources.Memory, o.Resources.CPUs, o.Resources.RootSize, o.Resources.DockerSize)
	fmt.Printf("==> sandbox name=%s  profile=%s\n", name, profileID)

	res := &Result{SandboxName: name, ProfileID: profileID, ProjectDir: absProject, Label: name}

	exists, err := r.Exists(name)
	if err != nil {
		fmt.Printf("==> warning: could not list sandboxes: %v\n", err)
		exists = false
	}

	restore := o.RestoreState
	if !exists {
		hasArch, err := statexfer.HasArchive(profileID)
		if err != nil {
			return res, err
		}
		if hasArch && !restore {
			if o.StaleArchiveFn != nil {
				wantRestore, discard, err := o.StaleArchiveFn(profileID)
				if err != nil {
					return res, err
				}
				if discard {
					if err := statexfer.DiscardArchive(profileID); err != nil {
						return res, err
					}
				} else if wantRestore {
					restore = true
				} else {
					return res, fmt.Errorf("create cancelled (archive kept at profile %s)", profileID)
				}
			} else if o.ConfirmCreate {
				return res, fmt.Errorf("archive exists for profile %s; pass --restore-state, discard it, or abort", profileID)
			} else {
				fmt.Printf("==> note: archive exists for profile %s (not restoring; pass --restore-state)\n", profileID)
			}
		}
	}

	if exists && o.CreateOnly {
		return res, fmt.Errorf("sandbox %q already exists\n  attach: sbx-kit run --name %s\n  or:    sbx-kit run   (from the project dir)\n  rename: sbx-kit rm --keep-state then recreate with a new name + --restore-state",
			name, name)
	}

	if !exists {
		ok, err := confirmCreate(o, absProject, name)
		if err != nil {
			return res, err
		}
		if !ok {
			return res, fmt.Errorf("create cancelled")
		}
	}

	// Persist mapping before sbx create/run (kit-owned; sbx only sees --name).
	if err := binding.Put(binding.Record{
		ProjectDir:  absProject,
		Agent:       o.AgentCatalogName,
		SandboxName: name,
		ProfileID:   profileID,
	}); err != nil {
		return res, fmt.Errorf("save binding: %w", err)
	}

	if restore {
		has, err := statexfer.HasArchive(profileID)
		if err != nil {
			return res, err
		}
		if !has {
			fmt.Printf("==> warning: --restore-state set but no archive at profile %s\n", profileID)
		}
		if !exists {
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

	if exists {
		fmt.Printf("==> re-attaching existing sandbox %s\n", name)
		return res, r.RunEnv(env, "run", "--name", name)
	}

	runArgs := buildArgs("run", o.SbxAgent, template, o.KitPaths, extra, absProject)
	return res, r.RunEnv(env, runArgs...)
}

func resolveFriendlyName(o Opts, absProject string, existing *binding.Record, passthroughHasName bool, passthroughName string) (string, error) {
	if passthroughHasName {
		return passthroughName, nil
	}
	if o.SandboxName != "" {
		return o.SandboxName, nil
	}

	def, err := sbxname.FromDir(absProject)
	if err != nil {
		return "", err
	}
	// Prefer previous friendly name as default when recreating the same binding.
	if existing != nil && existing.SandboxName != "" && !strings.HasPrefix(existing.SandboxName, "sbxk-") {
		def = existing.SandboxName
	} else if existing != nil && existing.SandboxName != "" {
		// Legacy binding used opaque id as sbx name — offer dirname default instead.
		def, err = sbxname.FromDir(absProject)
		if err != nil {
			return "", err
		}
	}

	if o.ConfirmCreate && o.NamePromptFn != nil {
		return o.NamePromptFn(def)
	}
	return def, nil
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
