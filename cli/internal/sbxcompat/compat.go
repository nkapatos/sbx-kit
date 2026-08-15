// Package sbxcompat gates sbx-kit against a tested Docker sbx CLI range.
// Kits/templates remain experimental upstream; lock the floor so behavior
// changes fail fast instead of with cryptic YAML errors.
package sbxcompat

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

// MinVersion is the oldest sbx CLI this toolkit tree is tested against.
// Bump when we rely on newer sbx behavior.
// Kits stay schemaVersion "1" until MinVersion is ≥ 0.36 and every spec.yaml
// is rewritten to the v2 grammar (v1 still loads via sbx's legacy path).
const MinVersion = "0.34.0"

// MaxVersion, if non-empty, is an exclusive upper bound (sbx < MaxVersion).
// Leave empty while tracking latest; set when a known break lands upstream.
const MaxVersion = ""

const skipEnv = "SBX_KIT_SKIP_SBX_CHECK"

var (
	versionToken = regexp.MustCompile(`(?i)\bv?(\d+\.\d+\.\d+(?:[-+][0-9A-Za-z.-]+)?)\b`)

	mu       sync.Mutex
	cachedOK bool
	cached   string
	cachedErr error
)

// Reset clears the cached check (tests).
func Reset() {
	mu.Lock()
	defer mu.Unlock()
	cachedOK = false
	cached = ""
	cachedErr = nil
}

// Check parses installed version text from `sbx version` / `sbx --version`
// output and ensures it falls in [MinVersion, MaxVersion).
func Check(versionOutput string) error {
	if os.Getenv(skipEnv) != "" {
		return nil
	}
	ver, err := ParseVersion(versionOutput)
	if err != nil {
		return fmt.Errorf("cannot parse sbx version from %q: %w\n  tip: install Docker sbx >= %s, or set %s=1 to skip",
			strings.TrimSpace(versionOutput), err, MinVersion, skipEnv)
	}
	if cmpSemver(ver, MinVersion) < 0 {
		return fmt.Errorf("sbx %s is too old for this sbx-kit (need >= %s)\n  upgrade: brew upgrade docker/tap/sbx   # or your distro package\n  override: %s=1",
			ver, MinVersion, skipEnv)
	}
	if MaxVersion != "" && cmpSemver(ver, MaxVersion) >= 0 {
		return fmt.Errorf("sbx %s is newer than this sbx-kit supports (need < %s)\n  upgrade sbx-kit, or set %s=1 to skip",
			ver, MaxVersion, skipEnv)
	}
	return nil
}

// Ensure runs Check once per process against output from getVersion.
func Ensure(getVersion func() (string, error)) error {
	mu.Lock()
	defer mu.Unlock()
	if cachedOK {
		return cachedErr
	}
	out, err := getVersion()
	if err != nil {
		cachedOK = true
		cachedErr = fmt.Errorf("sbx not usable (%v)\n  install Docker sbx >= %s; override: %s=1", err, MinVersion, skipEnv)
		return cachedErr
	}
	cached = strings.TrimSpace(out)
	cachedErr = Check(cached)
	cachedOK = true
	return cachedErr
}

// LastVersion returns the last successfully parsed version string (may be empty).
func LastVersion() string {
	mu.Lock()
	defer mu.Unlock()
	if v, err := ParseVersion(cached); err == nil {
		return v
	}
	return cached
}

// ParseVersion extracts the first semver-like token from sbx version output.
func ParseVersion(s string) (string, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", fmt.Errorf("empty version output")
	}
	m := versionToken.FindStringSubmatch(s)
	if m == nil {
		return "", fmt.Errorf("no semver token found")
	}
	// Strip pre-release/build for comparison floor (0.34.0-rc3 >= 0.34.0 for our purposes? 
	// Actually rc should be < release. Keep core for cmp; store full core.major.minor.patch)
	core := m[1]
	if i := strings.IndexAny(core, "-+"); i >= 0 {
		core = core[:i]
	}
	parts := strings.Split(core, ".")
	if len(parts) != 3 {
		return "", fmt.Errorf("invalid semver %q", m[1])
	}
	return core, nil
}

func cmpSemver(a, b string) int {
	ap := mustParts(a)
	bp := mustParts(b)
	for i := 0; i < 3; i++ {
		if ap[i] < bp[i] {
			return -1
		}
		if ap[i] > bp[i] {
			return 1
		}
	}
	return 0
}

func mustParts(v string) [3]int {
	var out [3]int
	parts := strings.Split(v, ".")
	for i := 0; i < 3 && i < len(parts); i++ {
		n, _ := strconv.Atoi(parts[i])
		out[i] = n
	}
	return out
}

// RequirementSummary is shown by `sbx-kit version`.
func RequirementSummary() string {
	s := fmt.Sprintf("requires sbx >= %s (kits authored as schemaVersion 1)", MinVersion)
	if MaxVersion != "" {
		s = fmt.Sprintf("requires sbx >= %s and < %s", MinVersion, MaxVersion)
	}
	return s
}
