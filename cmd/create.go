package cmd

import (
	"fmt"
  "io/ioutil"

	"github.com/spf13/cobra"
  "github.com/ghodss/yaml"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a resource from a file.",
	Long: ``,

	Run: func(cmd *cobra.Command, args []string) {
    path, _ := cmd.Flags().GetString("file")

    fmt.Printf("%+v\n", path)

    data, err := ioutil.ReadFile(path)

    if err != nil {
      fmt.Printf("%+v\n", err)
      return
    }

    parse(data)

	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	createCmd.PersistentFlags().StringP("file", "f", "", "Filename, directory, or URL to files to use to create the resource")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func parse(data []byte) {
  j, err := yaml.YAMLToJSON(data)

  if err != nil {
    fmt.Printf("%#v\n", err)
  }

  fmt.Printf("%s\n", j)
}
