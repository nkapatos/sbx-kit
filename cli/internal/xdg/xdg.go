package xdg

import (
	"os"
	"path/filepath"
)

const appName = "sbx-kit"

// ShareDir is ~/.local/share/sbx-kit (or $XDG_DATA_HOME/sbx-kit).
func ShareDir() (string, error) {
	if v := os.Getenv("XDG_DATA_HOME"); v != "" {
		return filepath.Join(v, appName), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".local", "share", appName), nil
}

// StateDir is ~/.local/state/sbx-kit (or $XDG_STATE_HOME/sbx-kit).
func StateDir() (string, error) {
	if v := os.Getenv("XDG_STATE_HOME"); v != "" {
		return filepath.Join(v, appName), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".local", "state", appName), nil
}

// ProfilesDir is share/profiles.
func ProfilesDir() (string, error) {
	share, err := ShareDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(share, "profiles"), nil
}

// ProfileArchive is share/profiles/<id>/state.tgz.
func ProfileArchive(profileID string) (string, error) {
	dir, err := ProfilesDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, profileID, "state.tgz"), nil
}

// Ensure creates share/profiles and state directories.
func Ensure() error {
	share, err := ShareDir()
	if err != nil {
		return err
	}
	state, err := StateDir()
	if err != nil {
		return err
	}
	profiles, err := ProfilesDir()
	if err != nil {
		return err
	}
	for _, d := range []string{share, state, profiles} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			return err
		}
	}
	return nil
}
