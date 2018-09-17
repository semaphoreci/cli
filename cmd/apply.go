package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	client "github.com/semaphoreci/cli/api/client"
	models "github.com/semaphoreci/cli/api/models"
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
	RootCmd.AddCommand(applyCmd)

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

	// apiVersion := resource["apiVersion"].(string)
	kind := resource["kind"].(string)

	switch kind {
	case "Project":
		fmt.Fprintln(os.Stderr, "Unsupported action for Projects")
		os.Exit(1)
	case "Secret":
		secret, err := models.NewSecretV1BetaFromYaml(data)

		utils.Check(err)

		c := client.NewSecretV1BetaApi()

		secret, err = c.UpdateSecret(secret)

		utils.Check(err)

		fmt.Printf("Secret %s created.", secret.Metadata.Name)
	case "Dashboard":
		dash, err := models.NewDashboardV1AlphaFromYaml(data)

		utils.Check(err)

		c := client.NewDashboardV1AlphaApi()

		dash, err = c.UpdateDashboard(dash)

		utils.Check(err)

		fmt.Printf("Dashboard %s created.", dash.Metadata.Name)
	default:
		fmt.Fprintf(os.Stderr, "Unknown resource kind '%s'", kind)
		os.Exit(1)
	}
}
