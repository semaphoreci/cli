package cmd

import (
	"fmt"
  "encoding/json"

	"github.com/renderedtext/sem/client"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a resource of a list of resources.",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
    RunGet(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}

func RunGet(cmd *cobra.Command, args []string) {
  kind := args[0]

  switch kind {
    case "secret", "secrets":
			c := client.FromConfig()

      body, _ := c.List("secrets")

      var secrets []map[string]interface{}

      json.Unmarshal([]byte(body), &secrets)

      fmt.Println("NAME")

      for _, secret := range secrets {
        fmt.Println(secret["metadata"].(map[string]interface{})["name"])
      }
    default:
      panic("Unsupported type")
  }
}
