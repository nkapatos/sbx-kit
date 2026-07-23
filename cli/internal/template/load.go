package template

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// LoadOpts controls template import into sbx.
type LoadOpts struct {
	Root     string
	Engine   string // docker | container
	NameOrPath string
	ImageTag string
}

// Load builds a template image and imports it via `sbx template load`.
func Load(o LoadOpts) error {
	engine := strings.ToLower(strings.TrimSpace(o.Engine))
	switch engine {
	case "docker", "container":
	default:
		return fmt.Errorf("unknown engine %q (use docker or container)", o.Engine)
	}

	if _, err := exec.LookPath("sbx"); err != nil {
		return fmt.Errorf("sbx not found on PATH (run this on the host)")
	}

	b, err := ResolveBuild(o.Root, o.NameOrPath, o.ImageTag)
	if err != nil {
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

	fmt.Println("==> sbx template load")
	if err := runLogged("sbx", "template", "load", dockerTar); err != nil {
		return err
	}

	fmt.Println("==> Verifying sbx can see the template:")
	_ = runLogged("sbx", "template", "ls")

	fmt.Println()
	fmt.Println("Done. Use the exact REPOSITORY:TAG from sbx template ls with --template if needed.")
	fmt.Printf("Typical next step:\n  sbx-kit run cursor .\n")
	fmt.Printf("(image tag used for build: %s)\n", b.ImageTag)
	return nil
}

func loadDocker(b *Build, dockerTar string) error {
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

	fmt.Printf("==> [1/3] docker build -t %s\n", b.ImageTag)
	if err := runLogged("docker", args...); err != nil {
		return err
	}

	fmt.Printf("==> [2/3] docker image save -> %s\n", dockerTar)
	return runLogged("docker", "image", "save", b.ImageTag, "-o", dockerTar)
}

func loadContainer(b *Build, ociTar, dockerTar string) error {
	if _, err := exec.LookPath("container"); err != nil {
		return fmt.Errorf("Apple container CLI not found on PATH")
	}
	if _, err := exec.LookPath("skopeo"); err != nil {
		return fmt.Errorf("skopeo not found on PATH (brew install skopeo; needed for OCI → docker-archive)")
	}

	arch, err := hostArch()
	if err != nil {
		return err
	}

	args := []string{"build", "-t", b.ImageTag, "-f", b.Dockerfile}
	for _, a := range b.BuildArgs {
		args = append(args, "--build-arg", a)
	}
	args = append(args, b.Context)

	fmt.Printf("==> [1/4] container build -t %s\n", b.ImageTag)
	if err := runLogged("container", args...); err != nil {
		return err
	}

	fmt.Printf("==> [2/4] container image save (OCI) -> %s\n", ociTar)
	if err := runLogged("container", "image", "save", b.ImageTag, "-o", ociTar); err != nil {
		return err
	}

	fmt.Printf("==> [3/4] skopeo: OCI -> docker-archive -> %s\n", dockerTar)
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
