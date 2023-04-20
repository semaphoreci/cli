package cmd

import (
	"fmt"

	client "github.com/semaphoreci/cli/api/client"
	"github.com/semaphoreci/cli/cmd/utils"
	"github.com/spf13/cobra"
)

var troubleshootCmd = &cobra.Command{
	Use:   "troubleshoot [KIND]",
	Short: "Troubleshoot resource.",
	Long:  ``,
	Args:  cobra.RangeArgs(1, 2),
}

var troubleshootPipelineCmd = &cobra.Command{
	Use:     "pipeline [id]",
	Short:   "Troubleshoot pipeline.",
	Long:    ``,
	Aliases: []string{"pipelines", "ppl"},
	Args:    cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		troubleshootClient := client.NewTroubleshootV1AlphaApi()
		t, err := troubleshootClient.TroubleshootPipeline(id)
		utils.Check(err)

		v, err := t.ToYaml()
		utils.Check(err)
		fmt.Printf("%s", v)
	},
}

var troubleshootJobCmd = &cobra.Command{
	Use:     "job [id]",
	Short:   "Troubleshoot job.",
	Long:    ``,
	Aliases: []string{"jobs"},
	Args:    cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		troubleshootClient := client.NewTroubleshootV1AlphaApi()
		t, err := troubleshootClient.TroubleshootJob(id)
		utils.Check(err)

		v, err := t.ToYaml()
		utils.Check(err)
		fmt.Printf("%s", v)
	},
}

var troubleshootWorkflowCmd = &cobra.Command{
	Use:     "workflow [id]",
	Short:   "Troubleshoot workflow.",
	Long:    ``,
	Aliases: []string{"workflows", "wf"},
	Args:    cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		troubleshootClient := client.NewTroubleshootV1AlphaApi()
		t, err := troubleshootClient.TroubleshootWorkflow(id)
		utils.Check(err)

		v, err := t.ToYaml()
		utils.Check(err)
		fmt.Printf("%s", v)
	},
}

func init() {
	RootCmd.AddCommand(troubleshootCmd)

	troubleshootCmd.AddCommand(troubleshootPipelineCmd)
	troubleshootCmd.AddCommand(troubleshootJobCmd)
	troubleshootCmd.AddCommand(troubleshootWorkflowCmd)
}
