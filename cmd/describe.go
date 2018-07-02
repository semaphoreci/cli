package cmd

import (
	"github.com/renderedtext/sem/cmd/handler"
	"github.com/renderedtext/sem/cmd/utils"

	"github.com/spf13/cobra"
)

// describeCmd represents the describe command
var describeCmd = &cobra.Command{
	Use:   "describe [KIND] [NAME]",
	Short: "Describe a resource",
	Long:  ``,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		RunDescribe(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(describeCmd)
}

func RunDescribe(cmd *cobra.Command, args []string) {
	kind := args[0]
	name := args[1]

	params := handler.DescribeParams{Name: name}
	handler, err := handler.FindHandler(kind)

	utils.Check(err)

	handler.Describe(params)
}
