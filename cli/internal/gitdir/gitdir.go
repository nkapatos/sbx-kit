package gitdir

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// DirName is the checkout directory name implied by a git remote URL.
func DirName(url string) string {
	s := strings.TrimSpace(url)
	s = strings.TrimSuffix(s, "/")
	s = strings.TrimSuffix(s, ".git")
	switch {
	case strings.HasPrefix(s, "https://"):
		s = strings.TrimPrefix(s, "https://")
	case strings.HasPrefix(s, "http://"):
		s = strings.TrimPrefix(s, "http://")
	case strings.HasPrefix(s, "ssh://"):
		s = strings.TrimPrefix(s, "ssh://")
	default:
		if i := strings.Index(s, "@"); i >= 0 {
			s = s[i+1:]
			s = strings.Replace(s, ":", "/", 1)
		}
	}
	return filepath.Base(s)
}

func IsRepo(dir string) bool {
	_, err := os.Stat(filepath.Join(dir, ".git"))
	return err == nil
}

func Clone(url, dest string) error {
	return git("clone", "--", url, dest)
}

func Fetch(dir string) error {
	return git("-C", dir, "fetch", "-q")
}

func Pull(dir string) error {
	return git("-C", dir, "pull", "--ff-only")
}

func Status(dir string) (string, error) {
	if err := lookGit(); err != nil {
		return "", err
	}
	out, err := exec.Command("git", "-C", dir, "status", "-sb").Output()
	if err != nil {
		return "", fmt.Errorf("git status: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

func RemoteURL(dir string) string {
	out, err := exec.Command("git", "-C", dir, "remote", "get-url", "origin").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func lookGit() error {
	if _, err := exec.LookPath("git"); err != nil {
		return fmt.Errorf("git not found on PATH")
	}
	return nil
}

func git(args ...string) error {
	if err := lookGit(); err != nil {
		return err
	}
	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	verb := "git"
	if len(args) > 0 {
		verb = args[0]
		if verb == "-C" && len(args) > 2 {
			verb = args[2]
		}
	}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git %s: %w", verb, err)
	}
	return nil
}
