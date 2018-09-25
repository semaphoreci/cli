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
		_, err := c.GetJob(id)

		utils.Check(err)

		ssh_path, err := exec.LookPath("ssh")

		utils.Check(err)

		ip := "94.130.129.220"
		username := "builder"
		port := fmt.Sprintf("-p %d", 29920)
		user_and_ip := fmt.Sprintf("%s@%s", username, ip)

		ssh_cmd := exec.Command(ssh_path, port, "-o", "StrictHostKeyChecking=no", user_and_ip)

		ssh_cmd.Stdin = os.Stdin
		ssh_cmd.Stdout = os.Stdout
		err = ssh_cmd.Start()

		utils.Check(err)

		err = ssh_cmd.Wait()

		utils.Check(err)
	},
}

func init() {
	RootCmd.AddCommand(attachCmd)
}
