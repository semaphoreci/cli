package utils

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	client "github.com/semaphoreci/cli/api/client"
	models "github.com/semaphoreci/cli/api/models"

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

func GetPublicSshKey() (string, error) {
	publicKey := viper.GetString("debug.PublicSshKey")

	if publicKey == "" {
		err := fmt.Errorf("Public SSH key for debugging is not configured.\n\n%s", unsetPublicSshKeyMsg)

		return "", err
	}

	return publicKey, nil
}

func WaitForStartAndSsh(c *client.JobsApiV1AlphaApi, job *models.JobV1Alpha) {
	var err error

	for {
		time.Sleep(1000 * time.Millisecond)

		job, err = c.GetJob(job.Metadata.Id)

		Check(err)

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
		time.Sleep(2000 * time.Millisecond)

		err := SshIntoAJob(ip, ssh_port, "semaphore")

		Check(err)
	} else {
		fmt.Printf("Job %s has no exposed SSH port.\n", job.Metadata.Id)
		os.Exit(1)
	}
}

func SshIntoAJob(ip string, port int32, username string) error {
	ssh_path, err := exec.LookPath("ssh")

	if err != nil {
		return err
	}

	portFlag := fmt.Sprintf("-p%d", port)
	noStrictFlag := "-oStrictHostKeyChecking=no"
	timeoutFlag := "-oConnectTimeout=10"
	userAndIp := fmt.Sprintf("%s@%s", username, ip)

	fmt.Printf("%s %s %s %s %s\n", ssh_path, portFlag, timeoutFlag, noStrictFlag, userAndIp)

	ssh_cmd := exec.Command(ssh_path, portFlag, timeoutFlag, noStrictFlag, userAndIp)

	ssh_cmd.Stdin = os.Stdin
	ssh_cmd.Stdout = os.Stdout
	ssh_cmd.Stderr = os.Stderr

	err = ssh_cmd.Start()

	if err != nil {
		return err
	}

	return ssh_cmd.Wait()
}
