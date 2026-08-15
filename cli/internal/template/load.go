package template

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/nkapatos/sbx-kit/cli/internal/sbxcompat"
	"github.com/nkapatos/sbx-kit/cli/internal/sbxutil"
)

// LoadOpts controls template import into sbx.
type LoadOpts struct {
	Root       string
	Engine     string // docker | container
	NameOrPath string
	ImageTag   string
}

// Load builds a template image and imports it via `sbx template load`.
func Load(o LoadOpts) error {
	engine := strings.ToLower(strings.TrimSpace(o.Engine))
	switch engine {
	case "docker", "container":
	default:
		return fmt.Errorf("unknown engine %q (use docker or container)", o.Engine)
	}

	b, err := ResolveBuild(o.Root, o.NameOrPath, o.ImageTag)
	if err != nil {
		return err
	}
	if IsParentOnly(b.TemplateDir) {
		return errParentNotImported(engine, b.Name)
	}

	r := sbxutil.Default()
	if _, err := r.LookPath(); err != nil {
		return fmt.Errorf("sbx not found on PATH (run this on the host; need >= %s)", sbxcompat.MinVersion)
	}
	if err := sbxcompat.Ensure(func() (string, error) {
		return r.ProbeVersion()
	}); err != nil {
		return err
	}

	if err := buildParentIfNeeded(engine, o.Root, b); err != nil {
		return err
	}

	tmpDir, err := os.MkdirTemp("", "sbx-kit-template-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	dockerTar := filepath.Join(tmpDir, "sbx-"+b.Name+".docker.tar")

	switch engine {
	case "docker":
		if err := loadDocker(b, dockerTar); err != nil {
			return err
		}
	case "container":
		ociTar := filepath.Join(tmpDir, "sbx-"+b.Name+".oci.tar")
		if err := loadContainer(b, ociTar, dockerTar); err != nil {
			return err
		}
	}

	return importIntoSbx(b.ImageTag, dockerTar)
}

func errParentNotImported(engine, name string) error {
	return fmt.Errorf("%s is a parent image (Docker FROM base), not imported into sbx.\n  Minimum sandbox image: sbx-kit image load --engine %s kit-shell\n  Baked agent example:   sbx-kit image load --engine %s kit-cursor\nThe parent is docker-built automatically when you load those.", name, engine, engine)
}

func buildParentIfNeeded(engine, root string, child *Build) error {
	parent, err := ParentTemplateName(child.Dockerfile)
	if err != nil {
		return err
	}
	if parent == "" || parent == child.Name {
		return nil
	}
	pb, err := ResolveBuild(root, parent, "")
	if err != nil {
		return fmt.Errorf("parent image %s: %w", parent, err)
	}
	fmt.Printf("==> parent %s (Docker FROM; not imported into sbx)\n", pb.ImageTag)
	return buildImage(engine, pb)
}

func buildImage(engine string, b *Build) error {
	switch engine {
	case "docker":
		return dockerBuild(b)
	case "container":
		return containerBuild(b)
	default:
		return fmt.Errorf("unknown engine %q", engine)
	}
}

// PullOpts controls registry pull + import into sbx.
type PullOpts struct {
	Engine   string // docker (default)
	ImageTag string
}

// Pull fetches a registry image with docker and imports it via `sbx template load`.
func Pull(o PullOpts) error {
	engine := strings.ToLower(strings.TrimSpace(o.Engine))
	if engine == "" {
		engine = "docker"
	}
	if engine != "docker" {
		return fmt.Errorf("image pull currently supports --engine docker only")
	}
	imageTag := strings.TrimSpace(o.ImageTag)
	if imageTag == "" {
		return fmt.Errorf("image tag required (e.g. ghcr.io/org/sbx-kit-cursor:latest)")
	}

	r := sbxutil.Default()
	if _, err := r.LookPath(); err != nil {
		return fmt.Errorf("sbx not found on PATH (run this on the host; need >= %s)", sbxcompat.MinVersion)
	}
	if err := sbxcompat.Ensure(func() (string, error) {
		return r.ProbeVersion()
	}); err != nil {
		return err
	}

	if _, err := exec.LookPath("docker"); err != nil {
		return fmt.Errorf("docker CLI not found on PATH (Docker Desktop / Colima)")
	}
	if err := runLogged("docker", "info"); err != nil {
		return fmt.Errorf("docker daemon not reachable (is Docker Desktop / Colima running?): %w", err)
	}

	tmpDir, err := os.MkdirTemp("", "sbx-kit-pull-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	safe := strings.NewReplacer("/", "_", ":", "_").Replace(imageTag)
	dockerTar := filepath.Join(tmpDir, "sbx-"+safe+".docker.tar")

	fmt.Printf("==> [1/2] docker pull %s\n", imageTag)
	if err := runLogged("docker", "pull", imageTag); err != nil {
		return err
	}

	fmt.Printf("==> [2/2] docker image save -> %s\n", dockerTar)
	if err := runLogged("docker", "image", "save", imageTag, "-o", dockerTar); err != nil {
		return err
	}

	return importIntoSbx(imageTag, dockerTar)
}

func importIntoSbx(imageTag, dockerTar string) error {
	fmt.Println("==> sbx template load")
	if err := runLogged("sbx", "template", "load", dockerTar); err != nil {
		return err
	}

	fmt.Println("==> Verifying sbx can see the template:")
	_ = runLogged("sbx", "template", "ls")

	fmt.Println()
	fmt.Println("Done. Confirm with sbx template ls (engine store).")
	fmt.Printf("Typical next step:\n  sbx-kit run kit-shell --yes   # or: kit-cursor / kit-pi\n")
	fmt.Printf("(image tag: %s)\n", imageTag)
	return nil
}

func dockerBuild(b *Build) error {
	if _, err := exec.LookPath("docker"); err != nil {
		return fmt.Errorf("docker CLI not found on PATH (Docker Desktop / Colima, or use --engine container)")
	}
	if err := runLogged("docker", "info"); err != nil {
		return fmt.Errorf("docker daemon not reachable (is Docker Desktop / Colima running?): %w", err)
	}
	args := []string{"build", "-t", b.ImageTag, "-f", b.Dockerfile}
	for _, a := range b.BuildArgs {
		args = append(args, "--build-arg", a)
	}
	args = append(args, b.Context)
	fmt.Printf("==> docker build -t %s\n", b.ImageTag)
	return runLogged("docker", args...)
}

func containerBuild(b *Build) error {
	if _, err := exec.LookPath("container"); err != nil {
		return fmt.Errorf("Apple container CLI not found on PATH")
	}
	args := []string{"build", "-t", b.ImageTag, "-f", b.Dockerfile}
	for _, a := range b.BuildArgs {
		args = append(args, "--build-arg", a)
	}
	args = append(args, b.Context)
	fmt.Printf("==> container build -t %s\n", b.ImageTag)
	return runLogged("container", args...)
}

func loadDocker(b *Build, dockerTar string) error {
	if err := dockerBuild(b); err != nil {
		return err
	}
	fmt.Printf("==> docker image save -> %s\n", dockerTar)
	if err := runLogged("docker", "image", "save", b.ImageTag, "-o", dockerTar); err != nil {
		return err
	}
	return smokeAgentBinary(b)
}

func loadContainer(b *Build, ociTar, dockerTar string) error {
	if _, err := exec.LookPath("skopeo"); err != nil {
		return fmt.Errorf("skopeo not found on PATH (brew install skopeo; needed for OCI → docker-archive)")
	}
	arch, err := hostArch()
	if err != nil {
		return err
	}
	if err := containerBuild(b); err != nil {
		return err
	}

	// Smoke runs via docker CLI when available; Apple container path skips host probe.
	_ = smokeAgentBinary(b)

	fmt.Printf("==> container image save (OCI) -> %s\n", ociTar)
	if err := runLogged("container", "image", "save", b.ImageTag, "-o", ociTar); err != nil {
		return err
	}

	fmt.Printf("==> skopeo: OCI -> docker-archive -> %s\n", dockerTar)
	return runLogged("skopeo", "copy",
		"--override-os", "linux",
		"--override-arch", arch,
		"oci-archive:"+ociTar,
		"docker-archive:"+dockerTar+":"+b.ImageTag,
	)
}

// smokeAgentBinary verifies layered agent images expose their CLI on PATH before
// import. Catches "agent binary not found" failures early.
func smokeAgentBinary(b *Build) error {
	bin := ""
	switch b.Name {
	case "kit-cursor":
		bin = "cursor-agent"
	default:
		return nil
	}
	if _, err := exec.LookPath("docker"); err != nil {
		fmt.Printf("==> skip smoke (%s): docker CLI not available\n", bin)
		return nil
	}
	fmt.Printf("==> smoke: docker run --rm --entrypoint which %s %s\n", b.ImageTag, bin)
	if err := runLogged("docker", "run", "--rm", "--entrypoint", "which", b.ImageTag, bin); err != nil {
		return fmt.Errorf("image %s is missing %q on PATH (rebuild parent kit-core then this template): %w", b.ImageTag, bin, err)
	}
	return nil
}

func hostArch() (string, error) {
	switch runtime.GOARCH {
	case "arm64":
		return "arm64", nil
	case "amd64":
		return "amd64", nil
	default:
		return "", fmt.Errorf("unsupported machine arch: %s", runtime.GOARCH)
	}
}

func runLogged(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
