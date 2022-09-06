package cmd

import (
	"fmt"
	"os"

	client "github.com/semaphoreci/cli/api/client"
	"github.com/semaphoreci/cli/cmd/ssh"
	"github.com/semaphoreci/cli/cmd/utils"
	"github.com/spf13/cobra"
)

var attachCmd = &cobra.Command{
	Use:   "attach [JOB ID]",
	Short: "Attach to a running job.",
	Long:  ``,
	Args:  cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]

		c := client.NewJobsV1AlphaApi()
		job, err := c.GetJob(id)

		utils.Check(err)

		if job.Status.State == "FINISHED" {
			fmt.Printf("Job %s has already finished.\n", job.Metadata.Id)
			os.Exit(1)
		}

		if job.Status.State != "RUNNING" {
			fmt.Printf("Job %s has not yet started.\n", job.Metadata.Id)
			os.Exit(1)
		}

		/*
		 * If we want to attach to a machine where a self-hosted job is running, no SSH key will be available.
		 * We just give the agent name back to the user and return.
		 */
		if job.IsSelfHosted() {
			fmt.Printf("* Job '%s' is running in the self-hosted agent named '%s'.\n", id, job.AgentName())
			fmt.Printf("* Once you access the machine where that agent is running, make sure you are logged in as the same user the Semaphore agent is using.\n")
			fmt.Printf("* You can source the '/tmp/.env-*' file where the agent keeps all the environment variables exposed to the job.\n")
			return
		}

		/*
		 * If this is for a cloud job, we go for the SSH key.
		 */
		sshKey, err := c.GetJobDebugSSHKey(job.Metadata.Id)
		utils.Check(err)

		conn, err := ssh.NewConnectionForJob(job, sshKey.Key)
		utils.Check(err)
		defer conn.Close()

		conn.Session()
	},
}

func init() {
	RootCmd.AddCommand(attachCmd)
}
