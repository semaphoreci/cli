package cmd

import (
	"github.com/semaphoreci/cli/cmd/tasks"
	"github.com/semaphoreci/cli/cmd/utils"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run [KIND]",
	Short: "Run a resource.",
	Long:  ``,
}

var runTaskCmd = &cobra.Command{
	Use:     "task [id]",
	Short:   "Trigger a task run.",
	Long:    ``,
	Aliases: []string{"tasks"},
	Args:    cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]

		branch, err := cmd.Flags().GetString("branch")
		utils.Check(err)
		tag, err := cmd.Flags().GetString("tag")
		utils.Check(err)
		pipelineFile, err := cmd.Flags().GetString("pipeline-file")
		utils.Check(err)
		params, err := cmd.Flags().GetStringSlice("param")
		utils.Check(err)

		tasks.Run(id, branch, tag, pipelineFile, params)
	},
}

func init() {
	RootCmd.AddCommand(runCmd)

	runTaskCmd.Flags().String("branch", "", "git branch to use for the task run")
	runTaskCmd.Flags().String("tag", "", "git tag to use for the task run")
	runTaskCmd.Flags().String("pipeline-file", "", "pipeline file to use for the task run")
	runTaskCmd.Flags().StringSlice("param", []string{}, "parameter in KEY=VALUE format; can be specified multiple times")
	runTaskCmd.MarkFlagsMutuallyExclusive("branch", "tag")

	runCmd.AddCommand(runTaskCmd)
}
