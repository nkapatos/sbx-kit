package kitcreds

import (
	"fmt"
	"io"
	"strings"
)

// PrintHints writes create-time host secret setup lines for services declared in kits.
func PrintHints(w io.Writer, recipe string, needs []Need) {
	if len(needs) == 0 {
		return
	}
	fmt.Fprintln(w, "==> host secrets declared by this recipe's kits (set any you use)")
	if recipe != "" {
		fmt.Fprintf(w, "    recipe=%s\n", recipe)
	}
	for _, n := range needs {
		from := n.KitName
		if from == "" {
			from = n.KitPath
		}
		fmt.Fprintf(w, "  sbx secret set %s", n.Service)
		if len(n.Envs) > 0 {
			fmt.Fprintf(w, "    # %s ← kit %s\n", strings.Join(n.Envs, ", "), from)
		} else {
			fmt.Fprintf(w, "    # kit %s\n", from)
		}
	}
	fmt.Fprintln(w, "  # Not all services are required. Extra APIs belong in the kit (or a mixin),")
	fmt.Fprintln(w, "  # not a kit-per-provider. Global secrets apply at create; recreate after changes.")
	fmt.Fprintln(w, "  # Optional: pass show api/<service> | sbx secret set <service>")
}
