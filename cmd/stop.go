package cmd

import (
	"fmt"

	client "github.com/semaphoreci/cli/api/client"
	"github.com/semaphoreci/cli/cmd/pipelines"
	"github.com/semaphoreci/cli/cmd/utils"
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

func init() {
	RootCmd.AddCommand(stopCmd)

	stopCmd.AddCommand(StopPplCmd)
	stopCmd.AddCommand(StopJobCmd)
}
