package cmd

import (
	"fmt"

	"github.com/renderedtext/sem/client"
	"github.com/spf13/cobra"
  "github.com/ghodss/yaml"
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
      j, _ := yaml.JSONToYAML(body)

      fmt.Println(string(j))
    default:
      panic("Unsuported kind")
  }
}
