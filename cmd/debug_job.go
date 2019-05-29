package cmd

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	client "github.com/semaphoreci/cli/api/client"
	models "github.com/semaphoreci/cli/api/models"

	"github.com/semaphoreci/cli/cmd/ssh"
	"github.com/semaphoreci/cli/cmd/utils"
	"github.com/spf13/cobra"
)

func NewDebugJobCmd() *cobra.Command {
	var DebugJobCmd = &cobra.Command{
		Use:     "job [ID]",
		Short:   "Debug a job",
		Long:    ``,
		Aliases: []string{"job", "jobs"},
		Args:    cobra.ExactArgs(1),
		Run:     RunDebugJobCmd,
	}

	DebugJobCmd.Flags().Duration(
		"duration",
		60*time.Minute,
		"duration of the debug session in seconds")

	return DebugJobCmd
}

func RunDebugJobCmd(cmd *cobra.Command, args []string) {
	duration, err := cmd.Flags().GetDuration("duration")
	utils.Check(err)

	jobId := args[0]

	c := client.NewJobsV1AlphaApi()
	oldJob, err := c.GetJob(jobId)

	utils.Check(err)

	jobName := fmt.Sprintf("Debug Session for Job %s", jobId)
	job := models.NewJobV1Alpha(jobName)

	// Copy everything to new job, except commands
	job.Spec = oldJob.Spec
	job.Spec.EpilogueAlwaysCommands = []string{}
	job.Spec.EpilogueOnPassCommands = []string{}
	job.Spec.EpilogueOnFailCommands = []string{}

	// Construct a commands file and inject into job
	commandsFileContent := strings.Join([]string{
		strings.Join(oldJob.Spec.Commands, "\n"),
		"# epilogue always commands",
		strings.Join(oldJob.Spec.EpilogueAlwaysCommands, "\n"),
		"# epilogue on pass commands",
		strings.Join(oldJob.Spec.EpilogueOnPassCommands, "\n"),
		"# epilogue on fail commands",
		strings.Join(oldJob.Spec.EpilogueOnFailCommands, "\n"),
	}, "\n")

	job.Spec.Files = []models.JobV1AlphaSpecFile{
		models.JobV1AlphaSpecFile{
			Path:    "commands.sh",
			Content: base64.StdEncoding.EncodeToString([]byte(commandsFileContent)),
		},
	}

	// Overwrite commands with a sleep. This will keep the job up for N seconds.
	// Original commands are inserted into commands.sh.
	job.Spec.Commands = []string{
		fmt.Sprintf("sleep %d", int(duration.Seconds())),
	}

	fmt.Printf("* Creating debug session for job '%s'\n", jobId)
	fmt.Printf("* Setting duration to %d minutes\n", int(duration.Minutes()))

	sshIntroMessage := `
Semaphore CI Debug Session.

  - Checkout your code with ` + "`checkout`" + `
  - Run your CI commands with ` + "`source ~/commands.sh`" + `
  - Leave the session with ` + "`exit`" + `

Documentation: https://docs.semaphoreci.com/article/75-debugging-with-ssh-access.
`

	ssh.StartDebugSession(job, sshIntroMessage)
}
