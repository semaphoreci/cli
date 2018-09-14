package cmd

import (
	"io/ioutil"

	"github.com/semaphoreci/cli/cmd/handler"
	"github.com/semaphoreci/cli/cmd/utils"

	"github.com/spf13/cobra"
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Updates resource based on file.",
	Long:  ``,

	Run: func(cmd *cobra.Command, args []string) {
		RunApply(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)

	desc := "Filename, directory, or URL to files to use to update the resource"
	applyCmd.Flags().StringP("file", "f", "", desc)
}

func RunApply(cmd *cobra.Command, args []string) {
	path, err := cmd.Flags().GetString("file")

	utils.CheckWithMessage(err, "Path not provided")

	data, err := ioutil.ReadFile(path)

	utils.CheckWithMessage(err, "Failed to read from resource file.")

	resource, err := parse_yaml_to_map(data)

	utils.CheckWithMessage(err, "Failed to parse resource file.")

	apiVersion := resource["apiVersion"].(string)
	kind := resource["kind"].(string)

	params := handler.ApplyParams{ApiVersion: apiVersion, Resource: data}
	handler, err := handler.FindHandler(kind)

	utils.Check(err)

	handler.Apply(params)
}
