package ssh

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	models "github.com/semaphoreci/cli/api/models"
)

type Connection struct {
	IP       string
	Port     int32
	Username string
}

func NewConnection(ip string, port int32, username string) (*Connection, error) {
	return &Connection{
		IP:       ip,
		Port:     port,
		Username: username,
	}, nil
}

func NewConnectionForJob(job *models.JobV1Alpha) (*Connection, error) {
	ip := job.Status.Agent.Ip

	var port int32
	port = 0

	for _, p := range job.Status.Agent.Ports {
		if p.Name == "ssh" {
			port = p.Number
		}
	}

	if ip == "" || port == 0 {
		err := fmt.Errorf("Job %s has no exposed SSH ports.\n", job.Metadata.Id)

		return nil, err
	}

	return NewConnection(ip, port, "semaphore")
}

func (c *Connection) WaitUntilReady(attempts int, callback func()) error {
	var err error
	var ok bool

	for i := 0; i < attempts-1; i++ {
		time.Sleep(1 * time.Second)

		ok, err = c.IsReady()

		if ok && err == nil {
			return nil
		} else {
			callback()
		}
	}

	return fmt.Errorf("SSH connection can't be established; %s", err)
}

func (c *Connection) IsReady() (bool, error) {
	cmd, err := c.sshCommand("echo 'success'", false)
	output, err := cmd.CombinedOutput()

	if err == nil && strings.Contains(string(output), "success") {
		return true, nil
	} else {
		outputOneLine := string(output)
		// remove empty spaces from ends
		outputOneLine = strings.TrimSpace(outputOneLine)
		// remove \r
		outputOneLine = strings.Replace(outputOneLine, "\r", "", -1)
		// join lines
		outputOneLine = strings.Replace(outputOneLine, "\n", ", ", -1)
		// remove trailing '.'
		outputOneLine = strings.Trim(outputOneLine, ".")

		return false, fmt.Errorf("%s; %s", outputOneLine, err)
	}
}

func (c *Connection) Session() error {
	cmd, err := c.sshCommand("bash --login", true)

	if err != nil {
		return err
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Start()

	if err != nil {
		return err
	}

	return cmd.Wait()
}

func (c *Connection) sshCommand(directive string, interactive bool) (*exec.Cmd, error) {
	path, err := exec.LookPath("ssh")

	if err != nil {
		return nil, err
	}

	portFlag := fmt.Sprintf("-p%d", c.Port)

	interactiveFlag := ""
	if interactive {
		interactiveFlag = "-t"
	} else {
		interactiveFlag = "-T"
	}

	noStrictFlag := "-oStrictHostKeyChecking=no"
	timeoutFlag := "-oConnectTimeout=5"
	userAndIp := fmt.Sprintf("%s@%s", c.Username, c.IP)

	cmd := exec.Command(
		path,
		interactiveFlag,
		portFlag,
		timeoutFlag,
		noStrictFlag,
		userAndIp,
		directive)

	return cmd, nil
}
