package cmd

import (
	"github.com/semaphoreci/cli/cmd/handler"
	"github.com/semaphoreci/cli/cmd/utils"

	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete [KIND] [NAME]",
	Short: "Delete a resource.",
	Long:  ``,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		RunDelete(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}

func RunDelete(cmd *cobra.Command, args []string) {
	kind := args[0]
	name := args[1]

	params := handler.DeleteParams{Name: name}
	handler, err := handler.FindHandler(kind)

	utils.Check(err)

	handler.Delete(params)
}
