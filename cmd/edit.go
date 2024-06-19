package cmd

import (
	"errors"
	"fmt"

	client "github.com/semaphoreci/cli/api/client"
	models "github.com/semaphoreci/cli/api/models"
	"github.com/semaphoreci/cli/api/uuid"
	"github.com/semaphoreci/cli/cmd/utils"
	"github.com/spf13/cobra"
)

const secretEditDangerMessage = `
# WARNING! Secrets cannot be updated, only replaced. Once the change is applied, the old values will be lost forever.
# Note: You can exit without saving to skip.

`
const secretAskConfirmationMessage = `WARNING! Secrets cannot be updated, only replaced. Once the change is applied, the old values will be lost forever. To continue, please type in the (current) secret name:`

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

			if !secret.Editable() {
				notice := []byte(secretEditDangerMessage)
				content = append(notice, content...)
			}

			new_content, err := utils.EditYamlInEditor(secret.ObjectName(), string(content))

			utils.Check(err)

			updated_secret, err := models.NewSecretV1BetaFromYaml([]byte(new_content))

			utils.Check(err)

			if secret.Editable() {
				secret, err = c.UpdateSecret(updated_secret)
			} else {
				cmd.Println(secretAskConfirmationMessage)
				err = utils.Ask(secret.Metadata.Name)
				if err == nil {
					secret, err = c.FallbackUpdate(updated_secret)
				}
			}

			utils.Check(err)

			fmt.Printf("Secret '%s' updated.\n", secret.Metadata.Name)
		} else {
			name := args[0]

			c := client.NewProjectSecretV1Api(projectID)

			secret, err := c.GetSecret(name)

			utils.Check(err)

			content, err := secret.ToYaml()

			utils.Check(err)

			if !secret.Editable() {
				notice := []byte(secretEditDangerMessage)
				content = append(notice, content...)
			}

			new_content, err := utils.EditYamlInEditor(secret.ObjectName(), string(content))

			utils.Check(err)

			updated_secret, err := models.NewProjectSecretV1FromYaml([]byte(new_content))

			utils.Check(err)

			if secret.Editable() {
				secret, err = c.UpdateSecret(updated_secret)
			} else {
				cmd.Println(secretAskConfirmationMessage)
				err = utils.Ask(secret.Metadata.Name)
				if err == nil {
					secret, err = c.FallbackUpdate(updated_secret)
				}
			}
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
	Use:     "deployment_target [id or name]",
	Short:   "Edit a deployment target.",
	Long:    ``,
	Aliases: models.DeploymentTargetCmdAliases,
	Args:    cobra.RangeArgs(0, 1),

	Run: func(cmd *cobra.Command, args []string) {
		c := client.NewDeploymentTargetsV1AlphaApi()
		targetName, err := cmd.Flags().GetString("name")
		utils.Check(err)
		targetId, err := cmd.Flags().GetString("id")
		utils.Check(err)
		if len(args) == 1 {
			if uuid.IsValid(args[0]) {
				targetId = args[0]
			} else {
				targetName = args[0]
			}
		}
		shouldActivate, err := cmd.Flags().GetBool("activate")
		utils.Check(err)
		shouldDeactivate, err := cmd.Flags().GetBool("deactivate")
		utils.Check(err)

		var target *models.DeploymentTargetV1Alpha
		if targetId != "" {
			if !shouldActivate && !shouldDeactivate {
				target, err = c.DescribeWithSecrets(targetId)
			}
		} else if targetName != "" {
			target, err = c.DescribeByName(targetName, getPrj(cmd))
			if err == nil && target != nil {
				target, err = c.DescribeWithSecrets(target.Id)
			}
		} else {
			err = errors.New("target id or target name must be provided")
		}
		utils.Check(err)
		if target != nil {
			targetId = target.Id
		}
		if shouldActivate {
			succeeded, err := c.Activate(targetId)
			utils.Check(err)
			if !succeeded {
				utils.Check(errors.New("the deployment target wasn't activated successfully"))
			} else {
				fmt.Printf("The deployment target '%s' is active.\n", targetId)
				return
			}
		} else if shouldDeactivate {
			succeeded, err := c.Deactivate(targetId)
			utils.Check(err)
			if !succeeded {
				utils.Check(errors.New("the deployment target wasn't deactivated successfully"))
			} else {
				fmt.Printf("The deployment target '%s' is inactive.\n", targetId)
				return
			}
		}
		if target == nil {
			utils.Check(errors.New("valid target could not be retrieved"))
		}
		// TODO: Load secrets for deployment target to avoid requiring clients
		// to provide them every time they edit the request.
		request := models.DeploymentTargetUpdateRequestV1Alpha{
			DeploymentTargetV1Alpha: *target,
		}
		content, err := request.ToYaml()
		utils.Check(err)

		new_content, err := utils.EditYamlInEditor(request.ObjectName(), string(content))
		utils.Check(err)

		changedTarget, err := models.NewDeploymentTargetV1AlphaFromYaml([]byte(new_content))
		utils.Check(err)
		utils.Check(changedTarget.LoadFiles())

		updateRequest := &models.DeploymentTargetUpdateRequestV1Alpha{
			DeploymentTargetV1Alpha: *changedTarget,
		}

		updatedTarget, err := c.Update(updateRequest)
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
		"project name; if not specified will be inferred from git origin")
	EditDeploymentTargetCmd.Flags().StringP("project-id", "i", "",
		"project id; if not specified will be inferred from git origin")
	EditDeploymentTargetCmd.Flags().StringP("name", "n", "", "target name")
	EditDeploymentTargetCmd.Flags().StringP("id", "t", "", "target id")
	EditDeploymentTargetCmd.Flags().BoolP("activate", "a", false, "activates/uncordon the deployment target")
	EditDeploymentTargetCmd.Flags().BoolP("deactivate", "d", false, "deactivates/cordon the deployment target")
	editCmd.AddCommand(EditDeploymentTargetCmd)
}
