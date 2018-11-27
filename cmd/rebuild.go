package cmd

import (
	"github.com/semaphoreci/cli/cmd/pipelines"
	"github.com/semaphoreci/cli/cmd/workflows"
	"github.com/spf13/cobra"
)

var rebuildCmd = &cobra.Command{
	Use:   "rebuild [KIND]",
	Short: "Rebuild workflow or pipeline.",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
}

var rebuildPplCmd = &cobra.Command{
	Use:     "pipeline [id]",
	Short:   "Partially rebuild failed pipeline.",
	Long:    ``,
	Aliases: []string{"pipelines", "ppl"},
	Args:    cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		pipelines.Rebuild(id)
	},
}

var rebuildWfCmd = &cobra.Command{
	Use:     "workflow [id]",
	Short:   "Rebuild workflow.",
	Long:    `Create new workflow, as if new new github push arrived.`,
	Aliases: []string{"workflows", "wf"},
	Args:    cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		workflows.Rebuild(id)
	},
}

func init() {
	RootCmd.AddCommand(rebuildCmd)

	rebuildCmd.AddCommand(rebuildPplCmd)
	rebuildCmd.AddCommand(rebuildWfCmd)
}
