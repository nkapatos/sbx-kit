package template

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/nkapatos/sbx-kit/cli/internal/sbxcompat"
	"github.com/nkapatos/sbx-kit/cli/internal/sbxutil"
	"github.com/nkapatos/sbx-kit/cli/internal/stdio"
)

// LoadOpts controls template import into sbx.
type LoadOpts struct {
	Root       string
	Engine     string // docker | container
	NameOrPath string
	ImageTag   string
	Out        io.Writer
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
	if b.BakeBase != "" {
		fmt.Fprintf(stdio.Out(o.Out), "==> shared bake: BASE_IMAGE=%s\n", b.BakeBase)
	}
	if IsParentOnly(b.TemplateDir) {
		return errParentNotImported(b.Name)
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

	if err := buildParentIfNeeded(engine, o.Root, b, o.Out); err != nil {
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
		if err := loadDocker(b, dockerTar, o.Out); err != nil {
			return err
		}
	case "container":
		ociTar := filepath.Join(tmpDir, "sbx-"+b.Name+".oci.tar")
		if err := loadContainer(b, ociTar, dockerTar, o.Out); err != nil {
			return err
		}
	}

	return importIntoSbx(b.ImageTag, dockerTar, o.Out)
}

func errParentNotImported(name string) error {
	return fmt.Errorf("%s is a parent image (Docker FROM base), not imported into sbx.\n  Load a child image under images/ instead (sbx-kit recipes image load).\n  The parent is docker-built automatically when you load the child.", name)
}

func buildParentIfNeeded(engine, root string, child *Build, w io.Writer) error {
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
	fmt.Fprintf(stdio.Out(w), "==> parent %s (Docker FROM; not imported into sbx)\n", pb.ImageTag)
	return buildImage(engine, pb, w)
}

func buildImage(engine string, b *Build, w io.Writer) error {
	switch engine {
	case "docker":
		return dockerBuild(b, w)
	case "container":
		return containerBuild(b, w)
	default:
		return fmt.Errorf("unknown engine %q", engine)
	}
}

// PullOpts controls registry pull + import into sbx.
type PullOpts struct {
	Engine   string // docker (default)
	ImageTag string
	Out      io.Writer
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
		return fmt.Errorf("image tag required (e.g. ghcr.io/org/example-image:latest)")
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

	out := stdio.Out(o.Out)
	fmt.Fprintf(out, "==> [1/2] docker pull %s\n", imageTag)
	if err := runLogged("docker", "pull", imageTag); err != nil {
		return err
	}

	fmt.Fprintf(out, "==> [2/2] docker image save -> %s\n", dockerTar)
	if err := runLogged("docker", "image", "save", imageTag, "-o", dockerTar); err != nil {
		return err
	}

	return importIntoSbx(imageTag, dockerTar, o.Out)
}

func importIntoSbx(imageTag, dockerTar string, w io.Writer) error {
	out := stdio.Out(w)
	fmt.Fprintln(out, "==> sbx template load")
	if err := runLogged("sbx", "template", "load", dockerTar); err != nil {
		return err
	}

	fmt.Fprintln(out, "==> Verifying sbx can see the template:")
	_ = runLogged("sbx", "template", "ls")

	fmt.Fprintln(out)
	fmt.Fprintln(out, "Done. Confirm with sbx template ls (engine store).")
	fmt.Fprintf(out, "Typical next step:\n  sbx-kit box run <dir>/<recipe> --yes\n")
	fmt.Fprintf(out, "(image tag: %s)\n", imageTag)
	return nil
}

func dockerBuild(b *Build, w io.Writer) error {
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
	fmt.Fprintf(stdio.Out(w), "==> docker build -t %s\n", b.ImageTag)
	return runLogged("docker", args...)
}

func containerBuild(b *Build, w io.Writer) error {
	if _, err := exec.LookPath("container"); err != nil {
		return fmt.Errorf("Apple container CLI not found on PATH")
	}
	args := []string{"build", "-t", b.ImageTag, "-f", b.Dockerfile}
	for _, a := range b.BuildArgs {
		args = append(args, "--build-arg", a)
	}
	args = append(args, b.Context)
	fmt.Fprintf(stdio.Out(w), "==> container build -t %s\n", b.ImageTag)
	return runLogged("container", args...)
}

func loadDocker(b *Build, dockerTar string, w io.Writer) error {
	if err := dockerBuild(b, w); err != nil {
		return err
	}
	fmt.Fprintf(stdio.Out(w), "==> docker image save -> %s\n", dockerTar)
	return runLogged("docker", "image", "save", b.ImageTag, "-o", dockerTar)
}

func loadContainer(b *Build, ociTar, dockerTar string, w io.Writer) error {
	if _, err := exec.LookPath("skopeo"); err != nil {
		return fmt.Errorf("skopeo not found on PATH (brew install skopeo; needed for OCI → docker-archive)")
	}
	arch, err := hostArch()
	if err != nil {
		return err
	}
	if err := containerBuild(b, w); err != nil {
		return err
	}

	fmt.Fprintf(stdio.Out(w), "==> container image save (OCI) -> %s\n", ociTar)
	if err := runLogged("container", "image", "save", b.ImageTag, "-o", ociTar); err != nil {
		return err
	}

	fmt.Fprintf(stdio.Out(w), "==> skopeo: OCI -> docker-archive -> %s\n", dockerTar)
	return runLogged("skopeo", "copy",
		"--override-os", "linux",
		"--override-arch", arch,
		"oci-archive:"+ociTar,
		"docker-archive:"+dockerTar+":"+b.ImageTag,
	)
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
