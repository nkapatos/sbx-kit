package sbxname

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
)

var (
	validName = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.-]{0,62}$`)
)

// FromProject builds a stable sbx --name for agent + absolute project path.
func FromProject(agent, absProject string) string {
	sum := sha256.Sum256([]byte(absProject))
	short := hex.EncodeToString(sum[:])[:8]
	agent = sanitizeToken(agent)
	if agent == "" {
		agent = "agent"
	}
	name := fmt.Sprintf("sbxk-%s-%s", agent, short)
	return truncate(name, 63)
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
