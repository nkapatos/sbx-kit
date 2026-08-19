package cmd

import (
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
			w := UI()
			w.Println(version.Version)
			w.Println(sbxcompat.RequirementSummary())

			r := sbxutil.Default()
			raw, err := r.ProbeVersion()
			if err != nil {
				w.Printf("sbx: not found on PATH (need >= %s)\n", sbxcompat.MinVersion)
				return
			}
			ver, perr := sbxcompat.ParseVersion(raw)
			if perr != nil {
				trunc := strings.TrimSpace(raw)
				if len(trunc) > 80 {
					trunc = trunc[:77] + "..."
				}
				w.Printf("sbx: unparsed version output %q\n", trunc)
				return
			}
			w.Printf("sbx: %s\n", ver)
			if err := sbxcompat.Check(raw); err != nil {
				w.Printf("sbx compat: FAIL\n  %v\n", err)
				return
			}
			w.Println("sbx compat: ok")
		},
	}
}
