package cmd

import (
	"github.com/spf13/cobra"

	"github.com/renderedtext/sem/cmd/edit"
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit a resource from.",
	Long:  ``,
}

func init() {
	rootCmd.AddCommand(editCmd)

	editCmd.AddCommand(cmd_edit.EditSecretCmd)
}
