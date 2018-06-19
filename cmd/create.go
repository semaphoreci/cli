package cmd

import (
	"fmt"
  "os"
  "io/ioutil"

	"github.com/renderedtext/sem/client"
	"github.com/spf13/cobra"
  "github.com/ghodss/yaml"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a resource from a file.",
	Long: ``,

	Run: func(cmd *cobra.Command, args []string) {
    RunCreate(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

  desc := "Filename, directory, or URL to files to use to create the resource"
	createCmd.PersistentFlags().StringP("file", "f", "", desc)
}

func RunCreate(cmd *cobra.Command, args []string) {
    path, err := cmd.Flags().GetString("file")

    check(err, "Path not provided")

    data, err := ioutil.ReadFile(path)

    check(err, "Failed to read from resource file.")

    resource, err := parse(data)

    check(err, "Failed to parse resource file.")

    apiVersion := resource["apiVersion"].(string)
    kind := resource["kind"].(string)

    switch kind {
      case "Secret":
        c := client.FromConfig()
        c.SetApiVersion(apiVersion)

        resource, _ := yaml.YAMLToJSON(data)

        body, _ := c.Post("secrets", resource)

        j, _ := yaml.JSONToYAML(body)

        fmt.Println(string(j))
      default:
        panic("Unsuported kind")
    }
}

func parse(data []byte) (map[string]interface{}, error) {
  m := make(map[string]interface{})

  // fmt.Print(string(data))

  err := yaml.Unmarshal(data, &m)

  return m, err
}

func check(err error, message string) {
  if err != nil {
    fmt.Fprintf(os.Stderr, "%s\n", message)
    fmt.Fprintf(os.Stderr, "error: %v\n", err)

    os.Exit(1)
  }
}
