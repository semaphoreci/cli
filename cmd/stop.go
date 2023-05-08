package cmd

import (
	"fmt"

	client "github.com/semaphoreci/cli/api/client"
	"github.com/semaphoreci/cli/cmd/deployment_targets"
	"github.com/semaphoreci/cli/cmd/pipelines"
	"github.com/semaphoreci/cli/cmd/utils"
	"github.com/semaphoreci/cli/cmd/workflows"
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

var StopJobCmd = &cobra.Command{
	Use:     "job [id]",
	Short:   "Stop running job.",
	Long:    ``,
	Aliases: []string{"jobs"},
	Args:    cobra.RangeArgs(1, 1),

	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]

		jobClient := client.NewJobsV1AlphaApi()
		err := jobClient.StopJob(id)

		utils.Check(err)

		fmt.Printf("Job '%s' stopped.\n", id)
	},
}

var StopWfCmd = &cobra.Command{
	Use:     "workflow [id]",
	Short:   "Stop all running pipelines in the workflow.",
	Long:    ``,
	Aliases: []string{"workflows", "wf"},
	Args:    cobra.RangeArgs(1, 1),

	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		workflows.Stop(id)
	},
}

var stopTargetCmd = &cobra.Command{
	Use:     "deployment_target [id]",
	Short:   "Block deployment target in the workflow.",
	Long:    ``,
	Aliases: []string{"deployment_target", "dt", "target", "targets", "tgt", "targ"},
	Args:    cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		targetId := args[0]
		deployment_targets.Stop(targetId)
	},
}

func init() {
	RootCmd.AddCommand(stopCmd)

	stopCmd.AddCommand(StopPplCmd)
	stopCmd.AddCommand(StopJobCmd)
	stopCmd.AddCommand(StopWfCmd)
	stopCmd.AddCommand(stopTargetCmd)
}
