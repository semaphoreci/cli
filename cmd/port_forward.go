package cmd

import (
	"fmt"
	"io/ioutil"
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
		localPort := args[1]
		remotePort := args[2]

		c := client.NewJobsV1AlphaApi()
		job, err := c.GetJob(id)
		utils.Check(err)

		sshKey, err := c.GetJobDebugSSHKey(job.Metadata.Id)

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

		var sshPort int32
		sshPort = 0

		for _, p := range job.Status.Agent.Ports {
			if p.Name == "ssh" {
				sshPort = p.Number
			}
		}

		if ip != "" && sshPort != 0 {
			sshAndPortForward(ip, sshPort, "semaphore", localPort, remotePort, sshKey.Key)
		} else {
			fmt.Printf("Port forwarding is not possible for job %s.\n", job.Metadata.Id)
			os.Exit(1)
		}
	},
}

func sshAndPortForward(ip string, sshPort int32, username string, localPort string, remotePort string, sshKey string) {
	sshPath, err := exec.LookPath("ssh")

	utils.Check(err)

	portFlag := fmt.Sprintf("-p %d", sshPort)
	userAndIP := fmt.Sprintf("%s@%s", username, ip)
	portForwardRule := fmt.Sprintf("%s:0.0.0.0:%s", localPort, remotePort)
	noStrictRule := "-oStrictHostKeyChecking=no"
	identityOnlyrule := "-oIdentitiesOnly=yes"

	sshKeyFile, _ := ioutil.TempFile("", "sem-cli-debug-private-key")
	defer os.Remove(sshKeyFile.Name())
	_, err = sshKeyFile.Write([]byte(sshKey))
	utils.Check(err)
	sshKeyFile.Close()
	utils.Check(err)

	fmt.Printf("Forwarding %s:%s -> %s:%s...\n", ip, remotePort, "0.0.0.0", localPort)
	sshCmd := exec.Command(sshPath, "-L", portForwardRule, portFlag, noStrictRule, identityOnlyrule, userAndIP, "-i", sshKeyFile.Name(), "sleep infinity")

	sshCmd.Stdin = os.Stdin
	sshCmd.Stdout = os.Stdout
	err = sshCmd.Start()

	utils.Check(err)

	err = sshCmd.Wait()

	utils.Check(err)
}

func init() {
	RootCmd.AddCommand(portForwardCmd)
}
