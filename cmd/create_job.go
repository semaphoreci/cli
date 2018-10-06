package cmd

import (
	"fmt"

	client "github.com/semaphoreci/cli/api/client"
	models "github.com/semaphoreci/cli/api/models"
	"github.com/semaphoreci/cli/cmd/utils"
	"github.com/spf13/cobra"
)

func NewCreateJobCmd() *cobra.Command {
	commandNotSpecifiedMsg := `Command not specified

Job can't be created without a command. Example:

$ sem create job hello-world --project hello-world --command 'echo "Hello World"'
`

	projectNotSpecifiedMsg := `Project not specified

Example:

  $ sem create job hello-world --project hello-world --command 'echo "Hello World"'
`

	c := &cobra.Command{
		Use:     "job [NAME]",
		Short:   "Create a job.",
		Long:    ``,
		Aliases: []string{"jobs"},

		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]

			projectName, err := cmd.Flags().GetString("project")
			utils.Check(err)

			if projectName == "" {
				utils.Fail(projectNotSpecifiedMsg)
			}

			pc := client.NewProjectV1AlphaApi()
			project, err := pc.GetProject(projectName)

			utils.Check(err)

			command, err := cmd.Flags().GetString("command")
			utils.Check(err)

			if command == "" {
				utils.Fail(commandNotSpecifiedMsg)
			}

			fileFlags, err := cmd.Flags().GetStringArray("file")
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

			job := models.NewJobV1Alpha(name)
			job.Spec = &models.JobV1AlphaSpec{
				ProjectId: project.Metadata.Id,
				Files:     files,
				Secrets:   []models.JobV1AlphaSpecSecret{},
				Commands:  []string{command},
			}

			c := client.NewJobsV1AlphaApi()

			job, err = c.CreateJob(job)

			utils.Check(err)

			fmt.Printf("Job '%s' created.\n", job.Metadata.Id)
		},
	}

	desc := "File mapping <local-path>:<mount-path>"
	c.Flags().StringArrayP("file", "f", []string{}, desc)

	// desc = "Specify environment variable NAME=VALUE"
	// CreateJobCmd.Flags().StringArrayP("e", "env", []string{}, desc)

	desc = "Command to execute in the job"
	c.Flags().StringP("command", "c", "", desc)

	desc = "Project name"
	c.Flags().StringP("project", "p", "", desc)

	return c
}
