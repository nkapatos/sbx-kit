package statexfer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/nkapatos/sbx-kit/cli/internal/sbxutil"
	"github.com/nkapatos/sbx-kit/cli/internal/xdg"
)

const (
	// RemoteArchive is the in-VM pack/unpack path.
	RemoteArchive = "/tmp/sbx-kit-state.tgz"
	// Helper is installed by the agent-workspace kit.
	Helper = "sbx-kit-state"
	// DefaultStopWait is a short grace period when export sees status=running.
	DefaultStopWait = 15 * time.Second
)

// Export packs VM portable state and copies it to the host profile archive.
// If the sandbox is still "running", waits for it to stop so agent SQLite WALs
// can checkpoint cleanly inside sbx-kit-state pack.
func Export(r *sbxutil.Runner, sandbox, profileID string) error {
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

	if s, err := r.Get(sandbox); err == nil && s != nil && strings.EqualFold(s.Status, "running") {
		fmt.Printf("==> sandbox %s is running; waiting briefly for detach before packing (best-effort)\n", sandbox)
		if err := r.WaitNotRunning(sandbox, 15*time.Second); err != nil {
			fmt.Printf("==> warning: %v\n==> packing anyway; prefer detach so SQLite WALs can flush\n", err)
		}
	}

	fmt.Printf("==> packing state in %s via %s\n", sandbox, Helper)
	if err := r.ExecVisible(sandbox, Helper, "pack", RemoteArchive); err != nil {
		return fmt.Errorf("pack failed (is agent-workspace kit installed?): %w", err)
	}

	fmt.Printf("==> copying %s:%s -> %s\n", sandbox, RemoteArchive, hostArch)
	if err := r.Cp(sandbox+":"+RemoteArchive, hostArch); err != nil {
		return err
	}
	fmt.Printf("==> state saved: %s\n", hostArch)
	return nil
}

// Import copies the host profile archive into the VM and unpacks it.
func Import(r *sbxutil.Runner, sandbox, profileID string) error {
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

	fmt.Printf("==> copying %s -> %s:%s\n", hostArch, sandbox, RemoteArchive)
	if err := r.Cp(hostArch, sandbox+":"+RemoteArchive); err != nil {
		return err
	}
	fmt.Printf("==> unpacking state in %s via %s\n", sandbox, Helper)
	if err := r.ExecVisible(sandbox, Helper, "unpack", RemoteArchive); err != nil {
		return fmt.Errorf("unpack failed (is agent-workspace kit installed?): %w", err)
	}
	fmt.Printf("==> state restored into %s\n", sandbox)
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
