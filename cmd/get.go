package cmd

import (
	cmd_get "github.com/semaphoreci/cli/cmd/get"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get [KIND]",
	Short: "List of resources.",
	Long:  ``,
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		RunGet(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.AddCommand(cmd_get.GetDashboardCmd)
}

func RunGet(cmd *cobra.Command, args []string) {
	// kind := args[0]

	// if len(args) == 1 {
	// 	params := handler.GetParams{}
	// 	handler, err := handler.FindHandler(kind)

	// 	utils.Check(err)

	// 	handler.Get(params)
	// } else {
	// 	name := args[1]

	// 	params := handler.DescribeParams{Name: name}
	// 	handler, err := handler.FindHandler(kind)

	// 	utils.Check(err)

	// 	handler.Describe(params)
	// }

}
