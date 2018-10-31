package cmd

import (
	"fmt"

	client "github.com/semaphoreci/cli/api/client"
	models "github.com/semaphoreci/cli/api/models"
	"github.com/semaphoreci/cli/cmd/utils"
	"github.com/spf13/cobra"
)

func NewDebugProjectCmd() *cobra.Command {
	var DebugProjectCmd = &cobra.Command{
		Use:     "project [NAME]",
		Short:   "Debug a project",
		Long:    ``,
		Aliases: []string{"prj", "projects"},
		Args:    cobra.ExactArgs(1),
		Run:     RunDebugProjectCmd,
	}

	DebugProjectCmd.Flags().String(
		"machine-type",
		"e1-standard-2",
		"machine type to use for debugging; default: e1-standard-2")

	return DebugProjectCmd
}

func RunDebugProjectCmd(cmd *cobra.Command, args []string) {
	publicKey, err := utils.GetPublicSshKey()
	machineType, err := cmd.Flags().GetString("machine-type")

	utils.Check(err)

	project_name := args[0]
	pc := client.NewProjectV1AlphaApi()
	project, err := pc.GetProject(project_name)

	utils.Check(err)

	jobName := fmt.Sprintf("Debug Session for %s", project_name)
	job := models.NewJobV1Alpha(jobName)

	job.Spec = &models.JobV1AlphaSpec{}
	job.Spec.Agent.Machine.Type = machineType
	job.Spec.Agent.Machine.OsImage = "ubuntu1804"
	job.Spec.ProjectId = project.Metadata.Id

	job.Spec.Commands = []string{
		fmt.Sprintf("echo '%s' >> .ssh/authorized_keys", publicKey),
		"sleep infinity",
	}

	c := client.NewJobsV1AlphaApi()

	job, err = c.CreateJob(job)

	utils.Check(err)

	utils.WaitForStartAndSsh(&c, job)
}
