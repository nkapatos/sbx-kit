package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/nkapatos/sbx-kit/cli/internal/sbxcompat"
	"github.com/nkapatos/sbx-kit/cli/internal/sbxutil"
	"github.com/nkapatos/sbx-kit/cli/internal/version"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print sbx-kit version and required sbx CLI range",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			out := cmd.OutOrStdout()
			fmt.Fprintln(out, version.Version)
			fmt.Fprintln(out, sbxcompat.RequirementSummary())

			r := sbxutil.Default()
			raw, err := r.ProbeVersion()
			if err != nil {
				fmt.Fprintf(out, "sbx: not found on PATH (need >= %s)\n", sbxcompat.MinVersion)
				return
			}
			ver, perr := sbxcompat.ParseVersion(raw)
			if perr != nil {
				trunc := strings.TrimSpace(raw)
				if len(trunc) > 80 {
					trunc = trunc[:77] + "..."
				}
				fmt.Fprintf(out, "sbx: unparsed version output %q\n", trunc)
				return
			}
			fmt.Fprintf(out, "sbx: %s\n", ver)
			if err := sbxcompat.Check(raw); err != nil {
				fmt.Fprintf(out, "sbx compat: FAIL\n  %v\n", err)
				return
			}
			fmt.Fprintln(out, "sbx compat: ok")
		},
	}
}
