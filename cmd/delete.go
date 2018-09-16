package cmd

import (
	cmd_delete "github.com/semaphoreci/cli/cmd/delete"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [KIND] [NAME]",
	Short: "Delete a resource.",
	Long:  ``,
	Args:  cobra.ExactArgs(2),
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	deleteCmd.AddCommand(cmd_delete.DeleteDashboardCmd)
	deleteCmd.AddCommand(cmd_delete.DeleteProjectCmd)
	deleteCmd.AddCommand(cmd_delete.DeleteSecretCmd)
}
