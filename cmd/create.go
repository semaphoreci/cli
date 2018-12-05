package cmd

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

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

var CreateSecretCmd = &cobra.Command{
	Use:     "secret [NAME]",
	Short:   "Create a secret.",
	Long:    ``,
	Aliases: []string{"secrets"},
	Args:    cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		fileFlags, err := cmd.Flags().GetStringArray("file")
		utils.Check(err)

		var files []models.SecretV1BetaFile
		for _, fileFlag := range fileFlags {
			matchFormat, err := regexp.MatchString(`^[^: ]+:[^: ]+$`, fileFlag)
			utils.Check(err)

			if matchFormat == true {
				flagPaths := strings.Split(fileFlag, ":")

				file := models.SecretV1BetaFile{}
				file.Path = flagPaths[1]
				file.Content = encodeFromFileAt(flagPaths[0])
				files = append(files, file)
			} else {
				utils.Fail("The format of --file flag must be: <local-path>:<semaphore-path>")
			}
		}

		secret := models.NewSecretV1Beta(name, files)

		c := client.NewSecretV1BetaApi()

		_, err = c.CreateSecret(&secret)

		utils.Check(err)

		fmt.Printf("Secret '%s' created.\n", secret.Metadata.Name)
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
		RootCmd.Execute()
	}
}

func encodeFromFileAt(path string) string {
	content, err := ioutil.ReadFile(path)
	utils.Check(err)

	return base64.StdEncoding.EncodeToString(content)
}

func init() {
	createJobCmd := NewCreateJobCmd().Cmd
	createNotificationCmd := NewCreateNotificationCmd()

	RootCmd.AddCommand(createCmd)
	createCmd.AddCommand(CreateSecretCmd)
	createCmd.AddCommand(CreateDashboardCmd)
	createCmd.AddCommand(createJobCmd)
	createCmd.AddCommand(CreateWorkflowCmd)
	createCmd.AddCommand(createNotificationCmd)

	// Create Flags

	desc := "Filename, directory, or URL to files to use to create the resource"
	createCmd.Flags().StringP("file", "f", "", desc)

	// Secret Create Flags

	desc = "File mapping <local-path>:<mount-path>, used to create a secret with file"
	CreateSecretCmd.Flags().StringArrayP("file", "f", []string{}, desc)

	CreateWorkflowCmd.Flags().StringP("project-name", "p", "", "project name; if not specified will be inferred wrom git origin")
	CreateWorkflowCmd.Flags().StringP("label", "l", "", "workflow label")
	CreateWorkflowCmd.Flags().StringP("archive", "a", "", "snapshot archive; if not specified current dir will be compressed into .tgz file")
	CreateWorkflowCmd.Flags().BoolP("follow", "f", false,
		"run 'get ppl <ppl_id>' after create repeatedly until pipeline reaches terminal state")
}
