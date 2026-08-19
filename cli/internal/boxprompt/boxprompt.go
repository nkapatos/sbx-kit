// Package boxprompt parks one-shot prompts for the in-box agent.
//
// Run these inside the sandbox (not on the host CLI) for assisted
// maintenance, cleanup, state, and summaries. The CLI overlay and
// sbx-kit-state pack/unpack own what is portable; these prompts stay
// in tandem with that overlay.
//
// TODO: share the same prompt set between boxes (inject via the overlay,
// copy with statexfer). Do not fork per-agent copies.
package boxprompt

import (
	"fmt"
	"strings"
)

// Status is the parked-feature blurb (experimental prompts).
const Status = `In-box one-shot prompts: parked.

Paste into the agent attached to a box (sbx exec / agent session).
Later: same texts shared across boxes via the CLI overlay + statexfer.

Names: maintain, cleanup, state, summary`

// Prompt is one named one-shot for the in-box agent.
type Prompt struct {
	Name string
	Body string
}

// All is the parked set, in display order.
var All = []Prompt{
	{Name: "maintain", Body: maintain},
	{Name: "cleanup", Body: cleanup},
	{Name: "state", Body: state},
	{Name: "summary", Body: summary},
}

// Lookup returns a parked prompt by name (case-insensitive).
func Lookup(name string) (Prompt, error) {
	want := strings.ToLower(strings.TrimSpace(name))
	for _, p := range All {
		if p.Name == want {
			return p, nil
		}
	}
	return Prompt{}, fmt.Errorf("unknown prompt %q (try: maintain, cleanup, state, summary)", name)
}

const maintain = `One shot — box maintenance.

You are inside a Docker sbx sandbox. Do not change the host. Work under $PWD and /home/agent.

1. Read overlay index first: /etc/sbx-kit/context.md, then floor.md and cli.md. User kit agentContext is extra.
2. Report OS, disk, and whether mise / the agent CLI is on PATH.
3. If a mise kit is in play and mise.toml (or .tool-versions) exists at the project root, run mise trust + mise install only if pins are missing.
4. Note anything broken (missing binaries, failed installs, network denials) with the exact command that failed.
5. Do not recreate the sandbox. Do not write secrets into the VM.

Reply with a short checklist of what you checked and what still needs a human.`

const cleanup = `One shot — box cleanup.

You are inside a Docker sbx sandbox. Do not touch the host profile archive.

1. Clear safe caches only (package manager caches, /tmp junk you created). Do not delete /home/agent workplace state the next agent will need.
2. Leave agent history, editor config, and overlay files alone unless they are clearly trash.
3. If an agent uses SQLite under /home/agent, do not delete the DB; vacuum only if the tool supports it and the agent is idle.
4. If there is no SQLite DB, skip DB steps — default pack layout is the portable overlay dir.
5. Summarize bytes freed and paths you changed.

Stop after one pass. Do not install new software during cleanup.`

const state = `One shot — portable state.

You are inside a Docker sbx sandbox. Default portable state is /home/agent/.local/share/sbx-kit/portable/ (sbx-kit-state pack/unpack, CLI overlay). sbx-kit on the host only copies the archive; the overlay manifest decides what is in it.

1. Find sbx-kit-state on PATH. If missing, say so — pack/unpack will fail until the overlay is installed (recreate with sbx-kit box run).
2. List what would be packed: workplace dirs, and a DB only if this agent actually has one under an INCLUDE path. Do not assume SQLite or ~/.cursor.
3. Flag anything large or non-portable (build artifacts, caches, secrets in the VM).
4. If the agent is still running against a DB, recommend detach before pack so WALs can flush; if there is no DB, skip that.
5. Do not run pack/unpack yourself unless the human asked.

Reply with: helper present?, sqlite or not, pack-ready or not, and the top paths.`

const summary = `One shot — box summary (handoff).

You are inside a Docker sbx sandbox. Write a handoff another box or agent can reuse. Same prompt text will be shared across boxes later; keep it overlay-driven, not agent-branded.

Cover:
- Project path and what the sandbox is for
- Agent kind / image if you can tell
- Overlay present (floor.md / cli.md) and any extra user kits
- Portable state: what lives under the portable share, SQLite yes/no
- Tools installed that are not in the image
- Open problems and the next useful command on the host (usually sbx-kit box state export / box run)

Write a short markdown note. Do not dump secrets or host paths outside the project.`
