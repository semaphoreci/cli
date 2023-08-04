package cmd

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/semaphoreci/cli/cmd/utils"
	"github.com/semaphoreci/cli/cmd/workflows"

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

		// #nosec
		data, err := ioutil.ReadFile(path)

		utils.CheckWithMessage(err, "Failed to read from resource file.")
		_, kind, err := utils.ParseYamlResourceHeaders(data)
		utils.Check(err)

		switch kind {
		case "Project":
			project, err := models.NewProjectV1AlphaFromYaml(data)

			utils.Check(err)

			c := client.NewProjectV1AlphaApi()

			_, err = c.CreateProject(project)

			utils.Check(err)

			fmt.Printf("Project '%s' created.\n", project.Metadata.Name)
		case "Notification":
			notif, err := models.NewNotificationV1AlphaFromYaml(data)

			utils.Check(err)

			c := client.NewNotificationsV1AlphaApi()

			notif, err = c.CreateNotification(notif)

			utils.Check(err)

			fmt.Printf("Notification '%s' created.\n", notif.Metadata.Name)
		case "Secret":
			secret, err := models.NewSecretV1BetaFromYaml(data)

			utils.Check(err)

			c := client.NewSecretV1BetaApi()

			_, err = c.CreateSecret(secret)

			utils.Check(err)

			fmt.Printf("Secret '%s' created.\n", secret.Metadata.Name)
		case "ProjectSecret":
			secret, err := models.NewProjectSecretV1FromYaml(data)

			utils.Check(err)

			c := client.NewProjectSecretV1Api(secret.Metadata.ProjectIdOrName)

			_, err = c.CreateSecret(secret)

			utils.Check(err)

			fmt.Printf("Secret '%s' created in project '%s'.\n", secret.Metadata.Name, secret.Metadata.ProjectIdOrName)
		case "Dashboard":
			dash, err := models.NewDashboardV1AlphaFromYaml(data)

			utils.Check(err)

			c := client.NewDashboardV1AlphaApi()

			_, err = c.CreateDashboard(dash)

			utils.Check(err)

			fmt.Printf("Dashboard '%s' created.\n", dash.Metadata.Name)
		case "Job":
			job, err := models.NewJobV1AlphaFromYaml(data)

			utils.Check(err)

			c := client.NewJobsV1AlphaApi()

			job, err = c.CreateJob(job)

			utils.Check(err)

			fmt.Printf("Job '%s' created.\n", job.Metadata.Id)
		case "SelfHostedAgentType":
			at, err := models.NewAgentTypeV1AlphaFromYaml(data)
			utils.Check(err)

			c := client.NewAgentTypeApiV1AlphaApi()
			newAgentType, err := c.CreateAgentType(at)
			utils.Check(err)

			y, err := newAgentType.ToYaml()
			utils.Check(err)
			fmt.Printf("%s", y)
		case models.DeploymentTargetKindV1Alpha:
			target, err := models.NewDeploymentTargetV1AlphaFromYaml(data)
			utils.Check(err)
			if target == nil {
				utils.Check(errors.New("deployment target in the file is empty"))
				return
			}
			createRequest := &models.DeploymentTargetCreateRequestV1Alpha{
				DeploymentTargetV1Alpha: *target,
			}
			utils.Check(createRequest.LoadFiles())
			c := client.NewDeploymentTargetsV1AlphaApi()
			createdDeploymentTarget, err := c.Create(createRequest)
			utils.Check(err)

			y, err := createdDeploymentTarget.ToYaml()
			utils.Check(err)
			fmt.Printf("Deployment target '%s' (%s) created:\n%s\n", createdDeploymentTarget.Id, createdDeploymentTarget.Name, y)
		default:
			utils.Fail(fmt.Sprintf("Unsupported resource kind '%s'", kind))
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

var CreateAgentTypeCmd = &cobra.Command{
	Use:     "agent_type [NAME]",
	Short:   "Create a self-hosted agent type.",
	Long:    ``,
	Aliases: []string{"agenttype", "agentType"},
	Args:    cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		c := client.NewAgentTypeApiV1AlphaApi()
		at := models.NewAgentTypeV1Alpha(name)
		agentType, err := c.CreateAgentType(&at)
		utils.Check(err)

		y, err := agentType.ToYaml()
		utils.Check(err)
		fmt.Printf("%s", y)
	},
}

var CreateWorkflowCmd = &cobra.Command{
	Use:     "workflow [NAME]",
	Short:   "Create a workflow from snapshot.",
	Long:    ``,
	Aliases: []string{"workflows", "wf"},
	Args:    cobra.ExactArgs(0),

	Run: func(cmd *cobra.Command, args []string) {
		projectName, err := cmd.Flags().GetString("project-name")
		utils.Check(err)

		archiveName, err := cmd.Flags().GetString("archive")
		utils.Check(err)

		follow, err := cmd.Flags().GetBool("follow")
		utils.Check(err)

		if projectName == "" {
			projectName, err = utils.InferProjectName()
			utils.Check(err)
		}

		label, err := cmd.Flags().GetString("label")
		utils.Check(err)

		createSnapshot(projectName, label, archiveName, follow)
	},
}

func createSnapshot(projectName, label, archiveName string, follow bool) {
	log.Printf("Project name: %s\n", projectName)

	body, err := workflows.CreateSnapshot(projectName, label, archiveName)
	utils.Check(err)

	if follow == false {
		fmt.Println(string(body))
	} else {
		body, err := models.NewWorkflowSnapshotResponseV1AlphaFromJson(body)
		utils.Check(err)

		RootCmd.SetArgs([]string{"get", "ppl", body.PplID, "--follow"})
		_ = RootCmd.Execute()
	}
}

func encodeFromFileAt(path string) string {
	// #nosec
	content, err := ioutil.ReadFile(path)
	utils.Check(err)

	return base64.StdEncoding.EncodeToString(content)
}

func init() {
	createJobCmd := NewCreateJobCmd().Cmd
	createNotificationCmd := NewCreateNotificationCmd()
	createSecretCmd := NewCreateSecretCmd()

	RootCmd.AddCommand(createCmd)
	createCmd.AddCommand(createSecretCmd)
	createCmd.AddCommand(CreateDashboardCmd)
	createCmd.AddCommand(createJobCmd)
	createCmd.AddCommand(CreateWorkflowCmd)
	createCmd.AddCommand(createNotificationCmd)
	createCmd.AddCommand(CreateAgentTypeCmd)
	createCmd.AddCommand(NewCreateDeploymentTargetCmd())

	// Create Flags

	desc := "Filename, directory, or URL to files to use to create the resource"
	createCmd.Flags().StringP("file", "f", "", desc)

	CreateWorkflowCmd.Flags().StringP("project-name", "p", "", "project name; if not specified will be inferred wrom git origin")
	CreateWorkflowCmd.Flags().StringP("label", "l", "", "workflow label")
	CreateWorkflowCmd.Flags().StringP("archive", "a", "", "snapshot archive; if not specified current dir will be compressed into .tgz file")
	CreateWorkflowCmd.Flags().BoolP("follow", "f", false,
		"run 'get ppl <ppl_id>' after create repeatedly until pipeline reaches terminal state")
}
