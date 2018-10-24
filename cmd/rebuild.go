package cmd

import (
	"github.com/semaphoreci/cli/cmd/pipelines"
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

func init() {
	RootCmd.AddCommand(rebuildCmd)

	rebuildCmd.AddCommand(rebuildPplCmd)
}
