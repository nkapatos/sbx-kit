package run

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/nkapatos/sbx-kit/cli/internal/resources"
)

type Opts struct {
	Root             string
	SbxAgent         string
	ImageName        string
	TemplateFallback string
	TemplateOverride string
	KitPaths         []string
	ProjectDir       string
	Extra            []string
	Resources        *resources.Profile
	ResourcesProfile string
}

func Sbx(o Opts) error {
	if _, err := exec.LookPath("sbx"); err != nil {
		return fmt.Errorf("sbx not found on PATH")
	}

	absProject, err := filepath.Abs(o.ProjectDir)
	if err != nil {
		return err
	}
	st, err := os.Stat(absProject)
	if err != nil || !st.IsDir() {
		return fmt.Errorf("not a directory: %s", o.ProjectDir)
	}

	template := resolveTemplate(o.ImageName, o.TemplateFallback, o.TemplateOverride)
	extra := appendResourceFlags(o.Extra, o.Resources)

	args := []string{"run", o.SbxAgent, "--template", template}
	for _, k := range o.KitPaths {
		args = append(args, "--kit", k)
	}
	args = append(args, extra...)
	args = append(args, absProject)

	fmt.Printf("==> sbx %s\n", strings.Join(args, " "))
	fmt.Printf("==> resources profile=%s memory=%s cpus=%s root=%s docker=%s\n",
		o.ResourcesProfile, o.Resources.Memory, o.Resources.CPUs, o.Resources.RootSize, o.Resources.DockerSize)

	cmd := exec.Command("sbx", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	if o.Resources.RootSize != "" {
		cmd.Env = setEnv(cmd.Env, "DOCKER_SANDBOXES_ROOT_SIZE", o.Resources.RootSize)
	}
	if o.Resources.DockerSize != "" {
		cmd.Env = setEnv(cmd.Env, "DOCKER_SANDBOXES_DOCKER_SIZE", o.Resources.DockerSize)
	}
	return cmd.Run()
}

func resolveTemplate(imageName, fallback, override string) string {
	if override != "" {
		return override
	}
	out, err := exec.Command("sbx", "template", "ls").Output()
	if err != nil {
		return fallback
	}
	best := ""
	for _, line := range strings.Split(string(out), "\n") {
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
