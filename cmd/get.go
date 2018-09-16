package cmd

import (
	cmd_get "github.com/semaphoreci/cli/cmd/get"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get [KIND]",
	Short: "List resources.",
	Long:  ``,
	Args:  cobra.RangeArgs(1, 2),
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.AddCommand(cmd_get.GetDashboardCmd)
	getCmd.AddCommand(cmd_get.GetSecretCmd)
	getCmd.AddCommand(cmd_get.GetProjectCmd)
}
