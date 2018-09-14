package cmd

import (
	"io/ioutil"

	"github.com/semaphoreci/cli/cmd/create"
	"github.com/semaphoreci/cli/cmd/handler"
	"github.com/semaphoreci/cli/cmd/utils"

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

	resource, err := parse_yaml_to_map(data)

	utils.CheckWithMessage(err, "Failed to parse resource file.")

	apiVersion := resource["apiVersion"].(string)
	kind := resource["kind"].(string)

	params := handler.CreateParams{ApiVersion: apiVersion, Resource: data}
	handler, err := handler.FindHandler(kind)

	utils.Check(err)

	handler.Create(params)
}
