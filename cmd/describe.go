package cmd

import (
	"fmt"

	"github.com/renderedtext/sem/cmd/handler"
	"github.com/spf13/cobra"
)

// describeCmd represents the describe command
var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Describe a resource",
	Long: ``,

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

  params := handler.DescribeParams { Name: name }
  handler, err := handler.FindHandler(kind)

  if err != nil {
    fmt.Println(err);
    return;
  }

  handler.Describe(params);
}
