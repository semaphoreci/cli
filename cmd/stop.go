package cmd

import (
	"github.com/semaphoreci/cli/cmd/pipelines"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop [KIND]",
	Short: "Stop resource execution.",
	Long:  ``,
	Args:  cobra.RangeArgs(1, 2),
}

var StopPplCmd = &cobra.Command{
	Use:     "pipeline [id]",
	Short:   "Stop running pipeline.",
	Long:    ``,
	Aliases: []string{"pipelines", "ppl"},
	Args:    cobra.RangeArgs(1, 1),

	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		pipelines.Stop(id)
	},
}

func init() {
	RootCmd.AddCommand(stopCmd)

	stopCmd.AddCommand(StopPplCmd)
}
