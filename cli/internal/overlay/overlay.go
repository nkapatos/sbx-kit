// Package overlay installs the sbx-kit platform contract into a sandbox.
//
// This is not an sbx kit and is not listed in the user's agents.yaml.
// User kits cannot omit or replace it. Values that belong to the CLI
// (command names, version) are generated from templates at install time.
// /etc/sbx-kit/context.md is the box-agent index (overlay analog of kit
// agentContext). User kit agentContext may point at it; it must not own it.
package overlay

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text/template"

	"github.com/nkapatos/sbx-kit/cli/internal/sbxcompat"
	"github.com/nkapatos/sbx-kit/cli/internal/stdio"
	"github.com/nkapatos/sbx-kit/cli/internal/version"
)

//go:embed files/*
var files embed.FS

const (
	// Helper is the in-VM pack/unpack binary name.
	Helper = "sbx-kit-state"
	// ShareDir is the overlay root inside the sandbox.
	ShareDir = "/home/agent/.local/share/sbx-kit"
	// ManifestEnv is the env var the helper reads.
	ManifestEnv = "SBX_KIT_STATE_MANIFEST"
)

// Runner is the sbx copy/exec surface used to install the overlay.
type Runner interface {
	Cp(src, dst string) error
	Exec(sandbox string, command ...string) error
}

// Install copies the platform overlay into an existing sandbox.
func Install(r Runner, sandbox string, w io.Writer) error {
	out := stdio.Out(w)
	fmt.Fprintf(out, "==> sbx-kit overlay -> %s\n", sandbox)

	tmp, err := os.MkdirTemp("", "sbx-kit-overlay-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	helper, err := files.ReadFile("files/sbx-kit-state")
	if err != nil {
		return err
	}
	manifest, err := files.ReadFile("files/state.manifest")
	if err != nil {
		return err
	}
	portable, err := files.ReadFile("files/portable-README.md")
	if err != nil {
		return err
	}
	floor, err := files.ReadFile("files/floor.sh")
	if err != nil {
		return err
	}
	contextBody, err := files.ReadFile("files/context.md")
	if err != nil {
		return err
	}
	cliBody, err := renderCLI()
	if err != nil {
		return err
	}

	hostHelper := filepath.Join(tmp, Helper)
	hostManifest := filepath.Join(tmp, "state.manifest")
	hostPortable := filepath.Join(tmp, "README.md")
	hostFloor := filepath.Join(tmp, "floor.sh")
	hostContext := filepath.Join(tmp, "context.md")
	hostCLI := filepath.Join(tmp, "cli.md")
	for _, item := range []struct {
		path string
		body []byte
		mode os.FileMode
	}{
		{hostHelper, helper, 0o755},
		{hostManifest, manifest, 0o644},
		{hostPortable, portable, 0o644},
		{hostFloor, floor, 0o755},
		{hostContext, contextBody, 0o644},
		{hostCLI, []byte(cliBody), 0o644},
	} {
		if err := os.WriteFile(item.path, item.body, item.mode); err != nil {
			return err
		}
	}

	share := ShareDir
	if err := r.Exec(sandbox, "bash", "-lc",
		fmt.Sprintf("mkdir -p %s/portable %s/bin", share, share)); err != nil {
		return fmt.Errorf("overlay mkdir: %w", err)
	}

	copies := []struct{ src, dst string }{
		{hostHelper, sandbox + ":" + share + "/bin/" + Helper},
		{hostManifest, sandbox + ":" + share + "/state.manifest"},
		{hostPortable, sandbox + ":" + share + "/portable/README.md"},
		{hostContext, sandbox + ":" + share + "/context.md"},
		{hostCLI, sandbox + ":" + share + "/cli.md"},
		{hostFloor, sandbox + ":" + share + "/bin/sbx-kit-floor.sh"},
	}
	for _, c := range copies {
		if err := r.Cp(c.src, c.dst); err != nil {
			return fmt.Errorf("overlay cp %s: %w", c.dst, err)
		}
	}

	finish := fmt.Sprintf(`set -euo pipefail
SHARE=%q
chmod 755 "$SHARE/bin/%s" "$SHARE/bin/sbx-kit-floor.sh"
sudo mkdir -p /etc/sbx-kit
sudo install -m 755 "$SHARE/bin/%s" /usr/local/bin/%s
sudo install -m 755 "$SHARE/bin/sbx-kit-floor.sh" /usr/local/bin/sbx-kit-floor.sh
sudo install -m 644 "$SHARE/cli.md" /etc/sbx-kit/cli.md
sudo install -m 644 "$SHARE/context.md" /etc/sbx-kit/context.md
bash /usr/local/bin/sbx-kit-floor.sh
sudo touch /etc/sandbox-persistent.sh
if ! grep -qF '# sbx-kit-overlay-state' /etc/sandbox-persistent.sh; then
  printf '%%s\n' \
    '# sbx-kit-overlay-state' \
    'export %s="${%s:-%s/state.manifest}"' \
    | sudo tee -a /etc/sandbox-persistent.sh >/dev/null
fi
`, share, Helper, Helper, Helper, ManifestEnv, ManifestEnv, share)

	if err := r.Exec(sandbox, "bash", "-lc", finish); err != nil {
		return fmt.Errorf("overlay install: %w", err)
	}
	return nil
}

func renderCLI() (string, error) {
	b, err := files.ReadFile("files/cli.md.tmpl")
	if err != nil {
		return "", err
	}
	t, err := template.New("cli.md").Parse(string(b))
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, struct {
		Version string
		MinSbx  string
	}{Version: version.Version, MinSbx: sbxcompat.MinVersion})
	return buf.String(), err
}
