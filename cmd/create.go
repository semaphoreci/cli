package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	cmd_create "github.com/semaphoreci/cli/cmd/create"
	"github.com/semaphoreci/cli/cmd/utils"

	client "github.com/semaphoreci/cli/api/client"
	models "github.com/semaphoreci/cli/api/models"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a resource from a file.",
	Long:  ``,

	Run: func(cmd *cobra.Command, args []string) {
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
			project, err := models.NewProjectV1AlphaFromYaml(data)

			utils.Check(err)

			c := client.NewProjectV1AlphaApi()

			_, err = c.CreateProject(project)

			utils.Check(err)

			fmt.Printf("Project %s created.", project.Metadata.Name)
		case "Secret":
			secret, err := models.NewSecretV1BetaFromYaml(data)

			utils.Check(err)

			c := client.NewSecretV1BetaApi()

			_, err = c.CreateSecret(secret)

			utils.Check(err)

			fmt.Printf("Secret %s created.", secret.Metadata.Name)
		case "Dashboard":
			dash, err := models.NewDashboardV1AlphaFromYaml(data)

			utils.Check(err)

			c := client.NewDashboardV1AlphaApi()

			_, err = c.CreateDashboard(dash)

			utils.Check(err)

			fmt.Printf("Dashboard %s created.", dash.Metadata.Name)
		default:
			fmt.Fprintf(os.Stderr, "Unknown resource kind '%s'", kind)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(createCmd)
	createCmd.AddCommand(cmd_create.CreateSecretCmd)
	createCmd.AddCommand(cmd_create.CreateDashboardCmd)

	desc := "Filename, directory, or URL to files to use to create the resource"
	createCmd.Flags().StringP("file", "f", "", desc)
}
