package cmd

import (
	"io/ioutil"

	"github.com/renderedtext/sem/cmd/create"
	"github.com/renderedtext/sem/cmd/handler"
	"github.com/renderedtext/sem/cmd/utils"

	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a resource from a file.",
	Long:  ``,

	Run: func(cmd *cobra.Command, args []string) {
		RunCreate(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.AddCommand(cmd_create.CreateSecretCmd)

	desc := "Filename, directory, or URL to files to use to create the resource"
	createCmd.Flags().StringP("file", "f", "", desc)
}

func RunCreate(cmd *cobra.Command, args []string) {
	path, err := cmd.Flags().GetString("file")

	utils.CheckWithMessage(err, "Path not provided")

	data, err := ioutil.ReadFile(path)

	utils.CheckWithMessage(err, "Failed to read from resource file.")

	resource, err := parse(data)

	utils.CheckWithMessage(err, "Failed to parse resource file.")

	apiVersion := resource["apiVersion"].(string)
	kind := resource["kind"].(string)

	params := handler.CreateParams{ApiVersion: apiVersion, Resource: data}
	handler, err := handler.FindHandler(kind)

	utils.Check(err)

	handler.Create(params)
}

func parse(data []byte) (map[string]interface{}, error) {
	m := make(map[string]interface{})

	err := yaml.Unmarshal(data, &m)

	return m, err
}
