package cmd

import "github.com/spf13/cobra"

func newBoxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "box",
		Aliases: []string{"sandbox"},
		Short:   "Create and manage sbx boxes",
		Long: `Box lifecycle: run, bindings, check, upgrade, rm, and portable state.

Resolve sandboxes with --name or --recipe/--path on subcommands.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(newRunCmd())
	cmd.AddCommand(newBindingsCmd())
	cmd.AddCommand(newCheckCmd())
	cmd.AddCommand(newUpgradeCmd())
	cmd.AddCommand(newRmCmd())
	cmd.AddCommand(newStateCmd())
	return cmd
}
