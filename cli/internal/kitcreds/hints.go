package kitcreds

import (
	"fmt"
	"io"
	"strings"
)

// PrintHints writes novice-friendly sbx secret guidance. Does not store secrets.
func PrintHints(w io.Writer, recipe string, needs []Need) {
	if len(needs) == 0 {
		return
	}
	fmt.Fprintln(w, "==> host secrets this recipe expects (keys stay on the host; box sees a sentinel)")
	if recipe != "" {
		fmt.Fprintf(w, "    recipe=%s\n", recipe)
	}
	for _, n := range needs {
		from := n.KitName
		if from == "" {
			from = n.KitPath
		}
		fmt.Fprintf(w, "  sbx secret set %s\n", n.Service)
		if len(n.Envs) > 0 {
			fmt.Fprintf(w, "    # in-box env: %s  (from kit %s)\n", strings.Join(n.Envs, ", "), from)
		} else {
			fmt.Fprintf(w, "    # from kit %s\n", from)
		}
		fmt.Fprintf(w, "    # or:  pass show api/%s | sbx secret set %s\n", n.Service, n.Service)
	}
	fmt.Fprintln(w, "  # Global secrets apply at create — recreate the sandbox after changing them.")
	fmt.Fprintln(w, "  # Prefer sbx secret set-custom only when the agent validates key format / body injection.")
	fmt.Fprintln(w, "  # See: sbx secret --help")
}
