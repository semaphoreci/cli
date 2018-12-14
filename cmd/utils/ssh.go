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

const unsetPublicSshKeyMsg = `Before creating a debug session job, configure the debug.PublicSshKey value.

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

func WaitForStartAndSsh(c *client.JobsApiV1AlphaApi, job *models.JobV1Alpha, sshIntroMessage string) {
	var err error

	fmt.Printf("* Waiting for debug session to boot up .")

	for {
		time.Sleep(1000 * time.Millisecond)

		job, err = c.GetJob(job.Metadata.Id)

		Check(err)

		if job.Status.State == "FINISHED" {
			fmt.Printf("\nError while starting SSH session Job. Debug Job '%s' is in finished state.\n", job.Metadata.Id)
			os.Exit(1)
		}

		if job.Status.State != "RUNNING" {
			fmt.Print(".")
			continue
		}

		break
	}

	fmt.Printf("\n")

	ip := job.Status.Agent.Ip

	var ssh_port int32
	ssh_port = 0

	for _, p := range job.Status.Agent.Ports {
		if p.Name == "ssh" {
			ssh_port = p.Number
		}
	}

	if ip == "" || ssh_port == 0 {
		fmt.Printf("Job %s has no exposed SSH ports.\n", job.Metadata.Id)
		os.Exit(1)
	}

	fmt.Print(sshIntroMessage)
	fmt.Print("\n")

	time.Sleep(2000 * time.Millisecond)

	err = SshIntoAJob(ip, ssh_port, "semaphore")

	fmt.Printf("\n")
	fmt.Printf("* Stopping debug session ..\n")

	err = c.StopJob(job.Metadata.Id)

	if err != nil {
		Check(err)
	} else {
		fmt.Printf("* Session stopped\n")
	}
}

func SshIntoAJob(ip string, port int32, username string) error {
	ssh_path, err := exec.LookPath("ssh")

	if err != nil {
		return err
	}

	portFlag := fmt.Sprintf("-p%d", port)
	interactive := "-t"
	supressMotd := "-q"
	noStrictFlag := "-oStrictHostKeyChecking=no"
	timeoutFlag := "-oConnectTimeout=10"
	userAndIp := fmt.Sprintf("%s@%s", username, ip)

	ssh_cmd := exec.Command(ssh_path, interactive, supressMotd, portFlag, timeoutFlag, noStrictFlag, userAndIp, "bash --login")

	ssh_cmd.Stdin = os.Stdin
	ssh_cmd.Stdout = os.Stdout
	ssh_cmd.Stderr = os.Stderr

	err = ssh_cmd.Start()

	if err != nil {
		return err
	}

	return ssh_cmd.Wait()
}
