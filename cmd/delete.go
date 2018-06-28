package cmd

import (
	"fmt"

	"github.com/renderedtext/sem/cmd/handler"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a resource.",
	Long:  ``,
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

	if err != nil {
		fmt.Println(err)
		return
	}

	handler.Delete(params)
}
