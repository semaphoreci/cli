package cmd

import (
	"fmt"
	"time"

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

	debug := models.NewDebugJobV1Alpha(jobId, int(duration.Seconds()))

	fmt.Printf("* Creating debug session for job '%s'\n", jobId)
	fmt.Printf("* Setting duration to %d minutes\n", int(duration.Minutes()))

	sshIntroMessage := `
Semaphore CI Debug Session.

  - Checkout your code with ` + "`checkout`" + `
  - Run your CI commands with ` + "`source ~/commands.sh`" + `
  - Leave the session with ` + "`exit`" + `

Documentation: https://docs.semaphoreci.com/article/75-debugging-with-ssh-access.
`

	ssh.StartDebugJobSession(debug, sshIntroMessage)
}
