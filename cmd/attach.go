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

		// Get SSH key for job
		sshKey, err := c.GetJobDebugSSHKey(job.Metadata.Id)
		utils.Check(err)

		conn, err := ssh.NewConnectionForJob(job, sshKey)
		utils.Check(err)
		defer conn.Close()

		conn.Session()
	},
}

func init() {
	RootCmd.AddCommand(attachCmd)
}
