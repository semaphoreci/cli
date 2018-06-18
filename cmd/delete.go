package cmd

import (
	"fmt"

	"github.com/renderedtext/sem/client"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a resource.",
	Long: ``,
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

  switch kind {
    case "secret", "secrets":
			c := client.FromConfig()

      body, _ := c.Delete("secrets", name)

			fmt.Println(string(body))
    default:
      panic("Unsuported kind")
  }
}
