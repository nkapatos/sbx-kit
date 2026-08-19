// Package statexfer copies portable workplace state between a running
// sandbox and the host profile archive.
//
// Pack/unpack is performed in-VM by sbx-kit-state, installed by the CLI
// overlay (not a user kit). The default manifest includes only the portable
// share; extra INCLUDE lines are overlay/user concern. SQLite WAL checkpoint
// is best-effort when .db files exist under INCLUDE trees.
//
// Parked in-box prompts for assisted pack/handoff: experimental prompts
// (cli/internal/boxprompt); share that set between boxes later.
package statexfer

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/nkapatos/sbx-kit/cli/internal/sbxutil"
	"github.com/nkapatos/sbx-kit/cli/internal/stdio"
	"github.com/nkapatos/sbx-kit/cli/internal/xdg"
)

const (
	// RemoteArchive is the in-VM pack/unpack path.
	RemoteArchive = "/tmp/sbx-kit-state.tgz"
	// Helper is installed by the CLI overlay.
	Helper = "sbx-kit-state"
	// DefaultStopWait is a short grace period when export sees status=running.
	// Useful when an agent has SQLite WALs; a no-op for agents that do not.
	DefaultStopWait = 15 * time.Second
)

// Export packs VM portable state and copies it to the host profile archive.
// If the sandbox is still "running", waits briefly for detach (best-effort)
// so kits that checkpoint SQLite can flush; kits without a DB still pack.
func Export(r *sbxutil.Runner, sandbox, profileID string, w io.Writer) error {
	if err := xdg.Ensure(); err != nil {
		return err
	}
	hostArch, err := xdg.ProfileArchive(profileID)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(hostArch), 0o755); err != nil {
		return err
	}

	out := stdio.Out(w)
	if s, err := r.Get(sandbox); err == nil && s != nil && strings.EqualFold(s.Status, "running") {
		fmt.Fprintf(out, "==> sandbox %s is running; waiting briefly for detach before packing (best-effort)\n", sandbox)
		if err := r.WaitNotRunning(sandbox, DefaultStopWait); err != nil {
			fmt.Fprintf(out, "==> warning: %v\n==> packing anyway; detach first if this agent uses SQLite\n", err)
		}
	}

	fmt.Fprintf(out, "==> packing state in %s via %s\n", sandbox, Helper)
	if err := r.Exec(sandbox, Helper, "pack", RemoteArchive); err != nil {
		return fmt.Errorf("pack failed (is the sbx-kit overlay installed?): %w", err)
	}

	fmt.Fprintf(out, "==> copying %s:%s -> %s\n", sandbox, RemoteArchive, hostArch)
	if err := r.Cp(sandbox+":"+RemoteArchive, hostArch); err != nil {
		return err
	}
	fmt.Fprintf(out, "==> state saved: %s\n", hostArch)
	return nil
}

// Import copies the host profile archive into the VM and unpacks it.
func Import(r *sbxutil.Runner, sandbox, profileID string, w io.Writer) error {
	hostArch, err := xdg.ProfileArchive(profileID)
	if err != nil {
		return err
	}
	st, err := os.Stat(hostArch)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("no saved state at %s", hostArch)
		}
		return err
	}
	if st.Size() == 0 {
		return fmt.Errorf("empty state archive: %s", hostArch)
	}

	out := stdio.Out(w)
	fmt.Fprintf(out, "==> copying %s -> %s:%s\n", hostArch, sandbox, RemoteArchive)
	if err := r.Cp(hostArch, sandbox+":"+RemoteArchive); err != nil {
		return err
	}
	fmt.Fprintf(out, "==> unpacking state in %s via %s\n", sandbox, Helper)
	if err := r.Exec(sandbox, Helper, "unpack", RemoteArchive); err != nil {
		return fmt.Errorf("unpack failed (is the sbx-kit overlay installed?): %w", err)
	}
	fmt.Fprintf(out, "==> state restored into %s\n", sandbox)
	return nil
}

// HasArchive reports whether a host profile archive exists.
func HasArchive(profileID string) (bool, error) {
	hostArch, err := xdg.ProfileArchive(profileID)
	if err != nil {
		return false, err
	}
	st, err := os.Stat(hostArch)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return st.Size() > 0, nil
}

// DiscardArchive removes the host profile archive (and empty parent dir) if present.
func DiscardArchive(profileID string, w io.Writer) error {
	hostArch, err := xdg.ProfileArchive(profileID)
	if err != nil {
		return err
	}
	if err := os.Remove(hostArch); err != nil && !os.IsNotExist(err) {
		return err
	}
	_ = os.Remove(filepath.Dir(hostArch)) // best-effort remove profiles/<id>/
	fmt.Fprintf(stdio.Out(w), "==> discarded archive for profile %s\n", profileID)
	return nil
}
