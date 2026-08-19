package cmd

import (
	"github.com/spf13/cobra"

	"github.com/nkapatos/sbx-kit/cli/internal/boxprompt"
	"github.com/nkapatos/sbx-kit/cli/internal/experimental"
	"github.com/nkapatos/sbx-kit/cli/internal/recipespec"
)

func newExperimentalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "experimental",
		Short: "Parked helpers (spec, in-box prompts)",
		Long: `Work in progress that is not part of the stable CLI surface.

Agent skill: sbx-kit recipes skill
Recipe verify: sbx-kit recipes verify
In-box prompts: sbx-kit experimental prompts`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(newExperimentalSpecCmd())
	cmd.AddCommand(newExperimentalPromptsCmd())
	return cmd
}

func newExperimentalSpecCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "spec",
		Short: "Show recipe spec status (parked)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			UI().Println(recipespec.Status)
			UI().Println("Agent skill: sbx-kit recipes skill")
			return experimental.ErrNotReady{Feature: "recipe spec", Track: "spec"}
		},
	}
}

func newExperimentalPromptsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "prompts [name]",
		Short: "Dump parked in-box one-shot prompts",
		Long: `Parked prompts for the agent inside a box (maintenance, cleanup, state, summary).

Paste into the attached agent. Later these texts are shared between boxes
via the CLI overlay and statexfer — not copied by hand.

  sbx-kit experimental prompts
  sbx-kit experimental prompts summary`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				UI().Println(boxprompt.Status)
				for _, p := range boxprompt.All {
					UI().Println()
					UI().Header(p.Name)
					UI().Println(p.Body)
				}
				return experimental.ErrNotReady{Feature: "in-box prompts sharing", Track: "prompts"}
			}
			p, err := boxprompt.Lookup(args[0])
			if err != nil {
				return err
			}
			UI().Println(p.Body)
			return experimental.ErrNotReady{Feature: "in-box prompts sharing", Track: "prompts"}
		},
	}
}
