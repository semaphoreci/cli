package cmd

import (
	"fmt"
  "io/ioutil"

	"github.com/spf13/cobra"
  // "gopkg.in/yaml.v2"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a resource from a file.",
	Long: ``,

	Run: func(cmd *cobra.Command, args []string) {
    path, _ := cmd.Flags().GetString("file")

    fmt.Printf("%+v\n", path)

    dat, err := ioutil.ReadFile(path)

    if err != nil {
      fmt.Printf("%+v\n", err)
      return
    }

    fmt.Print(string(dat))
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
