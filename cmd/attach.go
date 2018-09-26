package cmd

import (
	"fmt"
	"os"
	"os/exec"

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

		ip := job.Status.Agent.Ip

		var ssh_port int32
		ssh_port = 0

		for _, p := range job.Status.Agent.Ports {
			if p.Name == "ssh" {
				ssh_port = p.Number
			}
		}

		if ip != "" && ssh_port != 0 {
			sshIntoAJob(ip, ssh_port, "semaphore")
		}
	},
}

func sshIntoAJob(ip string, port int32, username string) {
	ssh_path, err := exec.LookPath("ssh")

	utils.Check(err)

	portFlag := fmt.Sprintf("-p%d", port)
	noStrictFlag := "-oStrictHostKeyChecking=no"
	userAndIp := fmt.Sprintf("%s@%s", username, ip)

	ssh_cmd := exec.Command(ssh_path, portFlag, noStrictFlag, userAndIp)

	ssh_cmd.Stdin = os.Stdin
	ssh_cmd.Stdout = os.Stdout
	err = ssh_cmd.Start()

	utils.Check(err)

	err = ssh_cmd.Wait()

	utils.Check(err)
}

func init() {
	RootCmd.AddCommand(attachCmd)
}
