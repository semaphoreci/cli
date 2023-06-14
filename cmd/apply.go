package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"

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

	_, kind, err := utils.ParseYamlResourceHeaders(data)

	utils.Check(err)

	switch kind {
	case "Secret":
		secret, err := models.NewSecretV1BetaFromYaml(data)

		utils.Check(err)

		c := client.NewSecretV1BetaApi()

		secret, err = c.UpdateSecret(secret)

		utils.Check(err)

		fmt.Printf("Secret '%s' updated.\n", secret.Metadata.Name)
	case "ProjectSecret":
		secret, err := models.NewProjectSecretV1FromYaml(data)

		utils.Check(err)

		c := client.NewProjectSecretV1Api(secret.Metadata.ProjectIdOrName)

		secret, err = c.UpdateSecret(secret)

		utils.Check(err)

		fmt.Printf("Secret '%s' created in project '%s'.\n", secret.Metadata.Name, secret.Metadata.ProjectIdOrName)
	case "Dashboard":
		dash, err := models.NewDashboardV1AlphaFromYaml(data)

		utils.Check(err)

		c := client.NewDashboardV1AlphaApi()

		dash, err = c.UpdateDashboard(dash)

		utils.Check(err)

		fmt.Printf("Dashboard '%s' updated.\n", dash.Metadata.Name)
	case "Notification":
		notif, err := models.NewNotificationV1AlphaFromYaml(data)

		utils.Check(err)

		c := client.NewNotificationsV1AlphaApi()

		notif, err = c.UpdateNotification(notif)

		utils.Check(err)

		fmt.Printf("Notification '%s' updated.\n", notif.Metadata.Name)
	case "Project":
		proj, err := models.NewProjectV1AlphaFromYaml(data)

		utils.Check(err)

		c := client.NewProjectV1AlphaApi()

		proj, err = c.UpdateProject(proj)

		utils.Check(err)

		fmt.Printf("Project '%s' updated.\n", proj.Metadata.Name)
	case models.DeploymentTargetKindV1Alpha:
		target, err := models.NewDeploymentTargetV1AlphaFromYaml(data)
		utils.Check(err)
		if target == nil {
			utils.Check(errors.New("deployment target in the file is empty"))
			return
		}
		updateRequest := &models.DeploymentTargetUpdateRequestV1Alpha{
			DeploymentTargetV1Alpha: *target,
		}
		utils.Check(updateRequest.LoadFiles())

		c := client.NewDeploymentTargetsV1AlphaApi()
		updatedDeploymentTarget, err := c.Update(updateRequest)
		utils.Check(err)

		fmt.Printf("Deployment target '%s' ('%s') updated.\n", updatedDeploymentTarget.Id, updatedDeploymentTarget.Name)
	default:
		utils.Fail(fmt.Sprintf("Unsuported resource kind '%s'", kind))
	}
}
