package cmd

import (
	"github.com/semaphoreci/cli/cmd/tasks"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run [KIND]",
	Short: "Run a task.",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
}

var runTaskCmd = &cobra.Command{
	Use:     "task [id]",
	Short:   "Trigger a task run.",
	Long:    ``,
	Aliases: []string{"tasks"},
	Args:    cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]

		branch, _ := cmd.Flags().GetString("branch")
		tag, _ := cmd.Flags().GetString("tag")
		pipelineFile, _ := cmd.Flags().GetString("pipeline-file")
		params, _ := cmd.Flags().GetStringSlice("param")

		tasks.Run(id, branch, tag, pipelineFile, params)
	},
}

func init() {
	RootCmd.AddCommand(runCmd)

	runTaskCmd.Flags().String("branch", "", "git branch to use for the task run")
	runTaskCmd.Flags().String("tag", "", "git tag to use for the task run")
	runTaskCmd.Flags().String("pipeline-file", "", "pipeline file to use for the task run")
	runTaskCmd.Flags().StringSlice("param", []string{}, "parameter in KEY=VALUE format; can be specified multiple times")

	runCmd.AddCommand(runTaskCmd)
}
