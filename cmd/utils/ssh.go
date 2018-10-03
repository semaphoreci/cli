package utils

import (
	"fmt"
	"os"
	"os/exec"
)

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
