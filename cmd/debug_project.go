package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/semaphoreci/cli/api/models"

	"github.com/semaphoreci/cli/cmd/jobs"
	"github.com/semaphoreci/cli/cmd/ssh"
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

	DebugProjectCmd.Flags().Duration(
		"duration",
		60*time.Minute,
		"duration of the debug session in seconds")

	DebugProjectCmd.Flags().String(
		"machine-type",
		"e1-standard-2",
		"machine type to use for debugging; default: e1-standard-2")

	return DebugProjectCmd
}

func RunDebugProjectCmd(cmd *cobra.Command, args []string) {
	machineType, err := cmd.Flags().GetString("machine-type")
	utils.Check(err)

	if jobs.IsSelfHosted(machineType) {
		fmt.Printf("Self-hosted agent type '%s' can't be used to debug a project. Only cloud agent types are allowed.\n", machineType)
		os.Exit(1)
	}

	duration, err := cmd.Flags().GetDuration("duration")

	utils.Check(err)

	projectNameOrId := args[0]

	utils.Check(err)

	debugPrj := models.NewDebugProjectV1Alpha(projectNameOrId, int(duration.Seconds()), machineType)

	fmt.Printf("* Creating debug session for project '%s'\n", projectNameOrId)
	fmt.Printf("* Setting duration to %d minutes\n", int(duration.Minutes()))

	sshIntroMessage := `
Semaphore CI Debug Session.

  - Checkout your code with ` + "`checkout`" + `
  - Leave the session with ` + "`exit`" + `

Documentation: https://docs.semaphoreci.com/essentials/debugging-with-ssh-access/.
`

	err = ssh.StartDebugProjectSession(debugPrj, sshIntroMessage)
	utils.Check(err)
}
