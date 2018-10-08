package cmd

import (
	"encoding/base64"
	"fmt"
	"strings"

	client "github.com/semaphoreci/cli/api/client"
	models "github.com/semaphoreci/cli/api/models"

	"github.com/semaphoreci/cli/cmd/utils"
	"github.com/spf13/cobra"
)

var DebugJobCmd = &cobra.Command{
	Use:     "job [NAME]",
	Short:   "Debug a job",
	Long:    ``,
	Aliases: []string{"job", "jobs"},
	Args:    cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		publicKey, err := utils.GetPublicSshKey()

		utils.Check(err)

		jobId := args[0]

		c := client.NewJobsV1AlphaApi()
		oldJob, err := c.GetJob(jobId)

		utils.Check(err)

		jobName := fmt.Sprintf("Debug Session for Job %s", jobId)
		job := models.NewJobV1Alpha(jobName)

		// Copy everything to new job, except commands
		job.Spec = oldJob.Spec
		job.Spec.EpilogueCommands = []string{}

		job.Spec.Commands = []string{
			fmt.Sprintf("echo '%s' >> .ssh/authorized_keys", publicKey),
			"sleep infinity",
		}

		// Construct a commands file and inject into job
		commandsFileContent := fmt.Sprintf("%s\n%s",
			strings.Join(oldJob.Spec.Commands, "\n"),
			strings.Join(oldJob.Spec.EpilogueCommands, "\n"))

		job.Spec.Files = []models.JobV1AlphaSpecFile{
			models.JobV1AlphaSpecFile{
				Path:    "commands.sh",
				Content: base64.StdEncoding.EncodeToString([]byte(commandsFileContent)),
			},
		}

		job, err = c.CreateJob(job)

		utils.Check(err)

		utils.WaitForStartAndSsh(&c, job)
	},
}
