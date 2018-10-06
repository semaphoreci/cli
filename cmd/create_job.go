package cmd

import (
	"fmt"

	client "github.com/semaphoreci/cli/api/client"
	models "github.com/semaphoreci/cli/api/models"
	"github.com/semaphoreci/cli/cmd/utils"
	"github.com/spf13/cobra"
)

type CreateJobCmd struct {
	Cmd *cobra.Command
}

const jobCreateCommandNotSpecifiedMsg = `Command not specified

Job can't be created without a command. Example:

$ sem create job hello-world --project hello-world --command 'echo "Hello World"'
`

const jobCreateProjectNotSpecifiedMsg = `Project not specified

Example:

  $ sem create job hello-world --project hello-world --command 'echo "Hello World"'
`

func NewCreateJobCmd() *CreateJobCmd {
	c := &CreateJobCmd{}

	c.Cmd = &cobra.Command{
		Use:     "job [NAME]",
		Short:   "Create a job.",
		Long:    ``,
		Aliases: []string{"jobs"},
		Run: func(cmd *cobra.Command, args []string) {
			c.Run(args)
		},
	}

	desc := "File mapping <local-path>:<mount-path>"
	c.Cmd.Flags().StringArrayP("file", "f", []string{}, desc)

	// desc = "Specify environment variable NAME=VALUE"
	// CreateJobCmd.Flags().StringArrayP("e", "env", []string{}, desc)

	desc = "Command to execute in the job"
	c.Cmd.Flags().StringP("command", "c", "", desc)

	desc = "Project name"
	c.Cmd.Flags().StringP("project", "p", "", desc)

	return c
}

func (c *CreateJobCmd) Run(args []string) {
	name := args[0]

	job := models.NewJobV1Alpha(name)
	job.Spec = &models.JobV1AlphaSpec{
		ProjectId: c.FindProjectId(),
		Files:     c.ParseFileFlags(),
		Secrets:   []models.JobV1AlphaSpecSecret{},
		Commands:  c.ParseCommandFlag(),
	}

	jobClient := client.NewJobsV1AlphaApi()
	job, err := jobClient.CreateJob(job)

	utils.Check(err)

	fmt.Printf("Job '%s' created.\n", job.Metadata.Id)
}

func (c *CreateJobCmd) FindProjectId() string {
	projectName, err := c.Cmd.Flags().GetString("project")
	utils.Check(err)

	if projectName == "" {
		utils.Fail(jobCreateProjectNotSpecifiedMsg)
	}

	pc := client.NewProjectV1AlphaApi()
	project, err := pc.GetProject(projectName)

	utils.Check(err)

	return project.Metadata.Id
}

func (c *CreateJobCmd) ParseCommandFlag() []string {
	command, err := c.Cmd.Flags().GetString("command")
	utils.Check(err)

	if command == "" {
		utils.Fail(jobCreateCommandNotSpecifiedMsg)
	}

	return []string{command}
}

func (c *CreateJobCmd) ParseFileFlags() []models.JobV1AlphaSpecFile {
	fileFlags, err := c.Cmd.Flags().GetStringArray("file")
	utils.Check(err)

	var files []models.JobV1AlphaSpecFile

	for _, fileFlag := range fileFlags {
		remotePath, content, err := utils.ParseFileFlag(fileFlag)

		utils.Check(err)

		files = append(files, models.JobV1AlphaSpecFile{
			Path:    remotePath,
			Content: content,
		})
	}

	return files
}
