package cmd

import (
	"errors"
	"fmt"

	client "github.com/semaphoreci/cli/api/client"
	models "github.com/semaphoreci/cli/api/models"
	"github.com/semaphoreci/cli/cmd/utils"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit a resource from.",
	Long:  ``,
}

var EditDashboardCmd = &cobra.Command{
	Use:     "dashboard [name]",
	Short:   "Edit a dashboard.",
	Long:    ``,
	Aliases: []string{"dashboards", "dash"},
	Args:    cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		c := client.NewDashboardV1AlphaApi()

		dashboard, err := c.GetDashboard(name)

		utils.Check(err)

		content, err := dashboard.ToYaml()

		utils.Check(err)

		new_content, err := utils.EditYamlInEditor(dashboard.ObjectName(), string(content))

		utils.Check(err)

		updated_dashboard, err := models.NewDashboardV1AlphaFromYaml([]byte(new_content))

		utils.Check(err)

		dashboard, err = c.UpdateDashboard(updated_dashboard)

		utils.Check(err)

		fmt.Printf("Dashboard '%s' updated.\n", dashboard.Metadata.Name)
	},
}

var EditNotificationCmd = &cobra.Command{
	Use:     "notification [name]",
	Short:   "Edit a notification.",
	Long:    ``,
	Aliases: []string{"notifications", "notifs", "notif"},
	Args:    cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		c := client.NewNotificationsV1AlphaApi()

		notif, err := c.GetNotification(name)

		utils.Check(err)

		content, err := notif.ToYaml()

		utils.Check(err)

		newContent, err := utils.EditYamlInEditor(notif.ObjectName(), string(content))

		utils.Check(err)

		updatedNotif, err := models.NewNotificationV1AlphaFromYaml([]byte(newContent))

		utils.Check(err)

		notif, err = c.UpdateNotification(updatedNotif)

		utils.Check(err)

		fmt.Printf("Notification '%s' updated.\n", notif.Metadata.Name)
	},
}

var EditSecretCmd = &cobra.Command{
	Use:     "secret [name]",
	Short:   "Edit a secret.",
	Long:    ``,
	Aliases: []string{"secrets"},
	Args:    cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		projectID := GetProjectID(cmd)

		if projectID == "" {
			name := args[0]

			c := client.NewSecretV1BetaApi()

			secret, err := c.GetSecret(name)

			utils.Check(err)

			content, err := secret.ToYaml()

			utils.Check(err)

			new_content, err := utils.EditYamlInEditor(secret.ObjectName(), string(content))

			utils.Check(err)

			updated_secret, err := models.NewSecretV1BetaFromYaml([]byte(new_content))

			utils.Check(err)

			secret, err = c.UpdateSecret(updated_secret)

			utils.Check(err)

			fmt.Printf("Secret '%s' updated.\n", secret.Metadata.Name)
		} else {
			name := args[0]

			c := client.NewProjectSecretV1Api(projectID)

			secret, err := c.GetSecret(name)

			utils.Check(err)

			content, err := secret.ToYaml()

			utils.Check(err)

			new_content, err := utils.EditYamlInEditor(secret.ObjectName(), string(content))

			utils.Check(err)

			updated_secret, err := models.NewProjectSecretV1FromYaml([]byte(new_content))

			utils.Check(err)

			secret, err = c.UpdateSecret(updated_secret)

			utils.Check(err)

			fmt.Printf("Secret '%s' updated.\n", secret.Metadata.Name)
		}
	},
}

var EditProjectCmd = &cobra.Command{
	Use:     "project [name]",
	Short:   "Edit a project.",
	Long:    ``,
	Aliases: []string{"project", "prj"},
	Args:    cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		c := client.NewProjectV1AlphaApi()

		project, err := c.GetProject(name)

		utils.Check(err)

		content, err := project.ToYaml()

		utils.Check(err)

		new_content, err := utils.EditYamlInEditor(project.ObjectName(), string(content))

		utils.Check(err)

		updated_project, err := models.NewProjectV1AlphaFromYaml([]byte(new_content))

		utils.Check(err)

		project, err = c.UpdateProject(updated_project)

		utils.Check(err)

		fmt.Printf("Project '%s' updated.\n", project.Metadata.Name)
	},
}

var EditDeploymentTargetCmd = &cobra.Command{
	Use:     "deployment_target [name]",
	Short:   "Edit a deployment target.",
	Long:    ``,
	Aliases: []string{"deployment_target", "dt", "target", "targets", "tgt", "targ"},
	Args:    cobra.RangeArgs(0, 1),

	Run: func(cmd *cobra.Command, args []string) {
		projectId := GetProjectID(cmd)
		c := client.NewDeploymentTargetsV1AlphaApi()
		var target *models.DeploymentTargetV1Alpha
		var err error
		if len(args) == 1 {
			targetId := args[0]
			target, err = c.Describe(targetId, projectId)
			utils.Check(err)
		} else {
			targetName, err := cmd.Flags().GetString("target-name")
			utils.Check(err)
			if targetName == "" {
				utils.Check(errors.New("target id or name must be present"))
			}
			target, err = c.DescribeByName(targetName, projectId)
			utils.Check(err)
		}
		if target == nil {
			utils.Check(errors.New("target must be present"))
			return
		}
		// TODO: Load secrets for deployment target to avoid requiring clients
		// to provide them every time they edit the request.
		request := models.DeploymentTargetUpdateRequestV1Alpha{
			ProjectId:               projectId,
			DeploymentTargetV1Alpha: *target,
		}
		content, err := request.ToYaml()
		utils.Check(err)

		new_content, err := utils.EditYamlInEditor(request.ObjectName(), string(content))
		utils.Check(err)

		update_request, err := models.NewDeploymentTargetUpdateRequestV1AlphaFromYaml([]byte(new_content))
		utils.Check(err)

		update_request.Id = target.Id
		update_request.ProjectId = target.ProjectId
		updatedTarget, err := c.Update(update_request)
		utils.Check(err)

		fmt.Printf("Deployment target '%s' updated.\n", updatedTarget.Name)
	},
}

func init() {
	RootCmd.AddCommand(editCmd)

	EditSecretCmd.Flags().StringP("project-name", "p", "",
		"project name; if specified will edit project level secret, otherwise organization secret")
	EditSecretCmd.Flags().StringP("project-id", "i", "",
		"project id; if specified will edit project level secret, otherwise organization secret")
	editCmd.AddCommand(EditSecretCmd)
	editCmd.AddCommand(EditDashboardCmd)
	editCmd.AddCommand(EditNotificationCmd)
	editCmd.AddCommand(EditProjectCmd)

	EditDeploymentTargetCmd.Flags().StringP("project-name", "p", "",
		"project name; if specified will edit project level secret, otherwise organization secret")
	EditDeploymentTargetCmd.Flags().StringP("project-id", "i", "",
		"project id; if specified will edit project level secret, otherwise organization secret")
	EditDeploymentTargetCmd.Flags().StringP("target-name", "n", "", "target name")
	editCmd.AddCommand(EditDeploymentTargetCmd)
}
