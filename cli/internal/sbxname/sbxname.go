package sbxname

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

var (
	validName = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.-]{0,62}$`)
)

// NewProfileID builds a stable opaque vault id for recipe + absolute project path.
// Used under ~/.local/share/sbx-kit/profiles/<id>/ — not as the sbx --name.
func NewProfileID(recipe, absProject string) string {
	sum := sha256.Sum256([]byte(absProject))
	short := hex.EncodeToString(sum[:])[:8]
	recipe = sanitizeToken(recipe)
	if recipe == "" {
		recipe = "agent"
	}
	return truncate(fmt.Sprintf("sbxk-%s-%s", recipe, short), 63)
}

// FromProject is an alias for NewProfileID (legacy name).
func FromProject(agent, absProject string) string {
	return NewProfileID(agent, absProject)
}

// FromDir returns a friendly sbx sandbox name from the project directory basename.
func FromDir(absProject string) (string, error) {
	base := filepath.Base(absProject)
	if base == "" || base == "." || base == string(filepath.Separator) {
		return "", fmt.Errorf("cannot derive sandbox name from path %q", absProject)
	}
	name := Sanitize(base)
	if !Valid(name) {
		return "", fmt.Errorf("derived sandbox name %q is invalid for sbx (from %q); pass --sandbox-name", name, base)
	}
	return name, nil
}

// Sanitize maps an arbitrary string to an sbx-safe name (keeps letters, digits, _ . -).
func Sanitize(s string) string {
	s = strings.TrimSpace(s)
	var b strings.Builder
	lastDash := false
	for _, r := range s {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			b.WriteRune(r)
			lastDash = false
		case r == '_' || r == '.' || r == '-':
			b.WriteRune(r)
			lastDash = false
		default:
			if !lastDash && b.Len() > 0 {
				b.WriteByte('-')
				lastDash = true
			}
		}
	}
	out := strings.Trim(b.String(), "-._")
	return truncate(out, 63)
}

// ExtractFromArgs returns --name value from sbx passthrough args, if present.
func ExtractFromArgs(args []string) (name string, rest []string, ok bool) {
	rest = make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "--name" && i+1 < len(args):
			return args[i+1], append(rest, args[i+2:]...), true
		case strings.HasPrefix(a, "--name="):
			return strings.TrimPrefix(a, "--name="), append(rest, args[i+1:]...), true
		default:
			rest = append(rest, a)
		}
	}
	return "", args, false
}

// Inject ensures --name is present; returns the name and full args.
func Inject(args []string, name string) (string, []string) {
	if existing, rest, ok := ExtractFromArgs(args); ok {
		return existing, append([]string{"--name", existing}, rest...)
	}
	return name, append([]string{"--name", name}, args...)
}

// Valid reports whether name is acceptable for sbx.
func Valid(name string) bool {
	return name != "" && validName.MatchString(name)
}

func sanitizeToken(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, "/", "-")
	var b strings.Builder
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == '-' || r == '_':
			b.WriteRune(r)
		}
	}
	return b.String()
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}
