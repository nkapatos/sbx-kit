package initproj

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nkapatos/sbx-kit/cli/internal/catalog"
)

const (
	markerStart = "<!-- sbx-sandbox:start -->"
	markerEnd   = "<!-- sbx-sandbox:end -->"
)

// Opts stamps a Docker Sandbox section into a project README.
type Opts struct {
	Root       string
	Agent      string
	RecipeID   string
	ProjectDir string
	Catalog    *catalog.Catalog
}

// Run writes or updates the README section for the given catalog recipe.
func Run(o Opts) error {
	agentName := o.Agent
	if agentName == "" {
		agentName = "cursor"
	}
	agent, ok := o.Catalog.Agents[agentName]
	if !ok {
		return fmt.Errorf("unknown recipe %q (try: sbx-kit recipes)", agentName)
	}

	dir := o.ProjectDir
	if dir == "" {
		dir = "."
	}
	abs, err := filepath.Abs(dir)
	if err != nil {
		return err
	}
	st, err := os.Stat(abs)
	if err != nil || !st.IsDir() {
		return fmt.Errorf("not a directory: %s", dir)
	}

	display := o.RecipeID
	if display == "" {
		display = agentName
	}
	section := buildSection(display, agent, o.Catalog.Defaults.Kits)
	readme := filepath.Join(abs, "README.md")
	if _, err := os.Stat(readme); os.IsNotExist(err) {
		title := "# " + filepath.Base(abs) + "\n\n"
		if err := os.WriteFile(readme, []byte(title+section), 0o644); err != nil {
			return err
		}
		fmt.Printf("Created %s with Docker Sandbox section (recipe=%s)\n", readme, display)
		return nil
	}

	existing, err := os.ReadFile(readme)
	if err != nil {
		return err
	}
	body := string(existing)

	var out string
	if strings.Contains(body, markerStart) {
		out, err = replaceSection(body, section)
		if err != nil {
			return err
		}
		fmt.Printf("Updated Docker Sandbox section in %s (recipe=%s)\n", readme, display)
	} else {
		if !strings.HasSuffix(body, "\n") {
			body += "\n"
		}
		out = body + "\n" + section
		fmt.Printf("Appended Docker Sandbox section to %s (recipe=%s)\n", readme, display)
	}
	return os.WriteFile(readme, []byte(out), 0o644)
}

func buildSection(name string, a catalog.Agent, defaultKits []string) string {
	stubNote := ""
	if a.Stub {
		stubNote = " (kit still a stub in the catalog)"
	}
	image := a.ImageName
	if image == "" {
		image = a.SbxAgent
	}
	blurb := fmt.Sprintf("the `%s` recipe (`%s` + kits)%s", name, image, stubNote)

	runBlock := fmt.Sprintf("sbx-kit run %s --yes", name)
	if a.Stub {
		runBlock = fmt.Sprintf("# once the kit is ready:\n# sbx-kit run %s --yes", name)
	}

	miseBlock := ""
	for _, k := range catalog.ResolveKits(a.Kits, defaultKits) {
		if k == "mise-workspace" {
			miseBlock = `
1. Pin tools in ` + "`mise.toml`" + ` (or ` + "`.tool-versions`" + `) at the repo root.
2. From the project root:

` + "```bash" + `
` + runBlock + `
` + "```" + `

3. On a **new** sandbox, install pins once (agent follows kit ` + "`agentContext`" + `, or manually):

` + "```bash" + `
sbx exec -it <sandbox> -- bash -lc 'cd "$PWD" && mise trust mise.toml; mise install'
` + "```" + `

After removing pins: ` + "`mise install && mise prune -y`" + ` (prefer a fresh login shell if env looks stale).

The same sandbox keeps those tools across restarts.
`
			break
		}
	}
	if miseBlock == "" {
		miseBlock = `
From the project root:

` + "```bash" + `
` + runBlock + `
` + "```" + `

This recipe does **not** attach ` + "`mise-workspace`" + `. Use the image's
preinstalled tools (typical for official Hub kinds). The agent follows kit
` + "`agentContext`" + ` and ` + "`/etc/sbx-kit/floor.md`" + ` inside the box.
`
	}

	return strings.TrimSpace(fmt.Sprintf(`
%s
## Docker Sandbox

Run this project under [Docker AI Sandboxes](https://docs.docker.com/ai/sandboxes/) (`+"`sbx`"+`) with %s. That is **not** `+"`docker run`"+`.

**Host setup** (one-time): install [sbx-kit](https://github.com/nkapatos/sbx-kit) via Homebrew, ensure the Docker `+"`sbx`"+` CLI is signed in, and import a custom image if you are not using a stock Hub kind (`+"`sbx-kit image load`"+` or `+"`sbx-kit image pull`"+`).

### This project
%s
%s
`, markerStart, blurb, miseBlock, markerEnd)) + "\n"
}

func replaceSection(body, section string) (string, error) {
	start := strings.Index(body, markerStart)
	end := strings.Index(body, markerEnd)
	if start < 0 || end < 0 || end < start {
		return "", fmt.Errorf("malformed sbx-sandbox markers in README")
	}
	end += len(markerEnd)
	rest := body[end:]
	return body[:start] + section + rest, nil
}
