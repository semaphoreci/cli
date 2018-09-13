package cmd

import (
	"github.com/spf13/cobra"

	"github.com/semaphoreci/cli/cmd/edit"
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit a resource from.",
	Long:  ``,
}

func init() {
	rootCmd.AddCommand(editCmd)

	editCmd.AddCommand(cmd_edit.EditSecretCmd)
	editCmd.AddCommand(cmd_edit.EditDashboardCmd)
}
