package cmd

import (
	"fmt"

	"github.com/renderedtext/sem/client"
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

  switch kind {
    case "secret", "secrets":
			c := client.FromConfig()

			body, _ := c.Get("secrets", name)

			fmt.Println(string(body))
    default:
      panic("Unsuported kind")
  }
}
