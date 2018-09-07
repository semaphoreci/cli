package cmd

import (
	"github.com/renderedtext/sem/cmd/handler"
	"github.com/renderedtext/sem/cmd/utils"

	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit a resource from.",
	Long:  ``,

	Run: func(cmd *cobra.Command, args []string) {
		RunEdit(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}

func RunEdit(cmd *cobra.Command, args []string) {
	kind := args[0]
	name := args[1]

	params := handler.EditParams{Name: name}
	handler, err := handler.FindHandler(kind)

	utils.Check(err)

	handler.Edit(params)
}
