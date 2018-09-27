package cmd

import (
	"fmt"
	"os"
	"os/exec"

	client "github.com/semaphoreci/cli/api/client"
	"github.com/semaphoreci/cli/cmd/utils"
	"github.com/spf13/cobra"
)

var portForwardCmd = &cobra.Command{
	Use:   "port-forward [JOB ID] [LOCAL PORT] [JOB PORT]",
	Short: "Port forward a local port to a remote port on the job.",
	Long:  ``,
	Args:  cobra.ExactArgs(3),

	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		local_port := args[1]
		remote_port := args[2]

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
			sshAndPortForward(ip, ssh_port, "semaphore", local_port, remote_port)
		} else {
			fmt.Printf("Port forwarding is not possible for job %s.\n", job.Metadata.Id)
			os.Exit(1)
		}
	},
}

func sshAndPortForward(ip string, sshPort int32, username string, localPort string, remotePort string) {
	ssh_path, err := exec.LookPath("ssh")

	utils.Check(err)

	portFlag := fmt.Sprintf("-p %d", sshPort)
	userAndIp := fmt.Sprintf("%s@%s", username, ip)
	portForwardRule := fmt.Sprintf("%s:0.0.0.0:%s", localPort, remotePort)
	noStrictRule := "-oStrictHostKeyChecking=no"

	fmt.Printf("Forwarding %s:%s -> %s:%s...\n", ip, remotePort, "0.0.0.0", localPort)

	ssh_cmd := exec.Command(ssh_path, "-L", portForwardRule, portFlag, noStrictRule, userAndIp, "sleep infinity")

	ssh_cmd.Stdin = os.Stdin
	ssh_cmd.Stdout = os.Stdout
	err = ssh_cmd.Start()

	utils.Check(err)

	err = ssh_cmd.Wait()

	utils.Check(err)
}

func init() {
	RootCmd.AddCommand(portForwardCmd)
}
