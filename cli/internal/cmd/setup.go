package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/nkapatos/sbx-kit/cli/internal/toolkit"
)

func newSetupCmd() *cobra.Command {
	var catalogPath string
	cmd := &cobra.Command{
		Use:   "setup",
		Short: "Configure the catalog path",
		Long: `Writes the catalog path to ~/.config/sbx-kit/config.yaml.

The catalog holds directories of recipes (each with recipes/, and optional
kits/ and images/).

With no --catalog, setup asks for the path (default ~/sbx-kit-catalog).
SBX_KIT_CATALOG overrides the configured path for one shell.`,
		Example: `  sbx-kit setup
  sbx-kit setup --catalog ~/sbx-kit-catalog`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := catalogPath
			if dir == "" {
				dir = os.Getenv(toolkit.CatalogEnv)
			}
			if dir == "" {
				if !stdinIsTTY() {
					return fmt.Errorf("pass --catalog <dir> (or run setup from a terminal)")
				}
				def := toolkit.DefaultCatalog()
				if existing, err := toolkit.ConfiguredCatalog(); err == nil && existing != "" {
					def = existing
				}
				var err error
				dir, err = promptLine(cmd.InOrStdin(), cmd.OutOrStdout(), "Catalog path:", def)
				if err != nil {
					return err
				}
				if dir == "" {
					return fmt.Errorf("catalog path is required")
				}
			}
			abs, err := toolkit.WriteCatalog(dir)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "catalog: %s\n", abs)
			fmt.Fprintf(cmd.OutOrStdout(), "add a directory:  sbx-kit catalog add <url>\n")
			fmt.Fprintf(cmd.OutOrStdout(), "optional: export %s=%s\n", toolkit.CatalogEnv, abs)
			return nil
		},
	}
	cmd.Flags().StringVar(&catalogPath, "catalog", "", "catalog path (parent of recipe directories)")
	return cmd
}
