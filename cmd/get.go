package cmd

import (
	"fmt"

	"github.com/renderedtext/sem/cmd/handler"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "List of resources.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		RunGet(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}

func RunGet(cmd *cobra.Command, args []string) {
	kind := args[0]

	params := handler.GetParams{}
	handler, err := handler.FindHandler(kind)

	if err != nil {
		fmt.Println(err)
		return
	}

	handler.Get(params)
}
