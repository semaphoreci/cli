package cmd

import (
	"fmt"
	"os"

	client "github.com/semaphoreci/cli/api/client"
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

		ip := job.Status.Agent.Ip

		var ssh_port int32
		ssh_port = 0

		for _, p := range job.Status.Agent.Ports {
			if p.Name == "ssh" {
				ssh_port = p.Number
			}
		}

		if ip != "" && ssh_port != 0 {
			utils.SshIntoAJob(ip, ssh_port, "semaphore")
		} else {
			fmt.Printf("Job %s has no exposed SSH port.\n", job.Metadata.Id)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(attachCmd)
}
