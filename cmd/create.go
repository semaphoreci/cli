package cmd

import (
	"fmt"
	"io/ioutil"
	"strings"

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

			fmt.Printf("Project %s created.\n", project.Metadata.Name)
		case "Secret":
			secret, err := models.NewSecretV1BetaFromYaml(data)

			utils.Check(err)

			c := client.NewSecretV1BetaApi()

			_, err = c.CreateSecret(secret)

			utils.Check(err)

			fmt.Printf("Secret %s created.\n", secret.Metadata.Name)
		case "Dashboard":
			dash, err := models.NewDashboardV1AlphaFromYaml(data)

			utils.Check(err)

			c := client.NewDashboardV1AlphaApi()

			_, err = c.CreateDashboard(dash)

			utils.Check(err)

			fmt.Printf("Dashboard %s created.\n", dash.Metadata.Name)
		default:
			utils.Fail(fmt.Sprintf("Unknown resource kind '%s'", kind))
		}
	},
}

var CreateDashboardCmd = &cobra.Command{
	Use:     "dashboard [NAME]",
	Short:   "Create a dashboard.",
	Long:    ``,
	Aliases: []string{"dashboard", "dash"},
	Args:    cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		c := client.NewDashboardV1AlphaApi()

		dash := models.NewDashboardV1Alpha(name)
		_, err := c.CreateDashboard(&dash)

		utils.Check(err)

		fmt.Printf("Dashboard '%s' created.\n", dash.Metadata.Name)
	},
}

var CreateSecretCmd = &cobra.Command{
	Use:     "secret [NAME]",
	Short:   "Create a secret.",
	Long:    ``,
	Aliases: []string{"secrets"},
	Args:    cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		files, err := cmd.Flags().GetStringArray("file")
		utils.Check(err)

		fileSlices := make([][]string, len(files))

		for i := 0; i < len(files); i++ {
			fileSlices[i] = strings.Split(files[i], ":")
		}

		secret := models.NewSecretV1Beta(name, fileSlices)

		c := client.NewSecretV1BetaApi()

		_, err = c.CreateSecret(&secret)

		utils.Check(err)

		fmt.Printf("Secret '%s' created.\n", secret.Metadata.Name)
	},
}

func init() {
	RootCmd.AddCommand(createCmd)
	createCmd.AddCommand(CreateSecretCmd)
	createCmd.AddCommand(CreateDashboardCmd)

	desc := "Filename, directory, or URL to files to use to create the resource"
	createCmd.Flags().StringP("file", "f", "", desc)

	desc = "File mapping <local-path>:<mount-path>, used to create a secret with file"
	CreateSecretCmd.Flags().StringArrayP("file", "f", []string{}, desc)
}
