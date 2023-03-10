package ssh

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	client "github.com/semaphoreci/cli/api/client"
	models "github.com/semaphoreci/cli/api/models"

	"github.com/semaphoreci/cli/cmd/utils"
)

func StartDebugJobSession(debug *models.DebugJobV1Alpha, message string) error {
	c := client.NewJobsV1AlphaApi()
	job, err := c.CreateDebugJob(debug)
	utils.Check(err)

	return StartDebugSession(job, message)
}

func StartDebugProjectSession(debug_project *models.DebugProjectV1Alpha, message string) error {
	c := client.NewJobsV1AlphaApi()
	job, err := c.CreateDebugProject(debug_project)
	utils.Check(err)

	return StartDebugSession(job, message)
}

func StartDebugSession(job *models.JobV1Alpha, message string) error {
	c := client.NewJobsV1AlphaApi()

	defer func() {
		fmt.Printf("\n")
		fmt.Printf("* Stopping debug session ..\n")

		err := c.StopJob(job.Metadata.Id)

		if err != nil {
			utils.Check(err)
		} else {
			fmt.Printf("* Session stopped\n")
		}
	}()

	fmt.Printf("* Waiting for debug session to boot up .")
	err := waitUntilJobIsRunning(job, func() { fmt.Printf(".") })
	if err != nil {
		fmt.Printf("\n[ERROR] %s\n", err)
		return err
	}

	fmt.Println()

	// Reload Job
	job, err = c.GetJob(job.Metadata.Id)
	if err != nil {
		fmt.Printf("\n[ERROR] %s\n", err)
		return err
	}

	/*
	 * If this is a debug session for a self-hosted job, an SSH key will not be available.
	 * The only thing we need to do here is return the self-hosted agent name to the user and hang until the users stops it.
	 */
	if job.IsSelfHosted() {
		fmt.Println(selfHostedSessionMessage(job.AgentName()))
		fmt.Println("Once you are done with the debug session, stop this command and the debug session will be stopped.")
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		return nil
	}

	/*
	 * If this is for a cloud debug session, we grab the SSH key and SSH into the proper machine.
	 */
	sshKey, err := c.GetJobDebugSSHKey(job.Metadata.Id)
	if err != nil {
		fmt.Printf("\n[ERROR] %s\n", err)
		return err
	}

	conn, err := NewConnectionForJob(job, sshKey.Key)
	if err != nil {
		fmt.Printf("\n[ERROR] %s\n", err)
		return err
	}
	defer conn.Close()

	fmt.Printf("* Waiting for ssh daemon to become ready .")
	err = conn.WaitUntilReady(20, func() { fmt.Printf(".") })
	fmt.Println()

	if err != nil {
		fmt.Printf("\n[ERROR] %s\n", err)
		return err
	}

	fmt.Println(message)

	err = conn.Session()

	if err != nil {
		fmt.Printf("\n[ERROR] %s\n", err)
		return err
	}

	return nil
}

func waitUntilJobIsRunning(job *models.JobV1Alpha, callback func()) error {
	var err error

	c := client.NewJobsV1AlphaApi()

	for {
		time.Sleep(1000 * time.Millisecond)

		job, err = c.GetJob(job.Metadata.Id)

		if err != nil {
			continue
		}

		if job.Status.State == "FINISHED" {
			return fmt.Errorf("Job '%s' has already finished.\n", job.Metadata.Id)
		}

		if job.Status.State == "RUNNING" {
			return nil
		}

		// do some processing between ticks
		callback()
	}
}

func selfHostedSessionMessage(agentName string) string {
	return fmt.Sprintf(`
Semaphore CI Self-Hosted Debug Session.

  - The debug session you created is running in the self-hosted agent named '%s'.
  - Once you access the machine where that agent is running, make sure you are logged in as the same user the Semaphore agent is using.
  - Source the '/tmp/.env-*' file where the agent keeps all the environment variables exposed to the job.
  - Checkout your code with `+"`checkout`"+`.

Documentation: https://docs.semaphoreci.com/essentials/debugging-with-ssh-access/.
`, agentName)
}
