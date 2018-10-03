package cmd

import (
	"fmt"
	"os"
	"time"

	client "github.com/semaphoreci/cli/api/client"
	models "github.com/semaphoreci/cli/api/models"

	"github.com/semaphoreci/cli/cmd/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const unsetPublicSshKeyMsg = `
Before creating a debug session job, configure the debug.PublicSshKey value.

Examples:

  # Configuring public ssh key with a literal
  sem config set debug.PublicSshKey "ssh-rsa AX3....DD"

  # Configuring public ssh key with a file
  sem config set debug.PublicSshKey "$(cat ~/.ssh/id_rsa.pub)"

  # Configuring public ssh key with your GitHub keys
  sem config set debug.PublicSshKey "$(curl -s https://github.com/<username>.keys)"
`

var debugProjectCmd = &cobra.Command{
	Use:     "project [NAME]",
	Short:   "Debug a project",
	Long:    ``,
	Aliases: []string{"prj", "projects"},
	Args:    cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		publicKey := viper.GetString("debug.PublicSshKey")

		if publicKey == "" {
			fmt.Println("[ERROR] Public SSH key for debugging is not configured.")
			fmt.Println(unsetPublicSshKeyMsg)

			os.Exit(1)
		}

		project_name := args[0]
		pc := client.NewProjectV1AlphaApi()
		project, err := pc.GetProject(project_name)

		utils.Check(err)

		jobName := fmt.Sprintf("Debug Session for %s", project_name)
		job_req := models.NewJobV1Alpha(jobName)

		job_req.Spec = &models.JobV1AlphaSpec{}
		job_req.Spec.Agent.Machine.Type = "e2-standard-2"
		job_req.Spec.Agent.Machine.OsImage = "ubuntu1804"
		job_req.Spec.ProjectId = project.Metadata.Id

		job_req.Spec.Commands = []string{
			fmt.Sprintf("echo '%s' >> .ssh/authorized_keys", publicKey),
			"sleep infinity",
		}

		c := client.NewJobsV1AlphaApi()

		job, err := c.CreateJob(&job_req)

		utils.Check(err)

		for {
			time.Sleep(1000 * time.Millisecond)

			job, err = c.GetJob(job.Metadata.Id)

			utils.Check(err)

			if job.Status.State == "FINISHED" {
				fmt.Printf("Job %s has already finished.\n", job.Metadata.Id)
				os.Exit(0)
			}

			if job.Status.State != "RUNNING" {
				fmt.Printf("Waiting for Job %s to start.\n", job.Metadata.Id)
				continue
			}

			break
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
			time.Sleep(1000 * time.Millisecond)

			err := utils.SshIntoAJob(ip, ssh_port, "semaphore")

			utils.Check(err)
		} else {
			fmt.Printf("Job %s has no exposed SSH port.\n", job.Metadata.Id)
			os.Exit(1)
		}

	},
}

var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "Debug a resource.",
	Long:  ``,
}

func init() {
	debugCmd.AddCommand(debugProjectCmd)

	RootCmd.AddCommand(debugCmd)
}
