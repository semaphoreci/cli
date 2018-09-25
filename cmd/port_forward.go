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
		_, err := c.GetJob(id)

		utils.Check(err)

		ssh_path, err := exec.LookPath("ssh")

		utils.Check(err)

		ip := "94.130.129.220"
		username := "builder"
		port := fmt.Sprintf("-p %d", 29920)
		user_and_ip := fmt.Sprintf("%s@%s", username, ip)
		port_forward_rule := fmt.Sprintf("%s:0.0.0.0:%s", local_port, remote_port)

		// fmt.Println(ssh_path)
		// fmt.Println("-L")
		// fmt.Println(port_forward_rule)
		// fmt.Println(port)
		// fmt.Println("-o")
		// fmt.Println("StrictHostKeyChecking=no")
		// fmt.Println(user_and_ip)

		fmt.Printf("Forwarding %s:%s -> %s:%s", ip, remote_port, "0.0.0.0", local_port)

		ssh_cmd := exec.Command(ssh_path, "-L", port_forward_rule, port, "-o", "StrictHostKeyChecking=no", user_and_ip, "sleep infinity")

		ssh_cmd.Stdin = os.Stdin
		ssh_cmd.Stdout = os.Stdout
		err = ssh_cmd.Start()

		utils.Check(err)

		err = ssh_cmd.Wait()

		utils.Check(err)
	},
}

func init() {
	RootCmd.AddCommand(portForwardCmd)
}
