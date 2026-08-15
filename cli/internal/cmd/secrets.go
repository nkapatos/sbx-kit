package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/nkapatos/sbx-kit/cli/internal/catalog"
	"github.com/nkapatos/sbx-kit/cli/internal/kitcreds"
)

func newSecretsCmd() *cobra.Command {
	var agent string

	cmd := &cobra.Command{
		Use:   "secrets",
		Short: "Show sbx secret commands implied by a recipe's kits",
		Long: `Reads credential declarations from the kits in a catalog recipe and prints
the matching host-side sbx secret set … commands.

sbx-kit does not store secrets — that is sbx's job. This is a convenience
guide so novices do not dump API keys into the sandbox.`,
		Example: `  sbx-kit secrets --agent shell-hub
  sbx-kit secrets --agent cursor-hub`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if agent == "" {
				return fmt.Errorf("--agent is required")
			}
			root, err := requireToolkitRoot()
			if err != nil {
				return err
			}
			cat, err := catalog.Load(filepath.Join(root, "config", "agents.yaml"))
			if err != nil {
				return err
			}
			ag, ok := cat.Agents[agent]
			if !ok {
				return fmt.Errorf("unknown recipe %q (try: sbx-kit agents)", agent)
			}
			kits := ag.Kits
			if len(kits) == 0 {
				kits = cat.Defaults.Kits
			}
			kitPaths := make([]string, 0, len(kits))
			for _, k := range kits {
				kitPaths = append(kitPaths, filepath.Join(root, "kits", k))
			}
			needs, err := kitcreds.ScanSpecs(kitPaths)
			if err != nil {
				return err
			}
			if len(needs) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "No credential services declared in kits for recipe %q.\n", agent)
				return nil
			}
			kitcreds.PrintHints(cmd.OutOrStdout(), agent, needs)
			return nil
		},
	}

	cmd.Flags().StringVar(&agent, "agent", "", "catalog recipe id (required)")
	_ = cmd.MarkFlagRequired("agent")
	return cmd
}
